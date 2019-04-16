# IOEX Wallet Development Environment Document

## 開發伺服器獲取

1. 使用郵箱將 `ssh public key` 發送給 `service@ioex.co` 請求分配開發伺服器，管理員分配成功後回復相應的伺服器位址和存取方法。

    通常推薦使用證書方式登錄伺服器，如果需要使用密碼登錄也可以在郵件中進行說明。

2. 根據需求後會給您分配兩台伺服器：`鏈節點伺服器` 和 `錢包伺服器`。

## IOEX節點服務

### 環境說明

1. 應用程式預設放置在 `/opt/ioex_node` 目錄中，一個運行中的包含命令列用戶端的節點目錄內容如下：

    ```bash
    ls -al /opt/ioex_node

    total 24924
    drwxrwxr-x 4 ioex ioex     4096 Apr 20 17:19 ./
    drwxrwxr-x 5 ioex ioex     4096 Apr 20 11:50 ../
    drwxr-xr-x 2 ioex ioex     4096 Apr 20 12:08 Chain/
    -rw-rw-r-- 1 ioex ioex       34 Apr 20 11:51 cli-config.json
    -rw-rw-r-- 1 ioex ioex     1088 Apr 20 12:02 config.json
    -rwxrwxr-x 1 ioex ioex 13565344 Apr 20 11:51 ioex-cli*
    -rw-rw-r-- 1 ioex ioex      444 Apr 20 12:01 keystore.dat
    drwxrw-r-- 2 ioex ioex     4096 Apr 20 12:07 Log/
    -rwxrwxr-x 1 ioex ioex 11648196 Apr 20 11:50 ioex*
    -rw-r--r-- 1 ioex ioex   274432 Apr 20 17:19 wallet.db
    ```

   * `ioex` 是鏈節點的啟動程式，相應的設定檔為 `config.json`。

   * `ioex-cli` 是一個命令列管理工具，可以用來創建錢包，轉帳和查看區塊資料，設定檔為`cli-config.json`。

   * `keystore.dat` 和 `wallet.db` 是錢包的秘鑰和錢包數據。

   * `Chain` 和 `Log` 目錄是區塊鏈帳本資料以及日誌資料。

   * 鏈節點和命令列用戶端代碼和編譯方法以及配置和命令列參數可以參考：

     * [ioeX.MainChain](https://github.com/ioeXNetwork/ioeX.MainChain/blob/master/README.md)
     * [ioeX.Client](https://github.com/ioeXNetwork/ioeX.Client/blob/master/README.md)

2. 預設情況下分配給你的節點伺服器會連接到一個專門用於開發的一組節點伺服器，自動挖礦並且把所獲取的幣發送到本地的錢包中；你可以參考 [ioeX.Client](https://github.com/ioeXNetwork/ioeX.Client/blob/master/README.md) 中的相關查詢命令進行查詢或者其他操作。

   備註：部分命令列操作需要提供一個錢包密碼，預設是 `ioex`。

3. 節點伺服器使用如下埠提供服務：

    ```bash
    2[*]334  提供 Http Rest api
    2[*]335  提供 Web Socket api
    2[*]336  用於命令列用戶端連接到節點獲取資料
    2[*]338  提供節點之間的資料同步
    ```

   備註：`*` 表示這個數字不確定，會在成功分配伺服器後確定。

### Node API

可參考 API示例：

* [IOEX_Wallet_Node_API](IOEX_Wallet_Node_API_CN.md)

## IOEX錢包服務

1. IOEX錢包服務包含兩個前臺的服務：`browser` 和 `wallet`；以及另外三個後台工作的服務 `worker`、 `api`、 `notification`。

2. 服務簡介：

    ```bash
    browser:  區塊流覽器使用者介面，可以查看區塊資料，查詢交易
    wallet:   web錢包使用者介面，可以進行錢包創建，資產查詢，轉帳操作
    worker:   從節點獲取區塊資訊，寫入緩存資料庫
    api:      為 wallet 提供介面
    notification:
    ```

3. 上面的這些服務的應用程式都在 `/data/www/` 目錄下並使用 `xxxx.test.ioex.org` 作為目錄名稱；目前所有服務都是使用 `supervisor` 進行管理，相關設定檔在 `/etc/supervisor/conf.d` 目錄下，可以使用 `supervisorctl [start|stop|restart] <service_name>` 來啟動或者停止某個服務：

    ```bash
    supervisorctl stop worker.test.ioex.org

    supervisorctl restart all
    ```

   使用 `supervisor` 管理服務的詳細文檔可以在 [這裡](http://www.supervisord.org/) 查看。

4. IOEX錢包的 `web` 伺服器使用 `nginx`, 相關設定檔在 `/etc/nginx/conf.d` 目錄。

   開發環境的 `https` 服務使用了自簽名的證書，所以使用流覽器進行訪問時需要手動添加一下對證書的信任，另外分配的開發環境默認不會使用標準的`80` `443`埠，所以訪問時候要加上相應的埠。

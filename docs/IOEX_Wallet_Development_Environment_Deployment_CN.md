# IOEX Wallet Development Environment Deployment

## 錢包部署

1. 參考 `docker` 官方 [安裝文檔](https://docs.docker.com/install/linux/docker-ce/ubuntu/)，安裝 `Docker CE`

2. 準備相關文件

    ```bash
    $ sudo ls -al /opt/docker_ioex_wallet

    total 36
    drwxrwxr-x 6 ioex ioex 4096 Apr 20 10:18 ./
    drwxr-xr-x 3 root   root   4096 Apr 20 14:05 ../
    -rw-rw-r-- 1 ioex ioex  395 Apr 19 16:50 authorized_keys
    -rw-rw-r-- 1 ioex ioex  502 Apr 20 10:18 Dockerfile
    drwxrwxr-x 2 ioex ioex 4096 Apr 17 16:05 ioex.co/
    drwxrwxr-x 2 ioex ioex 4096 Apr 17 16:05 nginx/
    drwxrwxr-x 2 ioex ioex 4096 Apr 19 18:09 supervisor/
    -rwxrwxr-x 1 ioex ioex 1314 Apr 19 18:55 walletWrapper.sh*
    drwxr-xr-x 7 ioex ioex 4096 Apr 13 12:16 www/

    # authorized_keys 檔包含用戶提供的 openssh 公開金鑰，啟動鏡像時會映射到鏡像中，這樣用戶可以通過 openssh 連接到容器
    # ioex.co 目錄中包含 https 訪問錢包和流覽器所需的伺服器憑證和私密金鑰
    # nginx 目錄中包含錢包相關應用的設定檔，可以從錢包相關代碼中獲得
    # supervisor 目錄中包含 supervisor 相關設定檔，可以從錢包代碼中獲得
    # www 目錄中包含錢包相關應用的代碼
    ```

   **`Dockerfile 檔內容如下:`**

    ```Dockerfile
    FROM ubuntu:16.04

    RUN apt-get update && apt-get install -y vim net-tools openssh-server iputils-ping

    RUN apt-get install -y libxml2-dev libxslt-dev python-dev python-pip libjpeg-dev libcurl4-openssl-dev libgeos-dev libmysqlclient-dev supervisor python nginx mongodb-server redis-server curl rsyslog

    RUN curl -sL https://deb.nodesource.com/setup_6.x | bash - && apt-get install nodejs

    COPY walletWrapper.sh walletWrapper.sh

    CMD ["./walletWrapper.sh"]
    ```

   **`walletWrapper.sh 檔內容如下:`**

    ```bash
    #!/bin/bash

    # Start the sshd
    /etc/init.d/ssh start
    status=$?
    if [ $status -ne 0 ]; then
      echo "Failed to start sshd: $status"
      exit $status
    fi

    # Start the all service
    /etc/init.d/rsyslog start
    if [ $? -ne 0 ];then
      echo "Rsyslog start failed"
      exit 1
    fi

    /etc/init.d/nginx start
    if [ $? -ne 0 ];then
      echo "Nginx start failed"
      exit 1
    fi

    /etc/init.d/mongodb start
    if [ $? -ne 0 ];then
      echo "Mongodb start failed"
      exit 1
    fi

    /etc/init.d/redis-server start
    if [ $? -ne 0 ];then
      echo "Redis start failed"
      exit 1
    fi

    sleep 20

    /etc/init.d/supervisor start
    if [ $? -ne 0 ];then
      echo "Supervisor start failed"
      exit 1
    fi

    while sleep 60; do
      ps aux |grep sshd |grep -q -v grep
      PROCESS_1_STATUS=$?
      # If the greps above find anything, they exit with 0 status
      # If they are not both 0, then something is wrong
      if [ $PROCESS_1_STATUS -ne 0 ]; then
        echo "Process has already exited."
        exit 1
      fi
    done
    ```

3. 創建 `docker` 鏡像

    ```bash
    $ cd /opt/docker_ioex_wallet

    $ docker build -t ioex_wallet_run_01 .
    # 這裡在安裝nodejs時可能會失敗，需要通過代理連接
    ```

4. 啟動 `docker` 鏡像

    ```bash
    $ docker run -m 2g --cpus=2 -p 922:22 -p 20443:443 -p 20080:80 \
    -v /opt/docker_ioex_wallet/authorized_keys:/root/.ssh/authorized_keys \
    -v /opt/docker_ioex_wallet/ioex.co:/etc/ssl/ioex.co \
    -v /opt/docker_ioex_wallet/nginx:/etc/nginx/conf.d \
    -v /opt/docker_ioex_wallet/supervisor:/etc/supervisor/conf.d \
    -v /opt/docker_ioex_wallet/www:/data/www \
    ioex_wallet_run_01
    ```

   * 這裡限制容器可以使用的記憶體數量為2G，使用的cpu數量為2個.
   * 注意： 如果把所有服務都部署在同一台伺服器的話，需要修改 `notification` 的監聽埠（否則會和 `api` 埠衝突），然後修改 `api` 服務連接 `notificationUrl`的埠.

## 鏈節點部署

1. 參考 `docker` 官方 [安裝文檔](https://docs.docker.com/install/linux/docker-ce/ubuntu/)，安裝 `Docker CE`

2. 準備相關文件

    ```bash
    $ sudo ls -al /opt/docker_ioex_node/

    total 24
    drwxrwxr-x 3 ioex ioex 4096 Apr 20 10:17 ./
    drwxr-xr-x 4 root   root   4096 Apr 20 14:42 ../
    -rw-rw-r-- 1 ioex ioex  395 Apr 19 19:18 authorized_keys
    -rw-rw-r-- 1 ioex ioex  213 Apr 19 19:34 Dockerfile
    drwxrwxr-x 4 ioex ioex 4096 Apr 19 19:37 ela_node/
    -rwxrwxr-x 1 ioex ioex 1221 Apr 17 11:56 nodeWrapper.sh*

    # authorized_keys 檔包含用戶提供的 openssh 公開金鑰，創建鏡像時會複製到鏡像中，這樣用戶可以通過 openssh 連接到容器
    # ioex_node 目錄中包含鏈節點相關應用和設定檔
    ```

   * 鏈節點和命令列用戶端代碼和編譯方法以及配置和命令列參數可以參考:

     * [ioeX.MainChain](../README.md)

     * [ioeX.Client](https://github.com/ioeXNetwork/ioeX.Client/blob/master/README.md)

   **`Dockerfile檔內容如下`**

    ```Dockerfile
    FROM ubuntu:16.04

    RUN apt-get update && apt-get install -y vim net-tools openssh-server iputils-ping

    COPY nodeWrapper.sh nodeWrapper.sh

    CMD ["./nodeWrapper.sh"]
    ```

   **`nodeWrapper.sh檔內容如下`**

    ```bash
    #!/bin/bash

    # Start the sshd
    /etc/init.d/ssh start
    status=$?
    if [ $status -ne 0 ]; then
      echo "Failed to start sshd: $status"
      exit $status
    fi

    # Start the IOEX_NODE
    if [ ! -d "/opt/ioex_node" ];then
      echo "IOEX_NODE program not exists";
      exit
    fi

    cd /opt/ioex_node/ && ./node -p ioex > /dev/null &

    status=$?
    if [ $status -ne 0 ]; then
      echo "Failed to start IOEX_NODE: $status"
      exit $status
    fi

    while sleep 60; do
      ps aux |grep sshd |grep -q -v grep
      PROCESS_1_STATUS=$?
      ps aux |grep "./node" |grep -q -v grep
      PROCESS_2_STATUS=$?
      if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
        echo "One of the processes has already exited."
        exit 1
      fi
    done
    ```

3. 創建 `docker` 鏡像

    ```bash
    cd /opt/docker_ioex_node
    $ docker build -t ioex_node_run_01 .
    ```

4. 啟動 `docker` 鏡像

    ```bash
    $ docker run -m 512m --cpus=2 -p 922:22 -p 20338:20338 -p 20334:20334 -p 20335:20335 \
    -v /opt/docker_ioex_node/authorized_keys:/root/.ssh/authorized_keys \
    -v /opt/docker_ioex_node/ioex_node:/opt/ioex_node \
    ioex_node_run_01
    ```

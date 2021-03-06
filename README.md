# ioeX.MainChain

## Summary

ioeX leveraged Elastos functions to create its owned features and also business cases.

ioeXcoin is the digital currency solution within ioeX eco system.

This project is the source code that can build a full node of IOEX.

## Build on Mac

### Check OS version

Make sure the OSX version is 16.7+

```shell
$ uname -srm
Darwin 16.7.0 x86_64
```

### Install Go distribution 1.9

Use Homebrew to install Golang 1.9.

```shell
$ brew install go@1.9
```

> If you install older version, such as v1.8, you may get missing math/bits package error when build.

### Setup basic workspace
In this instruction we use ~/dev/src/github.com/ioeXNetwork as our working directory. If you clone the source code to a different directory, please make sure you change other environment variables accordingly (not recommended). 

```shell
$ mkdir -p ~/dev/bin
$ mkdir -p ~/dev/src/github.com/ioeXNetwork/
```

### Set correct environment variables

```shell
export GOROOT=/usr/local/opt/go@1.9/libexec
export GOPATH=$HOME/dev
export GOBIN=$GOPATH/bin
export PATH=$GOROOT/bin:$PATH
export PATH=$GOBIN:$PATH
```

### Install Glide

Glide is a package manager for Golang. We use Glide to install dependent packages.

```shell
$ brew install --ignore-dependencies glide
```

### Check Go version and glide version

Check the golang and glider version. Make sure they are the following version number or above.

```shell
$ go version
go version go1.9.2 darwin/amd64

$ glide --version
glide version 0.13.1
```

If you cannot see the version number, there must be something wrong when install.

### Clone source code to $GOPATH/src/github.com/ioex folder
Make sure you are in the folder of $GOPATH/src/github.com/ioeXNetwork
```shell
$ git clone http://github.com/ioeXNetwork/ioeX.MainChain.git
```

If clone works successfully, you should see folder structure like $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain/Makefile
### Glide install

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `glide update && glide install` to install dependencies.

### Make

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `make` to build files.

If you did not see any error message, congratulations, you have made the IOEX full node.

## Run on Mac

- run `./ioex` to run the node program.

## Build on Ubuntu

### Check OS version

Make sure your ubuntu version is 16.04+

```shell
$ cat /etc/issue
Ubuntu 16.04.3 LTS \n \l
```

### Install basic tools

```shell
$ sudo apt-get install -y git
```

### Install Go distribution 1.9

```shell
$ sudo apt-get install -y software-properties-common
$ sudo add-apt-repository -y ppa:gophers/archive
$ sudo apt update
$ sudo apt-get install -y golang-1.9-go
```

> If you install older version, such as v1.8, you may get missing math/bits package error when build.

### Setup basic workspace
In this instruction we use ~/$GOPATH/src/github.com/ioeXNetwork/ as our working directory. If you clone the source code to a different directory, please make sure you change other environment variables accordingly (not recommended). 

```shell
$ mkdir -p ~/$GOPATH/bin
$ mkdir -p ~/$GOPATH/src/github.com/ioeXNetwork
```

### Set correct environment variables

```shell
export GOROOT=/usr/lib/go-1.9
export GOPATH=$HOME/dev
export GOBIN=$GOPATH/bin
export PATH=$GOROOT/bin:$PATH
export PATH=$GOBIN:$PATH
```

### Install Glide

Glide is a package manager for Golang. We use Glide to install dependent packages.

```shell
$ cd ~/dev
$ curl https://glide.sh/get | sh
```

### Check Go version and glide version

Check the golang and glider version. Make sure they are the following version number or above.

```shell
$ go version
go version go1.9.2 linux/amd64

$ glide --version
glide version v0.13.1
```

If you cannot see the version number, there must be something wrong when install.

### Clone source code to $GOPATH/src/github.com/ioeXNetwork folder
Make sure you are in the folder of $GOPATH/src/github.com/ioeXNetwork
```shell
$ git clone http://github.com/ioeXNetwork/ioeX.MainChain.git
```

If clone works successfully, you should see folder structure like $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain/Makefile
### Glide install

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `glide update && glide install` to install dependencies.

### Make

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `make` to build files.

If you did not see any error message, congratulations, you have made the ioeXNetwork full node.

## Run on Ubuntu

- run `./ioex` to run the node program.

## Build on Windows

### Install Go distribution 1.9

download Go v1.9 package
install Go

> If you install older version, such as v1.8, you may get missing math/bits package error when build.

### Setup basic workspace
In this instruction we use $GOPATH/src/github.com/ioeXNetwork/ as our working directory. If you clone the source code to a different directory, please make sure you change other environment variables accordingly (not recommended). 

```shell
$ mkdir -p $GOPATH/bin
$ mkdir -p $GOPATH/src/github.com/ioeXNetwork
```

### Install Glide

Glide is a package manager for Golang. We use Glide to install dependent packages.

```shell
$ go get github.com/Masterminds/glide
```

### Glide install

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `glide update && glide install` to install dependencies.

### Clone source code to $GOPATH/src/github.com/ioeXNetwork folder
Make sure you are in the folder of $GOPATH/src/github.com/ioeXNetwork

```shell
$ git clone http://github.com/ioeXNetwork/ioeX.MainChain.git
```

If clone works successfully, you should see folder structure like $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain

### Make

cd $GOPATH/src/github.com/ioeXNetwork/ioeX.MainChain and Run `go build -o ioex.exe main.go` to build files.

If you did not see any error message, congratulations, you have made the ioeXNetwork full node.

# Config the node

See the [documentation](./docs/config.json.md) about config.json

## Bootstrap using docker

Alternatively if don't want to build it manually. We also provide a `Dockerfile` to help you (You need have a prepared docker env).

```bash
cd docker
docker build -t ioex_node_run .
```

#start container

docker run -p 20334:20334 -p 20335:20335 -p 20336:20336 -p 20338:20338 ioex_node_run
```

> Note: don't using Ctrl-C to terminate the output, just close this terminal and open another.

Now you can access IOEX Node's rest api:

```bash
curl http://localhost:20334/api/v1/block/height
```

In the above instruction, we use default configuration file `config.json` in the repository; If you familiar with docker you can change the docker file to use your own IOEX Node configuration.

## More

If you want to learn the API of ioeX.MainChain, please refer to the following:

- [IOEX_Wallet_Node_API](docs/IOEX_Wallet_Node_API_CN.md)

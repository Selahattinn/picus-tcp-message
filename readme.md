# tcp-chat-app

## For The General Project
Implementation of 1-1 messaging app over TCP/IP. There is a server, which acts as a broker between clients

```
├─ bin           //The folder where the binary files was created
├─ cmd           //The code that started it all
├─ config.yml    //Config file for backend
├─ go.mod        //3rd party libraries
├─ go.sum        //Sums and versions of 3rd party libraries
├─ makefile      //MakeFile for build,test and version control 
└─ pkg
   ├─ client                 // Client class files  
   ├─ crypto                 // for encrypt and decrypt      
   ├─ model                  //Models for every type of object
   ├─ repository             //DB Layer
   │  ├─ message
   ├─ server                 //Server Layer for all aplication.
   ├─ service                //Service Layer
   │  ├─ message
   └─ version                //Version control&save for git

```

## ⚡️ Quick start

First of all, [download](https://golang.org/dl/) and install **Go**. :)

## Pre-Req
> Update & Upgrade OS
```bash
sudo apt update && sudo apt upgrade -y
```
> Install Mysql
```bash
sudo apt install mysql-server
```
> mysql password configuration
```bash
mysql -u root -p
mysql> UPDATE mysql.user SET authentication_string = PASSWORD('passwd')
WHERE User = 'root' AND Host = 'localhost';
mysql>FLUSH PRIVILEGES;
```
## For build

```bash
make build
```
## For Test

```bash
make test
```

## Running tcp-chat-app-backend
Retrieves other information from the config.yml file
```shell
./bin/server [-config.file string] [-log.file string] [-debug]  [-version]


-config.file : Get neccessary information from this file (default: config.yml)
-log.file : Log all outputs and errors 
(default: tcp-message-server.log)
-debug : Changes to log level (default: false)
-version : shows version information (default: false)

Example version command:
./bin/server -version
tcp-message-server, version  (branch: master, revision: f1027dac56c17c35f29d8a4ee21e37f2da86c678)
  build user:       selo
  build date:       20220116-21:36:58
  go version:       go1.13.8
```


## Running tcp-chat-app-client
> 1- Using telnet
```shell
telnet localhost 8080
```

> 2- Using client binary
```shell
./bin/client [-addr string] [-name string]

-addr : server address (default: Empty)
-name : user name (default:Empty)

Example Commands:
./bin/client -name Test -addr :8080
./bin/client -name Test -addr localhost:8080
./bin/client -name Test -addr 127.0.0.1:8080
```




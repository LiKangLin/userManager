Simple User-Management system
===========================

汇报文档地址
------------
```
https://docs.google.com/document/d/1uz0_Un8GvUz6nId92yAtEmPf0huj26ijOHdPn4_e8OI/edit#heading=h.q5wxou2y8rlj
```

Requirements
------------
* bash
* Go(v1.13)
* Mysql(v5.5+)
* Redis(v3.4+)
* go moudle

Installation
------------
```$xslt
    #makesure database and redis is correctlly installed and started
    sudo sysctl -w kern.ipc.somaxconn=2048
    sudo sysctl -w kern.maxfiles=12288
    ulimit -n 10000


    #setup environment
    cd src
    sh install.sh
```

Start
------------
```$xslt
    sh start.sh

    initdb.go:用于创建数据库表和插入数据
```

Performance testing
------------
install apache tools
```
    性能测试工具：WRK
       WRK的安装：
       /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
       brew install wrk

```

test with web client
```
    http://localhost:8080/static/login.html
    username/pwd : username138/123456
```

test with curl
```
    curl -d "username=username138&passwd=123456" "http://localhost:8080/login"
```

performance test with wrk


for login test
```
wrk -t12 -c15 -d20s --latency -s post.lua http://localhost:8080/login
```

post.lua 脚本文件
```
wrk.method="POST"
wrk.body='username=username138&passwd=123456&cherr=1'
wrk.headers["Content-Type"]="application/x-www-form-urlencoded"
function request()
 return wrk.format('POST',nil,nil,body)
end
```

汇报文档地址：https://docs.google.com/document/d/1uz0_Un8GvUz6nId92yAtEmPf0huj26ijOHdPn4_e8OI/edit#
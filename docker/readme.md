### 开始
```shell script
    # 进入当前目录
    docker-compose build
    # 启动mysql 创建数据库
    docker-compose up mysql
    # 进入go环境容器打包
    docker run -it --rm -v  $(pwd)/../:/data  docker_go1  /bin/bash
    # 设置国内镜像
    # go env -w  GOPROXY=https://goproxy.cn
    # 构建
    go build main.go
    # 退出容器
    exit
    cd ../
    # 复制程序到./docker/data/go
    cp main ./docker/data/go/
    # 复制修改配置文件
    cp config.example.yaml ./docker/data/go/config.yaml
    # 重启
    docker-compose down
    docker-compose up
    # 进入go容器，创建表，假数据
    docker exec -it goserve1 /bin/bash
    ./main migrate 
    ./main fake
```
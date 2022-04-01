### 开始
```shell script
    # 进入当前目录
    docker-compose build
    # 启动mysql然后手动创建数据库
    docker-compose up mysql
    # 构建
    docker run -it --rm -v  $(pwd)/../:/data  docker-compose_go1  /bin/bash -c "go env -w  GOPROXY=https://goproxy.cn && go build main.go" && cp ../main ./data/go
    # 复制修改配置文件
    cp ../config.example.yaml ./data/go/config.yaml
    # 表迁移，填入假数据
    docker run -it --rm -v  $(pwd)/data/go:/data --network docker-compose_default docker-compose_go1 /bin/bash -c "./main migrate && ./main fake"
    # 启动
    docker-compose up
```
### 集群
默认情况下会启动4个go服务，如不需要则在docker-composer.yaml里面注释掉,
并修改go配置config.yaml以及services/nginx/nginx.conf里面upstream重新up即可

### 前端代码
克隆下来启动即可
- [客服端](https://github.com/zjwshisb/service-frontend.git)
- [用户端](https://github.com/zjwshisb/service-user.git)

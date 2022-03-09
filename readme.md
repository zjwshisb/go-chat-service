# gochat客服系统

### 介绍
基于golang，ant-design-pro，Taro实现的在线客服IM系统，通过websocket实现消息接受发送。支持集群部署。
- [客服前端代码仓库](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码仓库](https://github.com/zjwshisb/service-user) 

依赖
- mysql
- redis
- rabbitmq(可选)
- etcd(可选)

### 开始
```shell script
    # 复制，修改 config.yaml 配置信息
    cp config.example.yaml config.yaml
    go build main
    # 创建表
    ./main migrate
    # 插入测试数据
    ./main fake
    # 启动服务
    ./main serve
    # 停止服务
    ./main stop
```

### 部署
nginx部署可参考目录下nginx.conf

    
### 功能
- 图片发送，emoji表情，快捷回复
- 自定义自动回复
- 转接人工(排队位置显示)
- 客服转接
- 离线消息提醒
- 多开提醒(重复登录，多个tab等)
- 多租户


### 演示地址
用户端(移动端): [http://119.29.196.153/mobile](http://119.29.196.153/mobile)  
账号: user(1-20) #user1-user20  
密码: user(1-20) #user1-user20  
客服端(pc): [http://119.29.196.153/admin](http://119.29.196.153/admin)  
客服界面在管理后台右上角客服面板点击进入   
账号: admin(1-20) #admin1-admin20  
密码: admin(1-20) #admin1-admin20  
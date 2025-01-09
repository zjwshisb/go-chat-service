# gochat客服系统

### 介绍
基于goframe(v2)，ant-design-pro，Taro实现的在线客服IM系统，通过websocket实现消息接受发送，支持集群部署。

- [客服前端代码](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码](https://github.com/zjwshisb/service-user)

依赖
- mysql
- redis
- etcd|polaris|zookeeper(集群部署情况下,默认etcd)

集群部署情况下，通过grpc实现不同服务之间的websocket消息传递发送。

### 开始
```shell
git clone https://github.com/zjwshisb/go-chat-service.git
cd go-chat-service
cp manifest/config/config.example.yaml manifest/config/config.yaml # 复制并修改配置文件
go run main.go migrate # 迁移数据结构
go run main.go fake # 假数据填充
go run main.go http # 启动服务
```

### 功能

- [x] 文字/表情/图片/音频/视频消息发送
- [x] 自动回复
- [x] 用户在线状态显示
- [x] 消息已读未读显示(支持配置开关)
- [X] 转接人工(支持排队位置显示，支持配置开关)
- [x] 客服转接功能
- [x] 多开提醒(重复登录，多个tab等)
- [X] 多租户
- [x] 文件管理


### 演示地址

用户端(移动端): [http://120.77.242.145/mobile](http://120.77.242.145/mobile)  
账号: user(0-19) #user0-user19  
密码: user(0-19) #user1-user19  
客服端(pc): [http://120.77.242.145/server](http://120.77.242.145/server)  
账号: admin(0-19) #admin1-admin19  
密码: admin(0-19) #admin1-admin19

### 更新日志

2025.1.4 使用[GOFRAME](https://github.com/gogf/gf)进行了重构  
2025.1.9 支持集群部署

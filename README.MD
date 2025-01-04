# gochat客服系统

### 介绍
基于golang，ant-design-pro，Taro实现的在线客服IM系统，通过websocket实现消息接受发送。  
使用GoFrame进行了重构，更易于二次开发。  
旧版本在分支old上。   
- [客服前端代码](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码](https://github.com/zjwshisb/service-user)

依赖
- mysql
- redis

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
- [x] 用户在线状态显示
- [x] 消息已读未读显示(支持配置开关)
- [X] 转接人工(支持排队位置，支持配置开关)
- [x] 客服转接功能
- [x] 多开提醒(重复登录，多个tab等)
- [X] 多租户
- [ ] 集群部署 

持续完善中...

### 演示地址

用户端(移动端): [http://120.77.242.145/mobile](http://120.77.242.145/mobile)  
账号: user(1-20) #user1-user20  
密码: user(1-20) #user1-user20  
客服端(pc): [http://120.77.242.145/server](http://120.77.242.145/server)  
账号: admin(1-20) #admin1-admin20  
密码: admin(1-20) #admin1-admin20

### 更新日志

2025.1.4 使用[GOFRAME](https://github.com/gogf/gf)进行了重构

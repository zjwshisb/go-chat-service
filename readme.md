# gochat客服系统

### 介绍
基于golang，ant-design-pro，Taro实现的在线客服IM系统，通过websocket实现消息接受发送，支持集群部署。
- [客服前端代码](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码](https://github.com/zjwshisb/service-user) 

依赖
- mysql
- redis
- rabbitmq(可选)
- etcd(可选)

### 开始
详见[docker-compose/readme.md](https://github.com/zjwshisb/go-chat-service/tree/master/docker-compose)
，已集成所有所需服务，开箱即用。

### 部署
nginx部署可参考目录下nginx.conf

    
### 功能
- 图片发送，emoji表情，快捷回复
- 自定义自动回复
- 转接人工(排队位置显示)
- 客服转接
- 离线消息提醒
- 用户上下线提醒  
- 多开提醒(重复登录，多个tab等)
- 多租户等

### update
4.28 在本地环境下新增一个简易监控面板(localhost/monitor)，可查看所有websocket连接数

### 演示地址
用户端(移动端): [http://119.29.196.153/mobile](http://119.29.196.153/mobile)  
账号: user(1-20) #user1-user20  
密码: user(1-20) #user1-user20  
客服端(pc): [http://119.29.196.153/admin](http://119.29.196.153/admin)  
客服界面在管理后台右上角客服面板点击进入   
账号: admin(1-20) #admin1-admin20  
密码: admin(1-20) #admin1-admin20 


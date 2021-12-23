# gochat客服系统

### 介绍
基于golang，ant-design-pro，Taro实现的在线客服IM系统，通过websocket实现消息接受发送。支持集群部署。
- [客服前端代码仓库](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码仓库](https://github.com/zjwshisb/service-user) 

依赖
- mysql
- redis
- rabbitmq(可选)

### 开始
```shell script
    # 复制，修改 config.ini 配置信息
    cp config.ini.example config.ini
    # 执行./commands/migrate/main [-c=yourpath/config.ini] 创建表
    # 执行./commands/fake/main [-c=yourpath/config.ini] 创建测试数据
    go build main
    ./main [-c=yourpath/config.ini]
```

### 部署
nginx部署可参考目录下nginx.conf   
systemctl管理见chat.service

    
### 功能
- 自动回复
- 转接人工
- 图片发送
- 客服转接
- 离线消息提醒
- 多开提醒(重复登录，多个tab等)
- 其他


### 演示地址
用户端(移动端): [http://119.29.196.153/mobile](http://119.29.196.153/mobile)  
账号: user(1-20) #user1-user20  
密码: user(1-20) #user1-user20  
客服端(pc): [http://119.29.196.153/admin](http://119.29.196.153/admin)  
客服界面在管理后台右上角客服面板点击进入   
账号: admin(1-20) #admin1-admin20  
密码: admin(1-20) #admin1-admin20  
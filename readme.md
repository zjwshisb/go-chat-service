# websocket实现的在线客服系统
基于golang(websocket，http)实现的在线客服系统，已在生产环境中使用  
适用与不满足微信小程序自带的客服系统或多端需要共用一套客服系统的需求  
需要与下面配套前端代码一起使用   
- [客服前端代码](https://github.com/zjwshisb/service-frontend)
- [用户端前端代码](https://github.com/zjwshisb/service-user) 

依赖   
- mysql
- redis

### 开始
```shell script
    # 复制，修改 config.ini 配置信息
    cp config.ini.example config.ini 
    go build main
    # 创建表 -u(是否创建用户表)
    ./main -m=migrate -u 
    # 插入必备的数据
    ./main -m=seeder
    # 假用户数据 
    ./main -m=fake
    # 启动
    ./main [-m=start] [-c=yourpath/config.ini]
    # 停止
    ./main -m=stop [-c=yourpath/config.ini]
   
```

### 部署
详见目录下nginx.conf，chat.service


### 功能
- 自动回复
- 转接人工
- 用户分组(需要根据业务二开，详见代码)
- 图片发送(支持七牛，本地存储)
- 客服转接
- 用户离线消息提示(小程序订阅消息实现)
- 客服界面禁止多开(重复登录，多个tab等)
- 其他详见演示地址


### 演示地址
用户端(移动端): [http://119.29.196.153/mobile](http://119.29.196.153/mobile)  
账号: user(1-20) #user1-user20  
密码: user(1-20) #user1-user20  
客服端(pc): [http://119.29.196.153/admin](http://119.29.196.153/admin)  
客服界面在管理后台右上角客服面板点击进入   
账号: admin(1-20) #admin1-admin20  
密码: admin(1-20) #admin1-admin20  
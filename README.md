## 运行Spring Cloud Config Server
```bash
cd config-server
mvn spring-boot:run
```

## 运行Viper客户端
```bash
cd viper-client
go run *.go 
```


访问 `http://localhost:8080/resume`，以获取Spring Cloud Config Server中对应微服务的配置信息。


## 运行改进版客户端-支持热更新
```bash
cd update-client
go run *.go
```
更新`config-repo`之后，触发Actuator刷新指定客户端事件，发送:
```bash
curl http://localhost:8889/actuator/busrefresh/client-demo:8081 -X POST --noproxy "*"
```

排错及验证步骤：
```bash
# 验证config-server已经获取新配置
curl http://localhost:8888/client-demo/dev --noproxy "*"

# 验证客户端是否读取到更新后的信息
curl http://localhost:8081/resume --noproxy "*"
```
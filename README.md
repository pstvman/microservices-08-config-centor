运行Spring Cloud Config Server
```bash
cd config-server
mvn spring-boot:run
```

运行Viper客户端
```bash
cd viper-client
go run *.go 
```


访问 `http://localhost:8080/resume`，以获取Spring Cloud Config Server中对应微服务的配置信息。
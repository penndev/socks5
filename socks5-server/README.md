# socks5
rfc1928 协议的Java实现

> 支持 ipv6 和 CONNECT, UDP, BIND。

```bash
mvn clean 
mvn install
java -jar target/socks5-1.0.jar 1080 admin 123456
```

** 待实现 **

- GSSAPI 登录
- BID CMD
- UDP FLAG
- IPv6 验证
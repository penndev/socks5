## Windows

`tun`设备驱动

**https://github.com/WireGuard/wintun**

设置网卡IP命令

```
netsh interface ipv4 set address name="wintun" source=static addr=10.10.10.10 mask=255.255.255.255

\ netsh interface ipv4 add route 0.0.0.0/0 10.10.10.10
\ route add 0.0.0.0/0 10.10.10.10
```

防止环路

```
route add {SOCKS5_IP} mask 255.255.255.255 0.0.0.0 metric 1
```


切换路由

```
    route add 0.0.0.0 mask 128.0.0.0 10.10.10.10 metric 1
```


## Android

> 编译为Android aar库使用

```
go install golang.org/x/mobile/cmd/gomobile@latest
go get golang.org/x/mobile/bind
gomobile init
gomobile bind -o mobileStack.aar -target android -androidapi 24 ./mobileStack
```
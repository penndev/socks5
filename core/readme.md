## 协议包含

- socks4 `socks4://127.0.0.1:1080`
- socks5(tls)
    - `socks5://127.0.0.1:1080`
    - `socks5tls://Domain:1080`
- http(tls)
    - `http://127.0.0.1:1080`
    - `https://Domain:1080`


## Windows

`tun`设备驱动 **https://github.com/WireGuard/wintun**


设置网卡IP命令

```
netsh interface ipv4 set address name="wintun" source=static addr=172.19.0.1 mask=255.255.255.255

route add 0.0.0.0/0 172.19.0.1
```

防止环路

```
route add {SOCKS5_IP} mask 255.255.255.255 0.0.0.0 metric 1
```


切换路由

```
    route add 0.0.0.0 mask 128.0.0.0 10.10.10.10 metric 1
```





## darwin

> ifconfig utun6 198.18.0.1 198.18.0.1 up
> ifconfig utun6 inet 172.19.0.1 172.19.0.1 netmask 255.255.255.255 up
netstat -rn -f inet

sudo route add -net 1.0.0.0/8 172.19.0.1
sudo route add -net 2.0.0.0/7 172.19.0.1
sudo route add -net 4.0.0.0/6 172.19.0.1
sudo route add -net 8.0.0.0/5 172.19.0.1
sudo route add -net 16.0.0.0/4 172.19.0.1
sudo route add -net 32.0.0.0/3 172.19.0.1
sudo route add -net 64.0.0.0/2 172.19.0.1
sudo route add -net 128.0.0.0/1 172.19.0.1
sudo route add -net 198.18.0.0/15 172.19.0.1


## Android

> 编译为Android aar库使用

```
go install golang.org/x/mobile/cmd/gomobile@latest
go get golang.org/x/mobile/bind
gomobile init
gomobile bind -o mobileStack.aar -target android -androidapi 24 ./mobileStack
```
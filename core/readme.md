
`tun`设备驱动

**https://github.com/WireGuard/wintun**



设置网卡IP

```
netsh interface ipv4 set address name="wintun" source=static addr=10.10.100.251 mask=255.255.255.255

netsh interface ipv4 add route 103.235.46.40/32  10.10.100.251
\ route add 103.235.46.40/32  10.10.100.251
```

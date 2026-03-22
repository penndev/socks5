package main

// netsh interface ipv4 set address name="PrismTUN" source=static addr=172.19.0.1 mask=255.255.255.255
const TUN_NAME = "utun"
const TUN_MTU = 1500
const TUN_OFFSET = 4

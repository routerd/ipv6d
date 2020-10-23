# Important

Set this in /etc/systctl.conf:

`net.ipv6.conf.eth0.accept_ra=2`


```
-A PREROUTING -d 2003:cd:d721:14fc::/64 -i eth0 -j NETMAP --to fd9c:fd74:6b8d:10::/64
-A PREROUTING -d 2003:cd:d721:14fe::/64 -i eth0 -j NETMAP --to fd9c:fd74:6b8d:3::/64

ip6tables -t nat -A POSTROUTING -s fd9c:fd74:6b8d:10::/64 -o eth0 -j NETMAP --to 2003:cd:d721:14fc::/64

-A POSTROUTING -s fd9c:fd74:6b8d:3::/64 -o eth0 -j NETMAP --to 2003:cd:d721:14fe::/64
```



ip6tables -t nat -N IPV6D-OUTBOUND
ip6tables -t nat -A POSTROUTING -j IPV6D-OUTBOUND
ip6tables -t nat -A IPV6D-OUTBOUND -s fd9c:fd74:6b8d:10::/64 -o eth0 -j NETMAP --to 2003:cd:d721:14fc::/64

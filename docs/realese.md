- Now version v0.3.5

## v0.3.5

add rotate report file,  by default set 10. 

```
local_trableshoot --count-rotate=20
```

## v0.3.4

add arg and env 

if you set container 

show utilization rates of the top 10 by cpu and mem container in docker 
```
local_trableshoot --container=docker
```

Tracing to DNS specified in /etc/resolv.conf, default set false
```
local_trableshoot --check-dns
```

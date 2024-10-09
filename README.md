## local_trableshoot
Local trableshoot linux system  and save to  local report for html format.

the file contains a data report:
 - atop  
 - session 
 - atop (cpu/mem/IOps) for top. 
 - top 10 containers used cpu&mem (for now docker)
 - process top cpu and process tree 
 - process top mem
 - tcp/udp/sockets connect 
 - connecting ip 
 - routes/Neighbours/Resolver
 - dns check  
 - device (used,innodes,mount)
 - kernel used modules and dmesg last messages

example report in [this](./docs/example/report_tooz-Aspire-V3-571G_09.10.2024_09:28:44.html) 

You can use it together with the [monit](https://www.mmonit.com/monit/) service. When check LOADAVG 

example:
```
check system {{ monit_hostname }}
    if loadavg (1min) > {{ monit_highload_la_1m }} then exec /usr/local/sbin/local_trableshoot
```
other examples, you can see [here](https://www.mmonit.com/monit/documentation/monit.html)

##  This project helps in troubleshooting the system

- when there is an abnormal load on the system and it is difficult to understand the true cause of the problem through standard monitoring systems
- written by go and description of the structure according to go [standards](https://github.com/golang-standards/project-layout/blob/master/README.md) 

## Howto work and use 

Run application 

```
./local_trableshoot
```

Apllication save report file in /var/log/report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html

## initialize the project and build 

```go
go mod init local_trableshoot
go mod tidy // это решение по зависимостей 
go build -o ./local_trableshoot ./cmd/app
```



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

## Flags&Env

| Args             | Variable         | Type    | Default | Description      |
|------------------|------------------|---------|---------|------------------|
| container        | container        | string  | ""      | Specify container runtime, top 10 cpu&mem usage, (e.g. docker) |
| check-dns      | check-dns          | Bool    | false   | checking dns availability from /etc/resolv.conf |
| count-rotate   | count-rotate       | int     | 10      | Apllication save report file in /var/log/report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html . The number of report files is no more than 10 pieces, all older files are of the following format: /var/log/report_*.html , will be deleted.  |



## Description of releases

Description of releases and new features [this](./docs/realese.md) 

##  This project helps in troubleshooting the system

- when there is an abnormal load on the system and it is difficult to understand the true cause of the problem through standard monitoring systems
- written by go and description of the structure according to go [standards](https://github.com/golang-standards/project-layout/blob/master/README.md) 

## Howto work and use 

Download and unpack

```sh 
# set version 
version_trableshoot="v0.3.5"
wget -qO- https://github.com/KKulishov/local_trableshoot/releases/download/$version_trableshoot/local_trableshoot.tar.gz | sudo tar xvz -C /usr/local/sbin --strip-components=1 && rm -f local_trableshoot.tar.gz
```

check Run application 

```
sudo local_trableshoot 
```

Apllication save report file in /var/log/report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html . The number of report files is no more than 10 pieces, all older files are of the following format: /var/log/report_*.html , will be deleted. 

If need change the number of old reports to delete, please set:
```
sudo local_trableshoot --count-rotate=20
```

If need check dns and top used cpu&mem in container, used args:
```
sudo local_trableshoot  --check-dns --container=docker
```

## initialize the project and build 

If you need re build , can use this man (go version 1.22):

```go
go mod init local_trableshoot
go mod tidy // это решение по зависимостей 
go build -o ./local_trableshoot ./cmd/app
```



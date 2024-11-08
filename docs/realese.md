- Now version v0.3.6

## v0.3.6

- adding 2 report forms, full and short (report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html and full_report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html)
- adding top 10 disk and network utilization in docker containers
- output of kernel errors and system messages
- Added version output to the report and cli command (--version)
- Added argument in which directory to store reports (--report-dir)

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

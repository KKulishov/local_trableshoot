- Now version v0.4.9

## v0.4.9

-- add:  performance test cpu for simple report. 

## v0.4.8

-- fix: save 2 local tcpdump. Don't upload tcpdump in s3. 

## v0.4.7

-- add: Experimental network analyze. use perf&tcpdump (for ksoftirqd analyzes pid) and upload to s3 

## v0.4.6

-- fix: parallel launch of report unloading in s3 as soon as they are ready

## v0.4.5

-- add: mapping pid of process by cpu&memory top utilization to name namespace/pods/container 
-- add: separate argument to start rotation s3 buсket 

if specified with the launch key
```
--run-rotate-s3
``` 
rotation will run in s3 bucket, but the report will not be generated

convenient to put on the cron. for weekly run to clean objects

example cron task (This task will be executed once on Saturday at midnight):
```
SHELL=/bin/bash
PATH=/sbin:/bin:/usr/sbin:/usr/bin
0 0 * * 6 root /usr/local/sbin/local_trableshoot --run-rotate-s3
```

## v0.4.4

-- changes: dns check new logic , add library connection [dns](https://github.com/miekg/dns)

add args check-dns-name
```
--check-dns-name="domain.com"
```
We get a list of /etc/resolv.conf , by nameserver.
Checking domain name resolution, if it doesn't work, it traces to the unavailable DNS server.  

## v0.4.3

-- add: args --atop-report if you need report atopsar top utilization CPU/MEM/IO/NET in the last 15 minutes
-- add: args --proxyS3Host if you use s3 proxy, by analogy https://github.com/nginxinc/nginx-s3-gateway
-- add: General network statistics 
-- changes: I also reworked the check-dns logic, now tracing to the host's DNS only occurs if the connection via udp port 53 does not reach the host.
-- changes: reduced the search depth in the 1 hour slice for errors in logs (/var/log/messages, /var/log/kernel.log, /var/log/kern.log)

## v0.4.0

- add: Io top disck utilization for 5 sec. 
- add: upload s3 report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html && /var/log/full_report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html. 
- add: Load Average system 
- add: Simple navigation menu

## v0.3.9

fix: s3 to upload charset, set UTF-8. 

## v0.3.8

Upload s3 report file

verified minio 

create file in ~/.config/report_send_s3

example:
```
endpoint_url = s3.ru-1.storage.selcloud.ru
access_key_id = login
secret_access_key = password
use_ssl =  true
bucket_name = name_bucket
```

if file ~/.config/report_send_s3  is in the system, then the program will try to download the report to s3 

By default, file rotation occurs when the value is above 30 files.


You can specify the quantity
```sh
sudo local_trableshoot --count-rotate-s3=20
```

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

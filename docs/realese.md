- Now version v0.4.0

## v0.4.0

- add: Io top disck utilization for 5 sec. 
- add: upload s3 report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html && /var/log/full_report_{{ name_host }}_{{ dd.mm.yyyy_hh.mm.ss }}.html. 
- add: Load Average system 

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

package flags

import (
	"local_trableshoot/configs"

	"github.com/alecthomas/kingpin/v2"
)

var (
	version        = configs.Version
	ContainerFlag  = kingpin.Flag("container", "Specify container runtime (e.g. docker)").Envar("container").Default("").String()
	CheckDns       = kingpin.Flag("check-dns", "Tracing to DNS specified in /etc/resolv.conf, default set true").Envar("CHECK_DNS").Bool()
	CheckNameDns   = kingpin.Flag("check-dns-name", "checking the DNS name resolution, default set ya.ru").Envar("CHECK_DNS_NAME").Default("ya.ru.").String()
	CountRotate    = kingpin.Flag("count-rotate", "Delete old files that are older than the specified number, default set 10").Envar("COUNT_ROTATE").Default("10").Int()
	CountRotate_S3 = kingpin.Flag("count-rotate-s3", "Delete old files in s3 that are older than the specified number, default set 60").Envar("COUNT_ROTATE_S3").Default("60").Int()
	ReportDir      = kingpin.Flag("report-dir", "Path to the save report directory").Envar("REPORT_DIR").Default("/var/log").String()
	ProxyS3Host    = kingpin.Flag("proxyS3Host", "Set s3 proxy host, if you use s3 proxy").Envar("PROXY_S3_HOST").Default("").String()
	AtopReport     = kingpin.Flag("atop-report", "include in the report from the information atop").Envar("ATOP_REPORT").Bool()
)

func init() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}

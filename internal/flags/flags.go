package flags

import (
	"local_trableshoot/configs"

	"github.com/alecthomas/kingpin/v2"
)

var (
	version       = configs.Version
	ContainerFlag = kingpin.Flag("container", "Specify container runtime (e.g. docker)").Envar("container").Default("").String()
	CheckDns      = kingpin.Flag("check-dns", "Tracing to DNS specified in /etc/resolv.conf, default set true").Envar("CHECK_DNS").Bool()
	CountRotate   = kingpin.Flag("count-rotate", "Delete old files that are older than the specified number, default set 10").Envar("COUNT_ROTATE").Default("10").Int()
	ReportDir     = kingpin.Flag("report-dir", "Path to the save report directory").Envar("REPORT_DIR").Default("/var/log").String()
)

func init() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}

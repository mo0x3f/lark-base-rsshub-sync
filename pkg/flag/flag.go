package flag

import (
	"log"
	"strconv"
)

var pageMonitor = 0

func SetPageMonitor(value string) {
	if value == "" {
		return
	}

	pageSize, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid page monitor: %s", value)
	}
	pageMonitor = pageSize
	log.Printf("set page mintior: %d\n", pageMonitor)
}

func PageMonitorEnable() bool {
	return pageMonitor != 0
}

func MonitorPageSize() int {
	return pageMonitor
}

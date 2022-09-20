package lib

import "FundingFeeCashout/db"

type DEBUG struct {
	VERBOSE string
	WARNING string
	INFO    string
}

var DebugLevel *DEBUG

/*
GetDebugLogLevel

获取日志级别
*/
func GetDebugLogLevel() *DEBUG {
	var info, verbose, warning string
	if db.Conf.DEBUG == "verbose" {
		info = "true"
		warning = "true"
		verbose = "true"
	} else if db.Conf.DEBUG == "warning" {
		warning = "true"
		info = "true"
	} else if db.Conf.DEBUG == "info" {
		info = "true"
	}

	return &DEBUG{
		VERBOSE: verbose,
		WARNING: warning,
		INFO:    info,
	}
}

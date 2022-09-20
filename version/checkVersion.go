package version

import (
	"FundingFeeCashout/appConf"
	"fmt"
	"net/http"
	"os"
	"srkTools"
)

const confUrl = "https://agit.ai/srank/srkConf/raw/branch/master/conf"
const logPrefix = "【版本检测】"
const upgradeUrl = "https://github.com/srankXY/FundingFeeCashout"

type versionConf struct {
	UPDATE  bool
	VERSION string
}

func CheckVersion() {

	var result versionConf

	resp, err := http.Get(confUrl)
	if err != nil {
		fmt.Println(logPrefix + "版本检查失败，请重新运行软件")
		os.Exit(1)
	}

	srkTools.DecodeJson(logPrefix, "", resp, &result)

	if result.UPDATE && result.VERSION != appConf.CurrentVersion {
		fmt.Println(logPrefix + "检测到有新版本，请升级之后运行，下载地址: " + upgradeUrl)
		os.Exit(1)
	}
}

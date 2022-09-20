package main

import (
	"FundingFeeCashout/app"
	"FundingFeeCashout/appConf"
	"FundingFeeCashout/db"
	"FundingFeeCashout/exchangeBase/ftx"
	"FundingFeeCashout/exchangeBase/gate"
	"FundingFeeCashout/exchangeBase/okx"
	"FundingFeeCashout/lib"
	"FundingFeeCashout/version"
	"fmt"
	"os"
	"sync"
)

//var Ex = map[string]interface{}{
//	"okx":  okx.NewClient,
//	"ftx":  ftx.Client(),
//	"gate": gate.Client(),
//}

func main() {

	// 版本检测
	version.CheckVersion()

	// cmd 参数判断
	if len(os.Args) > 1 {

		// 初始化数据库
		if os.Args[1] == "init" {
			db.InitDB()

			//	查看版本号
		} else if os.Args[1] == "version" {
			fmt.Println(appConf.CurrentVersion)
			os.Exit(0)
		}

	}

	// 判断db配置文件是否存在
	_, err := os.Stat(db.DbPath)
	if os.IsNotExist(err) {
		fmt.Printf("你应该先运行 %s init， 进行初始化\n", os.Args[0])
		os.Exit(1)
	}

	// 获取数据库配置
	db.Conf = db.QueryConf()
	// 获取日志级别
	lib.DebugLevel = lib.GetDebugLogLevel()

	// 获取手动传入的基础信息, moreExchange： 开多/平空. lowExchange： 开空/平多
	totalCoinCount, method, moreExchange, lowExchange, coin := app.GetBaseInfo()

	// 同步执行多空操作
	// 初始化进程同步通信通道
	lib.SyncStats = ""
	lib.PlaceOrderSyncSlice = []string{}

	// 加锁等待两个交易所执行完成
	var wg sync.WaitGroup
	wg.Add(2)
	if method == "open" {
		// 开多
		switch moreExchange {
		case "okx":
			go okx.OpenPosition("buy", coin, totalCoinCount, &wg)
		case "ftx":
			go ftx.OpenPosition("buy", coin, totalCoinCount, &wg)
		case "gate":
			go gate.OpenPosition("buy", coin, totalCoinCount, &wg)
		}

		// 开空
		switch lowExchange {
		case "okx":
			go okx.OpenPosition("sell", coin, totalCoinCount, &wg)
		case "ftx":
			go ftx.OpenPosition("sell", coin, totalCoinCount, &wg)
		case "gate":
			go gate.OpenPosition("sell", coin, totalCoinCount, &wg)
		}

	} else if method == "close" {
		// 平空
		switch moreExchange {
		case "okx":
			go okx.ClosePosition("buy", coin, &wg)
		case "ftx":
			go ftx.ClosePosition("buy", coin, &wg)
		case "gate":
			go gate.ClosePosition("buy", coin, &wg)
		}

		// 平多
		switch lowExchange {
		case "okx":
			go okx.ClosePosition("sell", coin, &wg)
		case "ftx":
			go ftx.ClosePosition("sell", coin, &wg)
		case "gate":
			go gate.ClosePosition("sell", coin, &wg)
		}

	}
	wg.Wait()
}

package app

import (
	"FundingFeeCashout/exchangeBase/ftx"
	"FundingFeeCashout/exchangeBase/gate"
	"FundingFeeCashout/exchangeBase/okx"
	"FundingFeeCashout/lib"
	"fmt"
	"os"
	"srkTools"
	"strconv"
)

const initTotalCount = 9999999999

/*
GetBaseInfo

获取用户输入的初始信息，根据初始信息判断开平仓方向，计算开仓数量

使用该方法，应先确定返回值中的method具体是开仓还是平仓，进而进行判断操作
*/

func GetBaseInfo() (int, string, string, string, string) {
	// 获取初始的基础信息
	// moreExchange： 开多/平空
	// lowExchange： 开空/平多
	var (
		coin, moreExchange, lowExchange, method                       string
		okxCount, ftxCount, gateCount, minCoinCount, intTotalCoinCout int
	)

	// 获取操作币种
	fmt.Printf("请确认要操作的币种\n>")
	_, _ = fmt.Scanln(&coin)

	// 获取操作模式
	fmt.Printf("请确认开仓还是平仓， 开仓：open，平仓：close\n>")
	_, _ = fmt.Scanln(&method)

	// 做多的情况
	if method == "open" {

		fmt.Printf("请输入做多的交易所，交易所名称示例：【okx，ftx，gate】\n>")
		_, _ = fmt.Scanln(&moreExchange)

		fmt.Printf("请输入做空的交易所，交易所名称示例：【okx，ftx，gate】\n>")
		_, _ = fmt.Scanln(&lowExchange)

		// 获取对应交易所可开仓数量
		if moreExchange == "okx" || lowExchange == "okx" {
			okxCount = int(okx.CalcOpenCount(coin))
		}
		if moreExchange == "ftx" || lowExchange == "ftx" {
			ftxCount = int(ftx.CalcOpenCount(coin))
		}
		if moreExchange == "gate" || lowExchange == "gate" {
			gateCount = int(gate.CalcOpenCount(coin))
		}
		// 获取初始的最小的开仓数量
		arr := []int{okxCount, ftxCount, gateCount}
		minCoinCount = initTotalCount

		for _, i := range arr {
			if i != 0 && minCoinCount > i {
				minCoinCount = i
			}
		}
		fmt.Printf("【ALL】使用最小规则选出的初始可开仓的币数量为： %d %s \n", minCoinCount, coin)
		intTotalCoinCout = CalcCount(coin, moreExchange, lowExchange, minCoinCount)

		// 做空的情况
	} else if method == "close" {
		fmt.Printf("请输入平空的交易所，交易所名称示例：【okx，ftx，gate】\n>")
		_, _ = fmt.Scanln(&moreExchange)

		fmt.Printf("请输入平多的交易所，交易所名称示例：【okx，ftx，gate】\n>")
		_, _ = fmt.Scanln(&lowExchange)
	}

	return intTotalCoinCout, method, moreExchange, lowExchange, coin
}

/*
CalcCount

统一各交易所币量，这里原本该有的小数位会全部被int取整省略，以保证最低合约张数一定能成功进仓

根据各个交易所合约币种，张数<->币量，的转换系数，计算最终适用于所有交易所下单的币量
*/

func CalcCount(coin, moreExchange, lowExchange string, minCoinCount int) int {

	var (
		gateContractSize, okxContractSize float64
		gateFloatPerContractCoinCount     float64
	)

	if moreExchange == "gate" || lowExchange == "gate" {
		// 币转张
		// 获取GATE交易所一张合约对应的币的数量
		gatePerContractCoinCount := gate.GetFutureInfo(gate.Settle, coin+gate.Prefix).QUANTO_MULTIPLIER
		gateFloatPerContractCoinCount, _ = strconv.ParseFloat(gatePerContractCoinCount, 32)
		srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】单张 %s 合约对应的 %s 币量为 %s 个", coin, coin, gatePerContractCoinCount))

		// 计算GATE总共能买多少张合约，这里单张合约可能可购买的币量为浮点数，所以必须用浮点数存储
		gateContractSize = float64(minCoinCount) / gateFloatPerContractCoinCount

		if gateContractSize < 1 {
			fmt.Println("【GATE】资金不足，可开仓量不足1张")
			os.Exit(1)
		}
		srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】%d %s 币量可开合约张数为 %f 张", minCoinCount, coin, gateContractSize))
	}

	if moreExchange == "okx" || lowExchange == "okx" {

		// 获取OKX交易所一张合约对应的币的数量
		// 张转币
		okxPerContractCoinCount := okx.CoinToSheet(coin+okx.Prefix, "1", "2").DATA[0].SZ
		okxFloatPerContractCoinCount, _ := strconv.ParseFloat(okxPerContractCoinCount, 32)
		srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【OKX】1张 %s 合约对应的 %s 币量为 %s 个", coin, coin, okxPerContractCoinCount))

		// 币转张
		//okxResult := okx.CoinToSheet(coin+okx.Prefix, strconv.Itoa(minCoinCount), "1").DATA[0].SZ
		// 单张合约可能可购买的币量为小数，必须用浮点数存储
		okxContractSize = float64(minCoinCount) / okxFloatPerContractCoinCount

		if okxContractSize < 1 {
			fmt.Println("【OKX】资金不足，可开仓量不足1张")
			os.Exit(1)
		}
		srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【OKX】%d %s 币量可开合约张数为 %f 张", minCoinCount, coin, okxContractSize))
	}
	// 对比合约张数大小，选择小的交易所计算最终数量（张数小表示面额大，面额大表示容错率低，所以选择小的）

	exNameMap := map[string]int{
		// 计算出来的合约张数实际上经过int转换之后不会有任何变化
		"okx":  int(okxContractSize),
		"gate": int(gateContractSize),
	}

	minContractSize := initTotalCount
	minContractEx := ""
	for k, v := range exNameMap {

		// 取最小合约张数
		if v != 0 && minContractSize > v {
			minContractSize = v
			minContractEx = k
		}
	}
	fmt.Printf("【ALL】采用 %s 交易所进行币量转换\n", minContractEx)
	fmt.Printf("【ALL】实际采用可开仓 %s 张数为： %d 张\n", coin, minContractSize)

	// 已获取最小张数的交易所名称，使用对应交易所的方法计算最终数量
	var intTotalCoinCount int

	if minContractEx != "" {

		if minContractEx == "okx" {
			// 张转币，这里okx的合约张数一定是整数，所以币量一定是整数
			totalCoinCount := okx.CoinToSheet(coin+okx.Prefix, strconv.Itoa(minContractSize), "2").DATA[0].SZ

			srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【OKX】%d 张 %s 合约， 转换为币量为 %s", minContractSize, coin, totalCoinCount))

			intTotalCoinCount, _ = strconv.Atoi(totalCoinCount)
		} else if minContractEx == "gate" {
			// 张转币，这里gate的合约张数本身是整数，但单张合约可购买的币量为浮点数，所以使用浮点数计算，用int转换为整数
			intTotalCoinCount = int(gateFloatPerContractCoinCount * float64(minContractSize))
			srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】%d 张 %s 合约， 转换为币量为 %d\n", minContractSize, coin, intTotalCoinCount))
		}
	} else {
		// 如果选择的交易所都不用合约张数转换，那么直接返回初始的币量
		intTotalCoinCount = minCoinCount
	}

	fmt.Printf("【ALL】已转换最大开仓的币数量为: %d %s\n", intTotalCoinCount, coin)
	return intTotalCoinCount
}

package okx

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"fmt"
	"math"
	"os"
	"srkTools"
	"strconv"
	"sync"
	"time"
)

const Prefix = "-USDT-SWAP"
const ExchangeName = "OKX"

var direction = map[string]string{
	"buy":  "多仓",
	"sell": "空仓",
}

/*
OpenPosition
开仓函数，开多开空均可

side: buy/sell

coin: USDT
*/
func OpenPosition(side, coin string, totalCoinCount int, wg *sync.WaitGroup) {

	defer wg.Done()
	openDirection := map[string]string{
		"buy":  "开多",
		"sell": "开空",
	}

	fmt.Printf("【OKX】正在进行[%s]操作\n", openDirection[side])

	// 币转张
	totalContractCount := CoinToSheet(coin+Prefix, strconv.Itoa(totalCoinCount), "1").DATA[0].SZ
	totalContractSize, _ := strconv.ParseFloat(totalContractCount, 32)

	fmt.Printf("【OKX】最大开仓币量：%d \n", totalCoinCount)

	// 开仓
	operatePosition(side, coin, int(totalContractSize), openDirection)

}

func ClosePosition(side, coin string, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		posiType          string
		totalContractSize int
		//coin       string
	)

	closeDirection := map[string]string{
		"buy":  "平空",
		"sell": "平多",
	}

	fmt.Printf("【OKX】正在进行[%s]操作\n", closeDirection[side])

	// 传入合约名称获取可操作持仓总量
	perPosiInfo := GetPositions("SWAP", coin+Prefix).DATA[0]
	totalContractSize, _ = strconv.Atoi(perPosiInfo.POS)

	if totalContractSize > 0 {
		posiType = "buy"
		side = "sell"
	} else {
		posiType = "sell"
		side = "buy"
	}

	fmt.Printf("【OKX】已获取 %s 的 %s 仓位信息，可操作持仓量[%d张]\n", coin, direction[posiType], totalContractSize)

	//// 自动获取持仓信息
	//totalPositions, _ := GetPosINFO("")
	//for _, i := range totalPositions.RESULT {
	//	if i.OPENSIZE != 0 {
	//		totalCount = i.OPENSIZE
	//		coin = i.FUTURE
	//		fmt.Printf("【OKX】检测到 %s 的 %s 仓位信息，可操作持仓量[%f]\n", coin, direction[i.SIDE], totalCount)
	//	}
	//}

	// 平仓
	operatePosition(side, coin, totalContractSize, closeDirection)

}

func operatePosition(side, coin string, totalContractSize int, operateDirection map[string]string) {

	// OKX开仓数量不存在正负问题
	fmt.Printf("【OKX】 最大开仓张数：%d 张\n", totalContractSize)
	fmt.Printf("【OKX】 使用 %d 次拆单进行下单\n", db.Conf.SPILT_COUNT)
	// 计算单次下单数量,省略小数位
	perPlaceCount := int(math.Abs(float64(totalContractSize / db.Conf.SPILT_COUNT)))
	if perPlaceCount < 1 {
		fmt.Println("【OKX】单次下单数量小于1，请检查资金是否不足，杠杆倍数是否合理，拆单次数是否过大")
		os.Exit(1)
	}

	fmt.Printf("【OKX】 单次下单张数：%d 张\n", perPlaceCount)

	// 仓位操作
	for i := 0; i < int(math.Abs(float64(totalContractSize))); i += perPlaceCount {

		// 如果当前单次下单量大于剩余单量，则设置下单量为剩余单量
		if int(math.Abs(float64(totalContractSize)))-i < perPlaceCount {
			perPlaceCount = totalContractSize - i
		}

		status, coinPrice := PostOrder(coin+Prefix, "isolated", side, "limit", strconv.Itoa(perPlaceCount), GetTicker(coin + Prefix).DATA[0].LAST)
		if status.CODE != "0" {
			fmt.Println("【OKX】" + status.MSG)
			os.Exit(1)
		}
		fmt.Printf("【OKX】已下单， 本次下单量为 %d 张， 下单价格: %f \n", perPlaceCount, coinPrice)

		// 判断未成交委托
		waitedTime := 0
		pendLoopCount := 0
		for pendingOrder := GetPendOrder(coin+Prefix, "limit"); pendingOrder.CODE != "0" || len(pendingOrder.DATA) != 0; {
			fmt.Printf("【OKX】[%s]订单未完全成交, 将等待%d秒 \n", operateDirection[side], db.Conf.PEND_TIMEOUT-waitedTime)
			time.Sleep(db.DefaultPendLoopWait * time.Second)
			waitedTime += db.DefaultPendLoopWait

			pendingOrder = GetPendOrder(coin+Prefix, "limit")
			// 当订单超过db.DefaultPendLoopWait分钟仍未成交，则重置订单，重置计时器
			if waitedTime >= db.Conf.PEND_TIMEOUT {
				srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【OKX】计时器已超时，将修改订单价格为新价格"))
				if len(pendingOrder.DATA) != 0 {
					orderId := pendingOrder.DATA[0].ORDID
					pendSide := pendingOrder.DATA[0].SIDE
					_, coinPrice = ModifyOrder(pendLoopCount, coin+Prefix, orderId, GetTicker(coin + Prefix).DATA[0].LAST, pendSide)
					srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【OKX】订单修改完成，id为：%s", orderId))
					waitedTime = 0
					pendLoopCount += 1
				}
			}
		}
		fmt.Printf("【OKX】一笔[%s]订单已成交，成交量[%d张%s]，当前成交价格：%f usdt，当前成交总量%s\n", operateDirection[side], perPlaceCount, coin+Prefix, coinPrice, CoinToSheet(coin+Prefix, strconv.Itoa(i+perPlaceCount), "2").DATA[0].SZ)
		lib.SyncProcess(&lib.PlaceOrderSyncSlice, ExchangeName)
	}
	lib.SyncStats = "ok"
}

func CalcOpenCount(coin string) float64 {

	// 修改合约倍数
	changeInfo := ChangeLeverage(strconv.Itoa(db.Conf.LEVERAGE), coin+Prefix)
	if changeInfo.CODE != "0" {
		fmt.Println("【OKX】修改合约倍数失败，请检查账户信息，如持仓模式是否为单向持仓模式")
		os.Exit(1)
	}

	// 获取账户可操作余额
	availableBalance := GetBalance("USDT").DATA[0].DETAILS[0].AVAILEQ
	if availableBalance == "" {
		fmt.Println("【OKX】获取账户可操作余额失败，请检查账户信息，如资金是否在交易账户，合约账户等")
		os.Exit(1)
	}

	floatAvailableBalance, _ := strconv.ParseFloat(availableBalance, 64)
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【OKX】可操作余额为 %f %s ", floatAvailableBalance, "USDT"))

	// 获取操作合约当前价格
	coinPrice := GetTicker(coin + Prefix).DATA[0].LAST
	floatCoinPrice, _ := strconv.ParseFloat(coinPrice, 64)
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【OKX】%s 当前价格为 %f %s ", coin, floatCoinPrice, "USDT"))

	// 计算可开仓数量
	// 计算规则: 总余额的可操作比例 / 对应合约币种上调比例之后的价格
	availableOpenCount := (float64(db.Conf.LEVERAGE) * (floatAvailableBalance * db.Conf.BALANCE_USED_RATIO)) / (floatCoinPrice / db.Conf.PRICE_RATIO)
	fmt.Printf("【OKX】%s可开仓币量为： %f %s \n", coin, availableOpenCount, "个")

	return availableOpenCount
}

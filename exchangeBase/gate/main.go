package gate

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

const Prefix = "_USDT"
const Settle = "usdt"
const ExchangeName = "GATE"

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

	fmt.Printf("【GATE】正在进行[%s]操作\n", openDirection[side])

	// 如果side=sell（开空）则把totalCount转换为负数
	if side == "sell" {
		totalCoinCount = totalCoinCount * -1
	}

	// 获取一张合约对应的币的数量， 可能是浮点数
	perContractCoinCount := GetFutureInfo(Settle, coin+Prefix).QUANTO_MULTIPLIER
	floatPerContractCoinCount, _ := strconv.ParseFloat(perContractCoinCount, 32)

	// 计算总共能买多少张合约
	totalContractSize := int(float64(totalCoinCount) / floatPerContractCoinCount)
	fmt.Printf("【GATE】最大开仓币量：%d \n", totalCoinCount)

	// 开仓
	operatePosition(side, coin, totalContractSize, openDirection)

}

func ClosePosition(side, coin string, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		posiType           string
		totalContractCount int
		//coin       string
	)

	closeDirection := map[string]string{
		"buy":  "平空",
		"sell": "平多",
	}

	fmt.Printf("【GATE】正在进行[%s]操作\n", closeDirection[side])

	// 传入合约名称获取可操作持仓总量
	perPosiInfo := GetPosition(Settle, coin+Prefix)
	// 因为是平仓，所以仓位需要反向操作，比如多仓为正数，实际下单就应为负数
	totalContractCount = perPosiInfo.SIZE * -1

	// 判断持仓信息为多仓，还是空仓
	if totalContractCount > 0 {
		posiType = "buy"
	} else {
		posiType = "sell"
	}
	fmt.Printf("【GATE】已获取 %s 的 %s 仓位信息，可操作持仓量[%d张]\n", coin, direction[posiType], totalContractCount)

	//// 自动获取持仓信息
	//totalPositions, _ := GetPosINFO("")
	//for _, i := range totalPositions.RESULT {
	//	if i.OPENSIZE != 0 {
	//		totalCount = i.OPENSIZE
	//		coin = i.FUTURE
	//		fmt.Printf("【GATE】检测到 %s 的 %s 仓位信息，可操作持仓量[%f]\n", coin, direction[i.SIDE], totalCount)
	//	}
	//}

	// 平仓
	operatePosition(side, coin, totalContractCount, closeDirection)

}

func operatePosition(side, coin string, totalContractSize int, operateDirection map[string]string) {

	// totalContractSize 可能是负数，负数表示持有空仓
	fmt.Printf("【GATE】 最大开仓张数：%d 张\n", totalContractSize)
	fmt.Printf("【GATE】 使用 %d 次拆单进行下单\n", db.Conf.SPILT_COUNT)

	// 获取一张合约对应的币的数量，可能是浮点数
	perContractCoinCount := GetFutureInfo(Settle, coin+Prefix).QUANTO_MULTIPLIER
	floatPerContractCoinCount, _ := strconv.ParseFloat(perContractCoinCount, 32)
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】1张 %s 合约，包含的币数量为：%f", coin, floatPerContractCoinCount))

	// 计算单次下单数量, 可能是浮点数和负数
	perPlaceCount := totalContractSize / db.Conf.SPILT_COUNT
	if math.Abs(float64(perPlaceCount)) < 1 {
		fmt.Println("【GATE】单次下单数量小于1，请检查资金是否不足，杠杆倍数是否合理，拆单次数是否过大")
		os.Exit(1)
	}
	fmt.Printf("【GATE】 单次下单张数：%d 张\n", perPlaceCount)

	// 仓位操作
	for i := 0.0; i < math.Abs(float64(totalContractSize)); i += math.Abs(float64(perPlaceCount)) {

		// 如果当前下单量大于剩余单量，则设置下单量为剩余单量
		if math.Abs(float64(totalContractSize))-i < math.Abs(float64(perPlaceCount)) {
			// 如果总量是负数（开空、平多），则把当前下单量重置为负数
			if totalContractSize < 0 {
				perPlaceCount = int(math.Abs(float64(totalContractSize))-i) * -1
			} else {
				perPlaceCount = int(math.Abs(float64(totalContractSize)) - i)
			}
		}

		status, coinPrice := PlaceOrder(Settle, GetFutureInfo(Settle, coin+Prefix).LAST_PRICE, coin+Prefix, perPlaceCount)
		if status.LABEL != "" {
			fmt.Println("【GATE】" + status.MESSAGE)
			os.Exit(1)
		}
		fmt.Printf("【GATE】已下单， 本次下单量为 %d 张， 下单价格: %f \n", perPlaceCount, coinPrice)

		// 判断未成交委托
		waitedTime := 0
		pendLoopCount := 0
		for code, pendingOrder := GetPendOrder(Settle, coin+Prefix); code != 200 || len(pendingOrder) != 0; {
			fmt.Printf("【GATE】[%s]订单未完全成交, 将等待%d秒 \n", operateDirection[side], db.Conf.PEND_TIMEOUT-waitedTime)
			time.Sleep(db.DefaultPendLoopWait * time.Second)
			waitedTime += db.DefaultPendLoopWait

			code, pendingOrder = GetPendOrder(Settle, coin+Prefix)
			// 当订单超过db.PendingWait分钟仍未成交，则重置订单，重置计时器
			if waitedTime >= db.Conf.PEND_TIMEOUT {
				srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【GATE】计时器已超时，将修改订单价格为新价格"))
				if len(pendingOrder) != 0 {
					orderId := pendingOrder[0].ID
					pendSize := pendingOrder[0].LEFT
					_, coinPrice = ModifyOrder(pendLoopCount, Settle, coin+Prefix, strconv.Itoa(orderId), GetFutureInfo(Settle, coin+Prefix).LAST_PRICE, pendSize)
					srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【GATE】订单修改完成，id为：%d", orderId))
					waitedTime = 0
					pendLoopCount += 1
				}
			}
		}
		fmt.Printf("【GATE】一笔[%s]订单已成交，成交量[%d张%s]，成交价格：%f usdt, 当前成交总量%f\n", operateDirection[side], perPlaceCount, coin+Prefix, coinPrice, (i+float64(perPlaceCount))*floatPerContractCoinCount)
		lib.SyncProcess(&lib.PlaceOrderSyncSlice, ExchangeName)
	}
	lib.SyncStats = "ok"
}

func CalcOpenCount(coin string) float64 {
	// 修改合约倍数
	changeInfo := UpdateLeverage(Settle, coin+Prefix, strconv.Itoa(db.Conf.LEVERAGE), "")
	if changeInfo.LEVERAGE != strconv.Itoa(db.Conf.LEVERAGE) {
		fmt.Println("【GATE】修改合约倍数失败，请检查账户信息，如持仓模式是否为单向持仓模式")
		os.Exit(1)
	}

	// 获取账户可操作余额
	strAvailableBalance := GetFutureAccountInfo(Settle).AVAILABLE
	if strAvailableBalance == "" {
		fmt.Println("【GATE】获取账户可操作余额失败，请检查账户信息，如资金是否在交易账户，合约账户等")
		os.Exit(1)
	}
	availableBalance, _ := strconv.ParseFloat(strAvailableBalance, 64)
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】可操作余额为 %f %s ", availableBalance, "USDT"))

	// 获取操作合约当前价格
	coinPrice, _ := strconv.ParseFloat(GetFutureInfo(Settle, coin+Prefix).LAST_PRICE, 64)
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【GATE】%s 当前价格为 %f %s ", coin, coinPrice, "USDT"))

	// 计算可开仓数量
	// 计算规则: 总余额的可操作比例 / 对应合约币种上调比例之后的价格
	availableOpenCount := (float64(db.Conf.LEVERAGE) * (availableBalance * db.Conf.BALANCE_USED_RATIO)) / (coinPrice / db.Conf.PRICE_RATIO)
	fmt.Printf("【GATE】%s可开仓币量为： %f %s \n", coin, availableOpenCount, "个")

	return availableOpenCount
}

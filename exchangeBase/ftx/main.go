package ftx

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"fmt"
	"os"
	"srkTools"
	"sync"
	"time"
)

const Prefix = "-PERP"
const ExchangeName = "FTX"

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

	fmt.Printf("【FTX】正在进行[%s]操作\n", openDirection[side])

	fmt.Printf("【FTX】最大开仓币量：%d \n", totalCoinCount)
	// 开仓
	operatePosition(side, coin, totalCoinCount, openDirection)

}

func ClosePosition(side, coin string, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		totalCoinCount int
		//coin       string
	)

	closeDirection := map[string]string{
		"buy":  "平空",
		"sell": "平多",
	}

	fmt.Printf("【FTX】正在进行[%s]操作\n", closeDirection[side])

	// 传入合约名称获取可操作持仓总量
	_, perPosiInfo := GetPosINFO(coin + Prefix)
	totalCoinCount = int(perPosiInfo.OPENSIZE)

	// 根据仓位方向，再次确认平仓方向
	if perPosiInfo.SIDE == "buy" {
		side = "sell"
	} else if perPosiInfo.SIDE == "buy" {
		side = "buy"
	}

	fmt.Printf("【FTX】已获取 %s 的 %s 仓位信息，可操作持仓量[%d]\n", coin, direction[perPosiInfo.SIDE], totalCoinCount)

	//// 自动获取持仓信息
	//totalPositions, _ := GetPosINFO("")
	//for _, i := range totalPositions.RESULT {
	//	if i.OPENSIZE != 0 {
	//		totalCount = i.OPENSIZE
	//		coin = i.FUTURE
	//		fmt.Printf("【FTX】检测到 %s 的 %s 仓位信息，可操作持仓量[%f]\n", coin, direction[i.SIDE], totalCount)
	//	}
	//}

	// 平仓
	operatePosition(side, coin, totalCoinCount, closeDirection)

}

func operatePosition(side, coin string, totalCoinCount int, operateDirection map[string]string) {
	// FTX 开仓数量不存在正负问题

	fmt.Printf("【FTX】 使用 %d 次拆单进行下单\n", db.Conf.SPILT_COUNT)
	// 计算单次下单数量，省略小数位
	perPlaceCount := totalCoinCount / db.Conf.SPILT_COUNT
	if perPlaceCount < 1 {
		fmt.Println("【FTX】单次下单数量小于1，请检查资金是否不足，杠杆倍数是否合理，拆单次数是否过大")
		os.Exit(1)
	}
	fmt.Printf("【FTX】 单次下单币数量：%d 张\n", perPlaceCount)

	/*
		由于goroutine 有可能退出的原因，如果需要接着运行，可考虑在这个位置（单次开单数量计算出来之后），获取当前仓位的持仓数量+未成交挂单的挂单数量
		用新总量（原有的总量-已下单成功的数量=新总量）去做循环下单接着处理

		另外当前算法可能导致最后剩下的合约张数不足以满足最低下单量，如果出现该问题，应考虑slice的方式去存储所有拆单之后的每一个下单量，然后循环该slice即可
		如：总合约张数120张，拆单25次，取余之后，单次单量为4张，余20，则需要两个循环
		循环1，以拆单次数为条件循环创建一个slice 包含 25个4
		循环2，以余数为条件循环，往上一个slice的前 20个下标(可以为i-1,i为当前循环的次数)的对应值+1
		最终slice的值应为：20个5 + 5个4
	*/

	// 仓位操作
	for i := 0; i < totalCoinCount; i += perPlaceCount {

		// 如果当前单次下单量大于剩余单量，则设置下单量为剩余单量
		if totalCoinCount-i < perPlaceCount {
			perPlaceCount = totalCoinCount - i
		}

		status, coinPrice := PlaceOrder(coin+Prefix, side, "limit", GetFuturePrice(coin+Prefix).RESULT.LAST, perPlaceCount)
		if !status.SUCCESS {
			fmt.Println("【FTX】" + status.ERROR)
			os.Exit(1)
		}
		fmt.Printf("【FTX】已下单， 本次下单量为 %d 个， 下单价格: %f \n", perPlaceCount, coinPrice)

		// 判断未成交委托
		waitedTime := 0
		pendLoopCount := 0
		for pendingOrder := GetPendindOrder(coin + Prefix); !pendingOrder.SUCCESS || len(pendingOrder.RESULT) != 0; {
			fmt.Printf("【FTX】[%s]订单未完全成交, 将等待%d秒 \n", operateDirection[side], db.Conf.PEND_TIMEOUT-waitedTime)
			time.Sleep(db.DefaultPendLoopWait * time.Second)
			waitedTime += db.DefaultPendLoopWait

			pendingOrder = GetPendindOrder(coin + Prefix)
			// 当订单超时仍未成交，则重置订单，重置计时器
			if waitedTime >= db.Conf.PEND_TIMEOUT {
				srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【FTX】计时器已超时，将修改订单价格为新价格"))
				if len(pendingOrder.RESULT) != 0 {
					orderId := pendingOrder.RESULT[0].ID
					pendSide := pendingOrder.RESULT[0].SIDE
					_, coinPrice = ModifyOrder(pendLoopCount, orderId, GetFuturePrice(coin).RESULT.LAST, coin+Prefix, pendSide)
					srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【FTX】订单修改完成，id为：%d", orderId))
					waitedTime = 0
					pendLoopCount += 1
				}
			}
		}
		fmt.Printf("【FTX】一笔[%s]订单已成交，成交量[%d%s]，当前成交价格：%f usdt， 当前成交总量%d\n", operateDirection[side], perPlaceCount, coin+Prefix, coinPrice, i+perPlaceCount)
		lib.SyncProcess(&lib.PlaceOrderSyncSlice, ExchangeName)

	}
	lib.SyncStats = "ok"
}

func CalcOpenCount(coin string) float64 {
	// 修改合约倍数
	changeInfo := ChangeLeverage(db.Conf.LEVERAGE)
	if !changeInfo.SUCCESS {
		fmt.Println("【FTX】修改合约倍数失败，请检查账户信息，如持仓模式是否为单向持仓模式")
		os.Exit(1)
	}

	// 获取账户可操作余额
	availableBalance := GetBalance("USDT").FREE
	if availableBalance == 0.0 {
		fmt.Println("【FTX】获取账户可操作余额失败，请检查账户信息，如资金是否在交易账户，合约账户等")
		os.Exit(1)
	}
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【FTX】可操作余额为 %f %s ", availableBalance, "USDT"))

	// 获取操作合约当前价格
	coinPrice := GetFuturePrice(coin + Prefix).RESULT.LAST
	srkTools.DebugLog(lib.DebugLevel.INFO, fmt.Sprintf("【FTX】%s 当前价格为 %f %s ", coin, coinPrice, "USDT"))

	// 计算可开仓数量
	// 计算规则: 总余额的可操作比例 / 对应合约币种上调比例之后的价格
	availableOpenCount := (float64(db.Conf.LEVERAGE) * (availableBalance * db.Conf.BALANCE_USED_RATIO)) / (coinPrice / db.Conf.PRICE_RATIO)
	fmt.Printf("【FTX】%s可开仓币量为： %f %s \n", coin, availableOpenCount, "个")

	return availableOpenCount
}

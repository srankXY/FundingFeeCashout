package main

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"fmt"
	"strings"
	"time"
)

func main() {
	db.Conf = db.QueryConf()
	lib.DebugLevel = lib.GetDebugLogLevel()
	//fmt.Println(srkTools.GetCstTime())
	//
	//ex := gate.GATE{}
	//fmt.Println(ex.GATE())
	//fmt.Println(okx.OKX())
	//fmt.Println(ftx.FTX())
	//coin := "SUSHI_USDT"
	//price := ftx.GetFuturePrice(coin).RESULT.PRICE
	//fmt.Println(price)
	//fmt.Println(ftx.PlaceOrder(coin, "buy","limit", price, 100000))
	//cancelid := ftx.GetPendindOrder(coin).RESULT[0].ID
	//fmt.Println(ftx.ModifyOrder(cancelid, price))
	//fmt.Println(ftx.PlaceOrder(coin, "sell", "limit", price, 100000))
	//price := gate.GetFutureInfo("usdt", coin)
	//count, _ := strconv.Atoi(gate.GetFutureInfo("usdt", coin).QUANTO_MULTIPLIER)
	//size := float32(100000 / count)
	//fmt.Println(gate.GetFutureAccountInfo("usdt").AVAILABLE)
	//price := gate.GetFutureInfo("usdt", coin).MARK_PRICE
	//fmt.Println(gate.UpdateLeverage("usdt", coin, "2", ""))
	//fmt.Println(gate.PlaceOrder("usdt", price, coin, size))
	//fmt.Println(gate.GetPosition("usdt", coin))
	//orderid := gate.GetPendOrder("usdt", coin)[0].ID
	//fmt.Println(gate.GetPendOrder("usdt", coin))
	//fmt.Println(gate.CacelOrder("usdt", orderid))
	//fmt.Println(gate.GetFutureAccountInfo("usdt"))
	//fmt.Println(gate.GetPendOrder("usdt", "SUSHI_USDT"))
	//fmt.Println(okx.ChangeLeverage("2", "SUSHI-USDT-SWAP"))
	//pendingOrder := gate.GetPendOrder("usdt", "SLP_USDT")
	//coinPrice := okx.GetTicker("DOGE-USDT-SWAP").DATA[0].LAST
	//fmt.Println(okx.PostOrder("DOGE-USDT-SWAP", "isolated", "buy", "limit", "1", coinPrice))
	//
	//fundingFee.PullFtxFundingToDB()
	//fundingFee.PullGateFundingToDB()
	//fundingFee.PullOkxFundingToDB()
	//fmt.Println(db.GetMaxCoinAndExchange())
	t1 := time.Now().UnixMilli()
	s := db.SplicingSQL("SELECT name, '%s' AS exName, MAX(%s) as rate FROM fund_fee_rate group by name UNION ", 2)
	fmt.Println(strings.TrimRight(s, "UNION "))
	t2 := time.Now().UnixMilli()

	fmt.Println(t2 - t1)

}

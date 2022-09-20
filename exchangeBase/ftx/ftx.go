package ftx

import (
	"FundingFeeCashout/appConf"
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"srkTools"
	"strconv"
	"time"
)

const logPrefix = "【FTX】"

// ApiRequest func ApiRequest(method string, path string, data []byte, sct interface{}) map[string]interface{} {

/*
	method	"GET POST"

	path	"请求的uri路径，首部不能包含/“

	data	"需要提交的json data数据“

	sct		“用于解析response的 struct结构体”
*/
func ApiRequest(method string, path string, data []byte, sct interface{}) {
	client := Client()
	var resp *http.Response
	//time.Sleep(1 * time.Second)
	srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【FTX】正在请求API: %s", path))

	err := fmt.Errorf("first")
	for err != nil {
		if method == "POST" {
			resp, err = client._post(path, data)
		} else if method == "DELETE" {
			resp, err = client._delete(path, data)
		} else {
			resp, err = client._get(path, data)
		}
	}

	// 如果日志级别不是verbose，则不打印解析json日志
	srkTools.DecodeJson(logPrefix, lib.DebugLevel.VERBOSE, resp, sct)
	/* res := srkTools.DecodeJson(resp, sct)

	result := res["result"].(map[string]interface{})
	return result */
}

/*
GetFuturePrice
查询单个币种的标记价格

ask: 卖一价，bid：买一价
*/
func GetFuturePrice(coin string) *Market {
	var market Market
	ApiRequest("GET", "markets/"+coin, []byte(""), &market)
	return &market
}

/*
GetFundRate
coin string [all|具体币种]

查询所有币种或者单个币种的资金费率
*/
func GetFundRate(coin string) *FundRate {

	var fundrate FundRate

	if coin == "all" {
		ApiRequest("GET", "funding_rates", []byte(""), &fundrate)
	} else {
		// ftx每个小时结算一次资金费，所以查看币种当前周期的资金费率这是3600
		startTime := strconv.FormatInt(time.Now().Unix()-3600, 10)
		stopTime := strconv.FormatInt(time.Now().Unix(), 10)
		ApiRequest("GET", "funding_rates?start_time="+startTime+"&end_time="+stopTime+"&future="+coin, []byte(""), &fundrate)
	}
	return &fundrate
}

/*
GetBalance
获取对应币种的余额情况

coin: USDT
*/
func GetBalance(coin string) *BalanceRsult {
	var balance Balance
	var result BalanceRsult
	ApiRequest("GET", "wallet/balances", []byte(""), &balance)
	for _, i := range balance.RESULT {
		if i.COIN == coin {
			result = i
		}
	}
	return &result
}

/*
PlaceOrder

coin: 币种

side: buy 平空/开多，sell 平多/开空

otype: limit 限价，makert 市价

price: 价格

count: 开仓数量
*/
func PlaceOrder(coin, side, otype string, price float64, count int) (*Order, float64) {
	var order Order

	if price > appConf.MinPrice {
		if side == "buy" {
			price *= db.Conf.PRICE_RATIO
		} else {
			price /= db.Conf.PRICE_RATIO
		}
	}

	//newPrice := fmt.Sprintf("%.10f", price)
	//fmt.Println(newPrice)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"market":     coin,
		"side":       side,
		"price":      price,
		"size":       count,
		"type":       otype,
		"reduceOnly": false,
	})

	ApiRequest("POST", "orders", requestBody, &order)

	return &order, price
}

/*
GetPendindOrder
查询提供的币种的委托信息/挂单情况
*/
func GetPendindOrder(coin string) *PendingOrder {
	var pendingorder PendingOrder
	ApiRequest("GET", "orders?market="+coin, []byte(""), &pendingorder)
	return &pendingorder
}

/*
CancelOrder
根据提供的orderId撤销未执行的挂单
*/
func CancelOrder(orderId int) *DelOrder {
	var delorder DelOrder
	cancelId := strconv.Itoa(orderId)
	ApiRequest("DELETE", "orders/"+cancelId, []byte(""), &delorder)
	return &delorder
}

/*
ModifyOrder
修改未执行的委托单子价格

orderId	需要修改的单子id

price	新价格
*/
func ModifyOrder(pendLoopCount, orderId int, price float64, coin, side string) (*ModifyOrderResult, float64) {
	var modifyorder ModifyOrderResult
	cancelId := strconv.Itoa(orderId)

	// 如果对方已经成交且已经超过MaxPendLoopCount次价格调整，则获取当前币种的卖一/买一价格修改订单
	if len(lib.PlaceOrderSyncSlice) == 1 && pendLoopCount >= appConf.MaxPendLoopCount {
		if side == "buy" {
			price = GetFuturePrice(coin).RESULT.ASK
		} else if side == "sell" {
			price = GetFuturePrice(coin).RESULT.BID
		}
	} else if db.Conf.OPEN_MODIFY_ORDER_PRICE_OFFSET == "true" {
		// 根据数据库配置决定是否调整价格
		if side == "buy" {
			price *= db.Conf.PRICE_RATIO
		} else {
			price /= db.Conf.PRICE_RATIO
		}
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"price": price,
	})

	ApiRequest("POST", "orders/"+cancelId+"/modify", requestBody, &modifyorder)

	return &modifyorder, price

}

/*
ChangeLeverage
设置杠杆倍数
*/
func ChangeLeverage(num int) *LeverageResult {
	var result LeverageResult
	leverage, _ := json.Marshal(map[string]interface{}{
		"leverage": num,
	})
	ApiRequest("POST", "account/leverage", leverage, &result)

	return &result
}

/*
GetPosINFO
获取对应币种的持仓情况，可以根据返回结果的SIZE情况，判断是否有持仓

coin: 可留空，可传具体的合约名称
*/
func GetPosINFO(coin string) (*Positions, *PositionResult) {
	var allPosi Positions
	var result PositionResult
	ApiRequest("GET", "positions", []byte(""), &allPosi)

	if coin != "" {
		for _, i := range allPosi.RESULT {

			if i.FUTURE == coin {
				result = i
			}
		}
	}
	return &allPosi, &result
}

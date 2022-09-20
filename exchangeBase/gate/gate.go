package gate

import (
	"FundingFeeCashout/appConf"
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"srkTools"
	"strconv"
)

const logPrefix = "【GATE】"

func ApiRequest(method string, path string, data []byte, sct interface{}) int {
	client := Client()
	var resp *http.Response
	//time.Sleep(1 * time.Second)
	srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【GATE】正在请求API: %s", path))

	//resp = new(http.Response)
	err := fmt.Errorf("first")
	for err != nil {
		if method == "POST" {
			resp, err = client._post(path, data)
		} else if method == "DELETE" {
			resp, err = client._delete(path, data)
		} else if method == "PUT" {
			resp, err = client._put(path, data)
		} else {
			resp, err = client._get(path, data)
		}
	}
	// 如果日志级别不是verbose，则不打印解析json日志

	srkTools.DecodeJson(logPrefix, lib.DebugLevel.VERBOSE, resp, sct)
	return resp.StatusCode

}

func GetAllFuturesInfo(settle string) []FuturesInfo {
	var result []FuturesInfo
	path := fmt.Sprintf("%s%s%s", "/futures/", settle, "/contracts")
	ApiRequest("GET", path, []byte(""), &result)
	return result
}

func GetFutureInfo(settle, contract string) *FuturesInfo {

	var result FuturesInfo
	path := fmt.Sprintf("%s%s%s%s", "/futures/", settle, "/contracts/", contract)
	ApiRequest("GET", path, []byte(""), &result)
	return &result
}

/*
GetOrderBook

获取某个币种成交深度

ask：卖一价，bid：买一价
*/
func GetOrderBook(coin, settle string) *OrderBookInfo {
	var result OrderBookInfo
	path := fmt.Sprintf("%s%s%s%s", "/futures/", settle, "/order_book", "?contract="+coin)
	ApiRequest("GET", path, []byte(""), &result)
	return &result
}

/*
GetFutureAccountInfo

获取合约账户可操作余额
*/
func GetFutureAccountInfo(settle string) *FutureAccountInfo {
	var result FutureAccountInfo
	ApiRequest("GET", "/futures/"+settle+"/accounts", []byte(""), &result)
	return &result

}

func CacelOrder(settle string, orderid int) *PlaceOrderResult {
	var result PlaceOrderResult
	path := fmt.Sprintf("%s%s%s%s", "/futures/", settle, "/orders/", strconv.Itoa(orderid))
	ApiRequest("DELETE", path, []byte(""), &result)
	return &result
}

/*
GetPendOrder

获取挂单情况

返回值中，left字段表示未成交合约张数
*/
func GetPendOrder(settle, contract string) (int, []OrderList) {

	var result []OrderList
	path := fmt.Sprintf("%s%s%s%s%s", "/futures/", settle, "/orders?contract=", contract, "&status=open")
	code := ApiRequest("GET", path, []byte(""), &result)

	return code, result

}

/*
PlaceOrder

下单接口， size为正为开多/平空，size为负为开空/平多

size 为合约张数，需根据合约详情或者单张合约对应的币数量转换之后进行下单
*/
func PlaceOrder(settle, price, contract string, size int) (*PlaceOrderResult, float64) {
	var result PlaceOrderResult

	newPrice, _ := strconv.ParseFloat(price, 32)
	if newPrice > appConf.MinPrice {
		if size > 0 {
			newPrice *= db.Conf.PRICE_RATIO
		} else {
			newPrice /= db.Conf.PRICE_RATIO
		}
	}
	//fmt.Println(price)
	//fmt.Println(strconv.FormatFloat(newPrice,'f',-1,32))
	requestBdoy, _ := json.Marshal(map[string]interface{}{
		"contract": contract,
		"size":     size,
		"price":    strconv.FormatFloat(newPrice, 'f', -1, 32),
	})
	ApiRequest("POST", "/futures/"+settle+"/orders", requestBdoy, &result)
	return &result, newPrice
}

/*
ModifyOrder

修改未成交订单，使用当前标记价格下单
*/
func ModifyOrder(pendLoopCount int, settle, coin, order_id, price string, size int) (*ModifyOrderResult, float64) {

	var result ModifyOrderResult

	// 如果对方已经成交，则获取当前币种的卖一/买一价格修改订单
	if len(lib.PlaceOrderSyncSlice) == 1 && pendLoopCount >= appConf.MaxPendLoopCount {
		if size > 0 {
			price = GetOrderBook(coin, settle).Asks[0].P
		} else {
			price = GetOrderBook(coin, settle).Bids[0].P
		}
	} else if db.Conf.OPEN_MODIFY_ORDER_PRICE_OFFSET == "true" {

		newPrice, _ := strconv.ParseFloat(price, 32)
		if newPrice > appConf.MinPrice {
			if size > 0 {
				newPrice *= db.Conf.PRICE_RATIO
			} else {
				newPrice /= db.Conf.PRICE_RATIO
			}
		}
		// 对价格类型进行转换
		price = strconv.FormatFloat(newPrice, 'f', -1, 32)
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"price": price,
	})
	ApiRequest("PUT", "/futures/"+settle+"/orders/"+order_id, requestBody, &result)

	// 返回float64类型的价格
	lastPrice, _ := strconv.ParseFloat(price, 32)

	return &result, lastPrice
}

/*
UpdateLeverage

更改合约倍数以及持仓模式（全仓或者逐仓）

leverage: 为0表示全仓模式，不为0为逐仓模式

cross_leverage_limit: leverage为0时生效，为全仓模式下的倍数，同时如leverage不为0，该值只能传空值
*/
func UpdateLeverage(settle, contract, leverage, cross_leverage_limit string) *PositionResult {
	var result PositionResult

	var uri string
	if leverage != "0" {
		uri = ""
	} else {
		uri = "&cross_leverage_limit=" + cross_leverage_limit
	}

	path := fmt.Sprintf("%s%s%s%s%s", "/futures/", settle, "/positions/", contract, "/leverage?leverage="+leverage+uri)
	ApiRequest("POST", path, []byte(""), &result)
	return &result
}

/*
GetPosition

获取单个币种仓位信息，根据返回参数判断是否为持仓还是挂单

返回的SIZE大于0，表示当前存在持仓

返回的pending_order大于0，表示存在挂单
*/
func GetPosition(settle, contract string) *PositionResult {
	var result PositionResult
	path := fmt.Sprintf("%s%s%s%s", "/futures/", settle, "/positions/", contract)
	ApiRequest("GET", path, []byte(""), &result)
	return &result
}

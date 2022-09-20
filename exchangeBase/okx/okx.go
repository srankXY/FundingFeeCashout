package okx

import (
	"FundingFeeCashout/appConf"
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"srkTools"
	"strconv"
	"time"
)

// 重试次数
var Retries int = 50
var Ok OkEx

// 查看账户余额
func GetBalance(ccy string) Balance {
	var bal Balance
	uri := fmt.Sprintf("/api/v5/account/balance?ccy=" + ccy)
	request(uri, http.MethodGet, nil, &bal)
	if !reflect.DeepEqual(bal, Balance{}) {
		return bal
	} else {
		return Balance{}
	}
}

// 查看持仓信息
//
// instType：产品类型 (MARGIN：币币杠杆 SWAP：永续合约 FUTURES：交割合约 OPTION：期权)
//
// instId：产品ID
func GetPositions(instType, instId string) Positions {
	var pos Positions
	uri := fmt.Sprintf("/api/v5/account/positions?instType=%s&instId=%s", instType, instId)
	request(uri, http.MethodGet, nil, &pos)
	if !reflect.DeepEqual(pos, Positions{}) {
		return pos
	} else {
		return Positions{}
	}
}

// 获取标记价格
//
// instType：产品类型 (MARGIN：币币杠杆 SWAP：永续合约 FUTURES：交割合约 OPTION：期权)
//
// instId：产品ID (空字符串表示全查)
//
func GetMarkPrice(instType, instId string) MarkPrice {
	var mark MarkPrice
	uri := fmt.Sprintf("/api/v5/public/mark-price?instType=%s&instId=%s", instType, instId)
	// 不传instId是全查
	if instId == "" {
		uri = fmt.Sprintf("/api/v5/public/mark-price?instType=%s", instType)
	} else {
		uri = fmt.Sprintf("/api/v5/public/mark-price?instType=%s&instId=%s", instType, instId)
	}
	request(uri, http.MethodGet, nil, &mark)
	if !reflect.DeepEqual(mark, MarkPrice{}) {
		return mark
	} else {
		return MarkPrice{}
	}
}

// 获取订单信息
//
// instId：产品ID
//
// ordId：订单ID
func GetOrder(instId, ordId string) OrderGet {
	var ord OrderGet
	uri := fmt.Sprintf("/api/v5/trade/order?instId=%s&ordId=%s", instId, ordId)
	request(uri, http.MethodGet, nil, &ord)
	if !reflect.DeepEqual(ord, OrderGet{}) {
		return ord
	} else {
		return OrderGet{}
	}
}

func GetPendOrder(instId, ordType string) *PendOrder {
	var result PendOrder
	uri := "/api/v5/trade/orders-pending?instId=" + instId + "&ordType=" + ordType
	request(uri, http.MethodGet, nil, &result)
	if !reflect.DeepEqual(result, SetLeverage{}) {
		return &result
	} else {
		return &PendOrder{}
	}
}

func CoinToSheet(instId, sz, transType string) *CoinSheet {
	var result CoinSheet
	uri := "/api/v5/public/convert-contract-coin?instId=" + instId + "&sz=" + sz + "&type=" + transType
	request(uri, http.MethodGet, nil, &result)
	if !reflect.DeepEqual(result, SetLeverage{}) {
		return &result
	} else {
		return &CoinSheet{}
	}
}

func ModifyOrder(pendLoopCount int, instId, ordId, price, side string) (*ModifyOrderRes, float64) {
	var result ModifyOrderRes
	uri := fmt.Sprintf("/api/v5/trade/amend-order")

	// 如果对方已经成交，则获取当前币种的卖一/买一价格修改订单
	if len(lib.PlaceOrderSyncSlice) == 1 && pendLoopCount >= appConf.MaxPendLoopCount {
		if side == "buy" {
			price = GetTicker(instId).DATA[0].ASKPX
		} else if side == "sell" {
			price = GetTicker(instId).DATA[0].BIDPX
		}
	} else if db.Conf.OPEN_MODIFY_ORDER_PRICE_OFFSET == "true" {

		newPrice, _ := strconv.ParseFloat(price, 32)

		if side == "buy" {
			newPrice *= db.Conf.PRICE_RATIO
		} else {
			newPrice /= db.Conf.PRICE_RATIO
		}

		price = strconv.FormatFloat(newPrice, 'f', -1, 32)
	}

	tranBody, _ := json.Marshal(map[string]interface{}{
		"ordId":  ordId,
		"newPx":  price,
		"instId": instId,
	})
	request(uri, http.MethodPost, tranBody, &result)

	// 返回float64类型的价格
	lastPrice, _ := strconv.ParseFloat(price, 32)

	if !reflect.DeepEqual(result, SetLeverage{}) {
		return &result, lastPrice
	} else {
		return &ModifyOrderRes{}, lastPrice
	}
}

// 下单
//
// instId：产品ID
//
// tdMode：交易模式,(保证金模式 isolated：逐仓 ；cross：全仓。非保证金模式 cash：非保证金)
//
// side：订单方向 (buy：买  sell：卖)
//
// posSide：持仓方向 (在双向持仓模式下必填，且仅可选择 long 或 short)
//
// ordType：订单类型 (market：市价单  limit：限价单  post_only：只做maker单  fok：全部成交或立即取消  ioc：立即成交并取消剩余 optimal_limit_ioc：市价委托立即成交并取消剩余（仅适用交割、永续）)
//
// sz：委托数量
//
// px：委托价格 (仅适用于limit、post_only、fok、ioc类型的订单)
func PostOrder(instId, tdMode, side, ordType, sz, px string) (Order, float64) {

	var ord Order
	uri := "/api/v5/trade/order"

	price, _ := strconv.ParseFloat(px, 32)
	if price > appConf.MinPrice {
		if side == "buy" {
			price *= db.Conf.PRICE_RATIO
		} else {
			price /= db.Conf.PRICE_RATIO
		}
	}

	ordBody := OrderBody{
		INSTID:  instId,
		TDMODE:  tdMode,
		SIDE:    side,
		ORDTYPE: ordType,
		SZ:      sz,
		PX:      strconv.FormatFloat(price, 'f', -1, 32),
	}

	reqBody, _ := json.Marshal(ordBody)
	request(uri, http.MethodPost, reqBody, &ord)
	if !reflect.DeepEqual(ord, Order{}) {
		return ord, price
	} else {
		return Order{}, price
	}
}

// 获取交易产品基础信息
//
// instType：产品类型 (SPOT：币币 MARGIN：币币杠杆 SWAP：永续合约 FUTURES：交割合约 OPTION：期权)
func GetInstruments(instType string) Instruments {
	var instr Instruments
	uri := fmt.Sprintf("/api/v5/public/instruments?instType=%s", instType)
	request(uri, http.MethodGet, nil, &instr)
	if !reflect.DeepEqual(instr, Instruments{}) {
		return instr
	} else {
		return Instruments{}
	}
}

// 账户内资金划转
//
// form：转出账户 (6：资金账户 18：交易账户)
//
// to：转入账户 (6：资金账户 18：交易账户)
//
// amt：转账数量
func PostTransfer(form, to, amt string) Transfer {
	var trs Transfer
	uri := "/api/v5/asset/transfer"
	tranBody := TransferBody{
		CCY:  "USDT",
		AMT:  amt,
		FROM: form,
		TO:   to,
		TYPE: "0",
	}
	reqBody, _ := json.Marshal(tranBody)
	request(uri, http.MethodPost, reqBody, &trs)
	if !reflect.DeepEqual(trs, Transfer{}) {
		return trs
	} else {
		return Transfer{}
	}
}

// 获取指数行情
//
// instId：产品ID
func GetIndexTickers(instId string) IndexTickers {
	var indexTickers IndexTickers
	uri := fmt.Sprintf("/api/v5/market/index-tickers?quoteCcy=USDT&instId=%s", instId)
	request(uri, http.MethodGet, nil, &indexTickers)
	// 判断结构体是否不为空结构体
	if !reflect.DeepEqual(indexTickers, IndexTickers{}) {
		return indexTickers
	} else {
		return IndexTickers{}
	}
}

// 获取永续合约当前资金费率
//
// instId：产品ID
func GetFundingRate(instId string) FundingRate {
	var fundingRate FundingRate
	uri := fmt.Sprintf("/api/v5/public/funding-rate?instId=%s", instId)
	request(uri, http.MethodGet, nil, &fundingRate)
	// 判断结构体是否不为空结构体
	if !reflect.DeepEqual(fundingRate, FundingRate{}) {
		return fundingRate
	} else {
		return FundingRate{}
	}

}

// 获取单个产品行情信息
//
// instId：产品ID
func GetTicker(instId string) Tickers {
	var ticker Tickers
	uri := fmt.Sprintf("/api/v5/market/ticker?instId=%s", instId)
	request(uri, http.MethodGet, nil, &ticker)
	//fmt.Printf("%#v\n", ticker)
	// 判断结构体是否不为空结构体
	if !reflect.DeepEqual(ticker, Tickers{}) {
		return ticker
	} else {
		return Tickers{}
	}
}

// 获取所有产品行情信息
//
// instType：产品类型 (SPOT：币币 MARGIN：币币杠杆 SWAP：永续合约 FUTURES：交割合约 OPTION：期权)
func GetTickers(instType string) Tickers {
	var tickers Tickers

	uri := fmt.Sprintf("/api/v5/market/tickers?instType=%s", instType)

	request(uri, http.MethodGet, nil, &tickers)
	// 判断结构体是否不为空结构体
	if !reflect.DeepEqual(tickers, Tickers{}) {
		return tickers
	} else {
		return Tickers{}
	}
}

func ChangeLeverage(lever, instId string) *SetLeverage {
	var result SetLeverage
	uri := fmt.Sprintf("/api/v5/account/set-leverage")
	tranBody, _ := json.Marshal(map[string]interface{}{
		"instId":  instId,
		"lever":   lever,
		"mgnMode": "isolated",
		"posSide": "net",
	})
	request(uri, http.MethodPost, tranBody, &result)
	if !reflect.DeepEqual(result, SetLeverage{}) {
		return &result
	} else {
		return &SetLeverage{}
	}
}

// 设置KEY相关数据
//
// apiKey：APIKey
//
// secretKey：SecretKey
//
// passhrasa：在创建API密钥时指定的Passphrase
func SetKey(apiKey, secretKey, passhrasa string) {
	Ok = OkEx{
		APIKEY:     apiKey,
		SECRETKEY:  secretKey,
		PASSPHRASE: passhrasa,
	}
}

func NewClient() {
	key := db.QueryDB("value", "conf", "name", "OKX_KEY")
	secert := db.QueryDB("value", "conf", "name", "OKX_SECERT")
	pass := db.QueryDB("value", "conf", "name", "OKX_PASS")

	SetKey(key, secert, pass)
}

// 调用请求
//
// uri：api地址
//
// method：http请求类型
//
// body：request结构体
//
// str：response结构体
func request(uri, method string, requestByte []byte, str interface{}) {
	// 初始化okx
	NewClient()
	//time.Sleep(1 * time.Second)
	// 请求前检测是否设置Key
	if reflect.DeepEqual(Ok, OkEx{}) {
		fmt.Println("请先使用okex.SetKey()设置apiKey")
		return
	}

	// 重试3次，否则失败
	if Retries < 0 {
		return
	}

	//实盘API交易地址如下：
	//REST：https://www.okx.com/
	//WebSocket公共频道：wss://ws.okx.com:8443/ws/v5/public
	//WebSocket私有频道：wss://ws.okx.com:8443/ws/v5/private
	//
	//AWS 地址如下：
	//REST：https://aws.okx.com
	//WebSocket公共频道：wss://wsaws.okx.com:8443/ws/v5/public
	//WebSocket私有频道：wss://wsaws.okx.com:8443/ws/v5/private
	requestPath := fmt.Sprintf("https://www.okx.com%s", uri)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	requestBody := bytes.NewReader(requestByte)

	req, err := http.NewRequest(method, requestPath, requestBody)
	if err != nil {
		fmt.Println("构建NewRequest对象失败", err)
		return
	}

	// 日志
	srkTools.DebugLog(lib.DebugLevel.WARNING, fmt.Sprintf("【OKX】当前正在请求的API为：%s", uri))
	var sig string
	sig = fmt.Sprintf("%v%v%v%v", timestamp, method, uri, string(requestByte))
	// hmac-sha256编码
	sha := hmac.New(sha256.New, []byte(Ok.SECRETKEY))
	sha.Write([]byte(sig))
	// base64编码
	sign := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	//gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	// 添加请求头
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-KEY", Ok.APIKEY)
	req.Header.Add("OK-ACCESS-PASSPHRASE", Ok.PASSPHRASE)
	req.Header.Add("OK-ACCESS-SIGN", sign)

	// 获取代理配置
	Proxy := db.QueryDB("value", "conf", "name", "PROXY")

	var t *http.Transport
	if Proxy != "" {
		proxy, _ := url.Parse(Proxy)
		t = &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
			Proxy:           http.ProxyURL(proxy),
		}
	} else {
		t = &http.Transport{}
	}

	// 创建http客户端
	client := &http.Client{
		Transport: t,
		Timeout:   time.Second * 5,
	}
	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("client请求失败", err)
		Retries -= 1
		request(uri, method, requestByte, str)
		return
	}
	// 请求成功，设置重试次数为3
	Retries = 3
	defer func() {
		time.Sleep(1)
		client.CloseIdleConnections()
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	// 解析为[]byte
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resBody解析失败", err)
		return
	}

	// 日志
	srkTools.DebugLog(lib.DebugLevel.VERBOSE, "【OKX】"+string(resBody))
	if err := json.Unmarshal(resBody, str); err != nil {
		fmt.Println("json解析失败", err)
		return
	}

	return
}

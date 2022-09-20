package db

type conf struct {
	LEVERAGE                       int
	SPILT_COUNT                    int
	PRICE_RATIO                    float64
	BALANCE_USED_RATIO             float64
	PEND_TIMEOUT                   int
	DEBUG                          string
	OPEN_MODIFY_ORDER_PRICE_OFFSET string
}

type FundingGap struct {
	Coin string
	Gap  string
}

type CoinExchange struct {
	Coin        string
	MaxExchange exNameRate
	MinExchange exNameRate
	Gap         string
}

// 以列表的形式存储最大资金费差的各交易所的名称和费率
type exNameRate struct {
	name string
	rate float64
}

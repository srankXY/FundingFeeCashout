package ftx

import "time"

type BASE struct {
	SUCCESS bool `json:"success"`
	ERROR   string
}

type Market struct {
	BASE
	RESULT MarketResult
}

type MarketResult struct {
	NAME                  string  `json:"name"`
	BASECURRENCY          string  `json:"baseCurrency"`
	QUOTECURRENCY         string  `json:"quoteCurrency"`
	QUOTEVOLUME24H        float64 `json:"quoteVolume24h"`
	CHANGE1H              float64 `json:"change1h"`
	CHANGE24H             float64 `json:"change24h"`
	CHANGEBOD             float64 `json:"changeBod"`
	HIGHLEVERAGEFEEEXEMPT bool    `json:"highLeverageFeeExempt"`
	MINPROVIDESIZE        float64 `json:"minProvideSize"`
	TYPE                  string  `json:"type"`
	UNDERLYING            string  `json:"underlying"`
	ENABLED               bool    `json:"enabled"`
	ASK                   float64 `json:"ask"`
	BID                   float64 `json:"bid"`
	LAST                  float64 `json:"last"`
	POSTONLY              bool    `json:"postOnly"`
	PRICE                 float64 `json:"price"`
	PRICEINCREMENT        float64 `json:"priceIncrement"`
	SIZEINCREMENT         float64 `json:"sizeIncrement"`
	RESTRICTED            bool    `json:"restricted"`
	VOLUMEUSD24H          float64 `json:"volumeUsd24h"`
	PRICEHIGH24H          float64 `json:"priceHigh24h"`
	PRICELOW24H           bool    `json:"priceLow24h"`
}

type FundRate struct {
	BASE
	RESULT []FundRateResult
}

type FundRateResult struct {
	FUTURE string
	RATE   float64
	TIME   time.Time
}

type Balance struct {
	BASE
	RESULT []BalanceRsult
}

type BalanceRsult struct {
	COIN                   string  `json:"coin"`
	FREE                   float64 `json:"free"`
	SPOTBORROW             float64 `json:"spotBorrow"`
	TOTAL                  float64 `json:"total"`
	USDVALUE               float64 `json:"usdValue"`
	AVAILABLEFORWITHDRAWAL float64 `json:"availableForWithdrawal"`
	AVAILABLEWITHOUTBORROW float64 `json:"availableWithoutBorrow"`
}

type Order struct {
	BASE
	RESULT OrderResult
}

type OrderResult struct {
	CREATEDAT     string  `json:"createdAt"`
	FILLEDSIZE    float64 `json:"filledSize"`
	FUTURE        string  `json:"future"`
	ID            int     `json:"id"`
	MARKET        string  `json:"market"`
	PRICE         float64 `json:"price"`
	REMAININGSIZE float64 `json:"remainingSize"`
	SIDE          string  `json:"side"`
	SIZE          float64 `json:"size"`
	STATUS        string  `json:"status"`
	TYPE          string  `json:"type"`
	REDUCEONLY    bool    `json:"reduceOnly"`
	LIQUIDATION   string  `json:"liquidation"`
	AVGFILLPRICE  float64 `json:"avgFillPrice"`
	IOC           bool    `json:"ioc"`
	POSTONLY      bool    `json:"postOnly"`
	CLIENTID      string  `json:"clientId"`
}

type PendingOrder struct {
	BASE
	RESULT []PendingOrderResult
}

type PendingOrderResult struct {
	CREATEDAT     string  `json:"createdAt"`
	FILLEDSIZE    float64 `json:"filledSize"`
	FUTURE        string  `json:"future"`
	ID            int     `json:"id"`
	MARKET        string  `json:"market"`
	PRICE         float64 `json:"price"`
	LIQUIDATION   bool    `json:"liquidation"`
	AVGFILLPRICE  float64 `json:"avgFillPrice"`
	REMAININGSIZE float64 `json:"remainingSize"`
	SIDE          string  `json:"side"`
	SIZE          float64 `json:"size"`
	STATUS        string  `json:"status"`
	TYPE          string  `json:"type"`
	REDUCEONLY    bool    `json:"reduceOnly"`
	IOC           bool    `json:"ioc"`
	POSTONLY      bool    `json:"postOnly"`
	CLIENTID      string  `json:"clientId"`
}

type DelOrder struct {
	BASE
	RESULT string `json:"result"`
}

type ModifyOrderResult struct {
	BASE
	RESULT PendingOrderResult
}

type LeverageResult struct {
	BASE
	RESULT string
}

type Positions struct {
	BASE
	RESULT []PositionResult
}

type PositionResult struct {
	FUTURE                       string  `json:"future"`
	SIZE                         float64 `json:"size"`
	SIDE                         string  `json:"side"`
	NETSIZE                      float64 `json:"netSize"`
	LONGORDERSIZE                float64 `json:"longOrderSize"`
	SHORTORDERSIZE               float64 `json:"shortOrderSize"`
	COST                         float64 `json:"cost"`
	ENTRYPRICE                   string  `json:"entryPrice"`
	UNREALIZEDPNL                float64 `json:"unrealizedPnl"`
	REALIZEDPNL                  float64 `json:"realizedPnl"`
	INITIALMARGINREQUIREMENT     float64 `json:"initialMarginRequirement"`
	MAINTENANCEMARGINREQUIREMENT float64 `json:"maintenanceMarginRequirement"`
	OPENSIZE                     float64 `json:"openSize"`
	COLLATERALUSED               float64 `json:"collateralUsed"`
	ESTIMATEDLIQUIDATIONPRICE    float64 `json:"estimatedLiquidationPrice"`
}

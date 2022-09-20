package gate

type BASE struct {
	LABEL   string
	MESSAGE string
}

//type FutureInfo struct {
//	NAME string
//	TYPE	string
//	QUANTO_MULTIPLIER	string
//	REF_DISCOUNT_RATE	string
//	ORDER_PRICE_DEVIATE	string
//	MAINTENANCE_RATE	string
//	MARK_TYPE	string
//	LAST_PRICE	string
//	MARK_PRICE	string
//	INDEX_PRICE	string
//	FUNDING_RATE_INDICATIVE	string
//	MARK_PRICE_ROUND	string
//	FUNDING_OFFSET	int
//	IN_DELISTING	bool
//	RISK_LIMIT_BASE	string
//	INTEREST_RATE	string
//	ORDER_PRICE_ROUND	string
//	ORDER_SIZE_MIN	int
//	REF_REBATE_RATE	string
//	FUNDING_INTERVAL	int
//	RISK_LIMIT_STEP	string
//	LEVERAGE_MIN	string
//	LEVERAGE_MAX	string
//	RISK_LIMIT_MAX	string
//	MAKER_FEE_RATE	string
//	TAKER_FEE_RATE	string
//	FUNDING_RATE	string
//	ORDER_SIZE_MAX	int
//	FUNDING_NEXT_APPLY	int
//	SHORT_USERS	int
//	CONFIG_CHANGE_TIME	int
//	TRADE_SIZE	int
//	POSITION_SIZE	int
//	LONG_USERS	int
//	FUNDING_IMPACT_VALUE	string
//	ORDERS_LIMIT	int
//	TRADE_ID	int
//	ORDERBOOK_ID	int
//}

type OrderBookInfo struct {
	Id      int     `json:"id"`
	Current float64 `json:"current"`
	Update  float64 `json:"update"`
	Asks    []struct {
		P string `json:"p"`
		S int    `json:"s"`
	} `json:"asks"`
	Bids []struct {
		P string `json:"p"`
		S int    `json:"s"`
	} `json:"bids"`
}

type FuturesInfo struct {
	NAME                    string
	TYPE                    string
	FUNDING_RATE_INDICATIVE float64
	FUNDING_OFFSET          float64
	INTEREST_RATE           string
	QUANTO_MULTIPLIER       string
	FUNDING_IMPACT_VALUE    string
	LEVERAGE_MIN            string
	LEVERAGE_MAX            string
	MAINTENANCE_RATE        string
	MARK_TYPE               string
	MARK_PRICE              string
	INDEX_PRICE             string
	LAST_PRICE              string
	MAKER_FEE_RATE          string
	TAKER_FEE_RATE          string
	ORDER_PRICE_ROUND       string
	MARK_PRICE_ROUND        string
	FUNDING_RATE            string
	FUNDING_INTERVAL        int
	FUNDING_NEXT_APPLY      float64
	SHORT_USERS             int
	RISK_LIMIT_BASE         string
	RISK_LIMIT_STEP         string
	ENABLE_BONUS            bool
	RISK_LIMIT_MAX          string
	ORDER_SIZE_MIN          int64
	ORDER_SIZE_MAX          int64
	ORDER_PRICE_DEVIATE     string
	REF_DISCOUNT_RATE       string
	REF_REBATE_RATE         string
	ORDERBOOK_ID            int64
	TRADE_ID                int64
	TRADE_SIZE              int64
	POSITION_SIZE           int64
	LONG_USERS              int
	CONFIG_CHANGE_TIME      float64
	IN_DELISTING            bool
	ORDERS_LIMIT            int
}

type FutureAccountInfo struct {
	TOTAL           string
	UNREALISED_PNL  string
	AVAILABLE       string
	ORDER_MARGIN    string
	POSITION_MARGIN string
	POINT           string
	CURRENCY        string
	IN_DUAL_MODE    bool
}

type PositionResult struct {
	VALUE                string
	LEVERAGE             string
	MODE                 string
	REALISED_POINT       string
	CONTRACT             string
	ENTRY_PRICE          string
	MARK_PRICE           string
	HISTORY_POINT        string
	REALISED_PNL         string
	CLOSE_ORDER          string
	SIZE                 int
	CROSS_LEVERAGE_LIMIT string
	PENDING_ORDERS       float64
	ADL_RANKING          int
	MAINTENANCE_RATE     string
	UNREALISED_PNL       string
	USER                 int
	LEVERAGE_MAX         string
	HISTORY_PNL          string
	RISK_LIMIT           string
	MARGIN               string
	LAST_CLOSE_PNL       string
	LIQ_PRICE            string
}

type PlaceOrderResult struct {
	BASE
	ID             int
	CONTRACT       string
	MKFR           string
	TKFR           string
	TIF            string
	IS_REDUCE_ONLY bool
	CREATE_TIME    float64
	FINISH_TIME    float64
	PRICE          string
	SIZE           float64
	REFR           string
	LEFT           float64
	TEXT           string
	FILL_PRICE     string
	USER           int
	FINISH_AS      string
	STATUS         string
	IS_LIQ         bool
	REFU           int
	IS_CLOSE       bool
	ICEBERG        int
}

type OrderList struct {
	STATUS         string
	SIZE           float64
	LEFT           int
	ID             int
	IS_LIQ         bool
	IS_CLOSE       bool
	CONTRACT       string
	TEXT           string
	FILL_PRICE     string
	FINISH_AS      string
	ICEBERG        int
	TIF            string
	IS_REDUCE_ONLY bool
	CREATE_TIME    float64
	FINISH_TIME    float64
	PRICE          string
}

type ModifyOrderResult struct {
	ID           int
	USER         int
	CONTRACT     string
	CREATETIME   int
	SIZE         float64
	ICEBERG      int
	LEFT         float64
	PRICE        string
	FILLPRICE    string
	MKFR         string
	TKFR         string
	TIF          string
	REFU         int
	ISREDUCEONLY bool
	ISCLOSE      bool
	ISLIQ        bool
	TEXT         string
	STATUS       string
	FINISHTIME   int
	FINISHAS     string
}

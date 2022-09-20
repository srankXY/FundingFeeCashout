package okx

// OkEx相关ApiKey
type OkEx struct {
	APIKEY     string
	SECRETKEY  string
	PASSPHRASE string
}

// SRK
type SetLeverage struct {
	CODE string
	MSG  string
	DATA []LeverageResult
}

type LeverageResult struct {
	LEVER   string `json:"lever"`
	MGNMODE string `json:"mgnMode"`
	INSTID  string `json:"instId"`
	POSSIDE string `json:"posSide"`
}

type ModifyOrderRes struct {
	CODE string
	MSG  string
	DATA []ModifyOrderResult
}

type ModifyOrderResult struct {
	CLORDID string `json:"clOrdId"`
	ORDID   string `json:"ordId"`
	REQID   string `json:"reqId"`
	SCODE   string `json:"sCode"`
	SMSG    string `json:"sMsg"`
}

// 币转张
type CoinSheet struct {
	CODE string
	MSG  string
	DATA []CoinSheetResult
}

type CoinSheetResult struct {
	INSTID string `json:"instId"`
	PX     string
	SZ     string
	TYPE   string
	UNIT   string
}

type PendOrder struct {
	CODE string
	MSG  string
	DATA []PendOrderResult
}

type PendOrderResult struct {
	ACCFILLSZ       string `json:"accFillSz"`
	AVGPX           string `json:"avgPx"`
	CTIME           string `json:"cTime"`
	CATEGORY        string `json:"category"`
	CCY             string `json:"ccy"`
	CLORDID         string `json:"clOrdId"`
	FEE             string `json:"fee"`
	FEECCY          string `json:"feeCcy"`
	FILLPX          string `json:"fillPx"`
	FILLSZ          string `json:"fillSz"`
	FILLTIME        string `json:"fillTime"`
	INSTID          string `json:"instId"`
	INSTTYPE        string `json:"instType"`
	LEVER           string `json:"lever"`
	ORDID           string `json:"ordId"`
	ORDTYPE         string `json:"ordType"`
	PNL             string `json:"pnl"`
	POSSIDE         string `json:"posSide"`
	PX              string `json:"px"`
	REBATE          string `json:"rebate"`
	REBATECCY       string `json:"rebateCcy"`
	SIDE            string `json:"side"`
	SLORDPX         string `json:"slOrdPx"`
	SLTRIGGERPX     string `json:"slTriggerPx"`
	SLTRIGGERPXTYPE string `json:"slTriggerPxType"`
	SOURCE          string `json:"source"`
	STATE           string `json:"state"`
	SZ              string `json:"sz"`
	TAG             string `json:"tag"`
	TDMODE          string `json:"tdMode"`
	TGTCCY          string `json:"tgtCcy"`
	TPORDPX         string `json:"tpOrdPx"`
	TPTRIGGERPX     string `json:"tpTriggerPx"`
	TPTRIGGERPXTYPE string `json:"tpTriggerPxType"`
	TRADEID         string `json:"tradeId"`
	UTIME           string `json:"uTime"`
}

// 产品行情
type Tickers struct {
	CODE string        `json:"code"`
	MSG  string        `json:"msg"`
	DATA []TickersData `json:"data"`
}

// 产品行情Data
type TickersData struct {
	INSTTYPE  string `json:"insttype"`
	INSTID    string `json:"instid"`
	LAST      string `json:"last"`
	LASTSZ    string `json:"lastsz"`
	ASKPX     string `json:"askpx"`
	ASKSZ     string `json:"asksz"`
	BIDPX     string `json:"bidpx"`
	BIDSZ     string `json:"bidsz"`
	OPEN24H   string `json:"open24h"`
	HIGHT24H  string `json:"hight24h"`
	LOW24H    string `json:"low24h"`
	VOLCCY24H string `json:"volccy24h"`
	VOL24H    string `json:"vol24h"`
	SODUTC0   string `json:"sodutc0"`
	SODUTC8   string `json:"sodutc8"`
	TS        string `json:"ts"`
}

// 永续合约
type FundingRate struct {
	CODE string            `json:"code"`
	MSG  string            `json:"msg"`
	DATA []FundingRateData `json:"data"`
}

// 永续合约Data
type FundingRateData struct {
	INSTTYPE        string `json:"insttype"`
	INSTID          string `json:"instid"`
	FUNDINGRATE     string `json:"fundingrate"`
	NEXTFUNDINGRATE string `json:"nextfundingrate"`
	FUNDINGTIME     string `json:"fundingtime"`
	NEXTFUNDINGTIME string `json:"nextfundingtime"`
}

// 指数行情
type IndexTickers struct {
	CODE string             `json:"code"`
	MSG  string             `json:"msg"`
	DATA []IndexTickersData `json:"data"`
}

// 指数行情Data
type IndexTickersData struct {
	INSTID  string `json:"instid"`
	IDXPX   string `json:"idxpx"`
	HIGH24H string `json:"high24h"`
	LOW24H  string `json:"low24h"`
	OPEN24H string `json:"open24h"`
	SODUTC0 string `json:"sodUtc0"`
	SODUTC8 string `json:"sodUtc8"`
	TS      string `json:"ts"`
}

// 资金划转RequestBody
type TransferBody struct {
	CCY       string `json:"ccy"`
	AMT       string `json:"amt"`
	FROM      string `json:"from"`
	TO        string `json:"to"`
	SUBACCT   string `json:"subacct"`
	TYPE      string `json:"type"`
	LOANTRANS bool   `json:"loantrans"`
	CLIENTID  string `json:"clientid"`
}

// 资金划转
type Transfer struct {
	CODE string         `json:"code"`
	MSG  string         `json:"msg"`
	DATA []TransferData `json:"data"`
}

// 资金划转Data
type TransferData struct {
	TRANSID  string `json:"transid"`
	CCY      string `json:"ccy"`
	FROM     string `json:"from"`
	AMT      string `json:"amt"`
	TO       string `json:"to"`
	CLIENTID string `json:"clientid"`
}

// 交易产品信息
type Instruments struct {
	CODE string            `json:"code"`
	MSG  string            `json:"msg"`
	DATA []InstrumentsData `json:"data"`
}

// 交易产品信息Data
type InstrumentsData struct {
	INSTTYPE     string `json:"insttype"`
	INSTID       string `json:"instid"`
	ULY          string `json:"uly"`
	CATEGORY     string `json:"category"`
	BASECCY      string `json:"baseccy"`
	QUOTECCY     string `json:"quoteccy"`
	SETTLECCY    string `json:"settleccy"`
	CTVAL        string `json:"ctval"`
	CTMULT       string `json:"ctmult"`
	CTVALCCY     string `json:"ctvalccy"`
	OPTTYPE      string `json:"opttype"`
	STK          string `json:"stk"`
	LISTTIME     string `json:"listtime"`
	EXPTIME      string `json:"exptime"`
	LEVER        string `json:"lever"`
	TICKSZ       string `json:"ticksz"`
	LOTSZ        string `json:"lotsz"`
	MINSZ        string `json:"minsz"`
	CTTYPE       string `json:"cttype"`
	ALIAS        string `json:"alias"`
	STATE        string `json:"state"`
	MAXLMTSZ     string `json:"maxlmtsz"`
	MAXMKTSZ     string `json:"maxmktsz"`
	MAXTWAPSZ    string `json:"maxtwapsz"`
	MAXICEBERGSZ string `json:"maxicebergsz"`
	MAXTRIGGERSZ string `json:"maxtriggersz"`
	MAXSTOPSZ    string `json:"maxstopsz"`
}

// 下单RequestBody
type OrderBody struct {
	INSTID     string `json:"instid"`
	TDMODE     string `json:"tdmode"`
	CCY        string `json:"ccy"`
	CLORDID    string `json:"clordid"`
	TAG        string `json:"tag"`
	SIDE       string `json:"side"`
	POSSIDE    string `json:"posside"`
	ORDTYPE    string `json:"ordtype"`
	SZ         string `json:"sz"`
	PX         string `json:"px"`
	REDUCEONLY bool   `json:"reduceonly"`
	TGTCCY     string `json:"tgtccy"`
	BANAMEND   bool   `json:"banamend"`
}

// 下单
type Order struct {
	CODE string         `json:"code"`
	MSG  string         `json:"msg"`
	DATA []OderDataPost `json:"data"`
}

// 下单Data
type OderDataPost struct {
	ORDID   string `json:"ordid"`
	CLORDID string `json:"clordid"`
	TAG     string `json:"tag"`
	SCODE   string `json:"scode"`
	SMSG    string `json:"smsg"`
}

// 订单信息
type OrderGet struct {
	CODE string        `json:"code"`
	MSG  string        `json:"msg"`
	DATA []OderDataGet `json:"data"`
}

// 订单信息Data
type OderDataGet struct {
	INSTTYPE        string `json:"insttype"`
	INSTID          string `json:"instid"`
	TGTCCY          string `json:"tgtccy"`
	CCY             string `json:"ccy"`
	ORDID           string `json:"ordid"`
	CLORDID         string `json:"clordid"`
	TAG             string `json:"tag"`
	PX              string `json:"px"`
	SZ              string `json:"sz"`
	PNL             string `json:"pnl"`
	ORDTYPE         string `json:"ordtype"`
	SIDE            string `json:"side"`
	POSSIDE         string `json:"posside"`
	TDMODE          string `json:"tdmode"`
	ACCFILLSZ       string `json:"accfillsz"`
	FILLPX          string `json:"fillpx"`
	TRADEID         string `json:"tradeid"`
	FILLSZ          string `json:"fillsz"`
	FILLTIME        string `json:"filltime"`
	AVGPX           string `json:"avgpx"`
	STATE           string `json:"state"`
	LEVER           string `json:"lever"`
	TPTRIGGERPX     string `json:"tptriggerpx"`
	TPTRIGGERPXTYPE string `json:"tptriggerpxtype"`
	TPORDPX         string `json:"tpordpx"`
	SLTRIGGERPX     string `json:"sltriggerpx"`
	SLTRIGGERPXTYPE string `json:"sltriggerpxtype"`
	SLORDPX         string `json:"slordpx"`
	FEECCY          string `json:"feeccy"`
	FEE             string `json:"fee"`
	REBATECCY       string `json:"rebateccy"`
	SOURCE          string `json:"source"`
	REBATE          string `json:"rebate"`
	CATEGORY        string `json:"category"`
	UTIME           string `json:"utime"`
	CTIME           string `json:"ctime"`
}

// 标记价格
type MarkPrice struct {
	CODE string          `json:"code"`
	MSG  string          `json:"msg"`
	DATA []MarkPriceData `json:"data"`
}

// 标记价格Data
type MarkPriceData struct {
	INSTTYPE string `json:"insttype"`
	INSTID   string `json:"instid"`
	MARKPX   string `json:"markpx"`
	TS       string `json:"ts"`
}

// 持仓信息
type Positions struct {
	CODE string          `json:"code"`
	MSG  string          `json:"msg"`
	DATA []PositionsData `json:"data"`
}

// 持仓信息Data
type PositionsData struct {
	INSTTYPE    string `json:"insttype"`
	MGNMODE     string `json:"mgnmode"`
	POSID       string `json:"posid"`
	POSSIDE     string `json:"posside"`
	POS         string `json:"pos"`
	BASEBAL     string `json:"basebal"`
	QUOTEBAL    string `json:"quotebal"`
	POSCCY      string `json:"posccy"`
	AVAILPOS    string `json:"availpos"`
	AVGPX       string `json:"avgpx"`
	UPL         string `json:"upl"`
	UPLRATIO    string `json:"uplratio"`
	INSTID      string `json:"instid"`
	LEVER       string `json:"lever"`
	LIQPX       string `json:"liqpx"`
	MARKPX      string `json:"markpx"`
	IMR         string `json:"imr"`
	MARGIN      string `json:"margin"`
	MGNRATIO    string `json:"mgnratio"`
	MMR         string `json:"mmr"`
	LIAB        string `json:"liab"`
	LIABCCY     string `json:"liabccy"`
	INTEREST    string `json:"interest"`
	TRADEID     string `json:"tradeid"`
	OPTVAL      string `json:"optval"`
	NOTIONALUSD string `json:"notionalusd"`
	ADL         string `json:"adl"`
	CCY         string `json:"ccy"`
	LAST        string `json:"last"`
	USDPX       string `json:"usdpx"`
	DELTABS     string `json:"deltabs"`
	DELTAPA     string `json:"deltapa"`
	GAMMABS     string `json:"gammabs"`
	GAMMAPA     string `json:"gammapa"`
	THETABS     string `json:"thetabs"`
	THETAPA     string `json:"thetapa"`
	VEGABS      string `json:"vegabs"`
	VEGAPA      string `json:"vegapa"`
	CTIME       string `json:"ctime"`
	UTIME       string `json:"utime"`
}

// 账户余额
type Balance struct {
	CODE string        `json:"code"`
	MSG  string        `json:"msg"`
	DATA []BalanceData `json:"data"`
}

// 账户余额Data
type BalanceData struct {
	UTIME       string           `json:"uTime"`
	TOTALEQ     string           `json:"totalEq"`
	ISOEQ       string           `json:"isoEq"`
	ADJEQ       string           `json:"adjEq"`
	ORDFROZ     string           `json:"ordFroz"`
	IMR         string           `json:"imr"`
	MMR         string           `json:"mmr"`
	MGNRATIO    string           `json:"mgnRatio"`
	NOTIONALUSD string           `json:"notionalUsd"`
	DETAILS     []BalanceDetails `json:"details"`
}

// 账户余额Details
type BalanceDetails struct {
	CCY           string `json:"ccy"`
	EQ            string `json:"eq"`
	CASHBAL       string `json:"cashbal"`
	UTIME         string `json:"utime"`
	ISOEQ         string `json:"isoeq"`
	AVAILEQ       string `json:"availEq"`
	DISEQ         string `json:"diseq"`
	AVAILBAL      string `json:"availbal"`
	FROZENBAL     string `json:"frozenbal"`
	ORDFROZEN     string `json:"ordfrozen"`
	LIAB          string `json:"liab"`
	UPL           string `json:"upl"`
	UPLLIAB       string `json:"uplliab"`
	CROSSLIAB     string `json:"crossliab"`
	ISOLIAB       string `json:"isoliab"`
	MGNRATIO      string `json:"mgnratio"`
	INTEREST      string `json:"interest"`
	TWAP          string `json:"twap"`
	MAXLOAN       string `json:"maxloan"`
	EQUSD         string `json:"equsd"`
	NOTIONALLEVER string `json:"notionallever"`
	STGYEQ        string `json:"stgyeq"`
	ISOUPL        string `json:"isoupl"`
}

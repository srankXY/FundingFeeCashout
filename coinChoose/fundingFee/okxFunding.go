package fundingFee

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/exchangeBase/okx"
	"fmt"
	"strings"
	"time"
)

func PullOkxFundingToDB() {

	for _, tickerData := range okx.GetTickers("SWAP").DATA {

		// 跳过币本位合约
		if strings.Split(tickerData.INSTID, "-")[1] == "USD" {
			continue
		}

		var rate string

		//tickerData.INSTID 即为合约名称

		// 判断数据是否成功获取，未成功获取则循环重新获取
		var result okx.FundingRate
		for result = okx.GetFundingRate(tickerData.INSTID); result.CODE != "0"; result = okx.GetFundingRate(tickerData.INSTID) {
			time.Sleep(100 * time.Millisecond)
		}
		rate = result.DATA[0].FUNDINGRATE

		//for {
		//	result := okx.GetFundingRate(tickerData.INSTID)
		//
		//	if result.CODE == "0" {
		//		rate = result.DATA[0].FUNDINGRATE
		//		break
		//	}
		//	time.Sleep(100 * time.Millisecond)
		//}

		coin := strings.Split(fmt.Sprintf("%s", tickerData.INSTID), okx.Prefix)[0]

		// 存储币种名称
		db.InsertDB(db.FundFeeTableName, "name", coin)

		// 存入资金费率
		db.UpdateDB(db.FundFeeTableName, okx.ExchangeName, rate, "name", coin)
	}

}

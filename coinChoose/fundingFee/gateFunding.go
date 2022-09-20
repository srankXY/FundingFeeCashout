package fundingFee

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/exchangeBase/gate"
	"fmt"
	"strings"
)

func PullGateFundingToDB() {

	allFundRate := gate.GetAllFuturesInfo(gate.Settle)
	for _, data := range allFundRate {
		coin, rate := strings.Split(fmt.Sprintf("%s", data.NAME), gate.Prefix)[0], fmt.Sprintf("%s", data.FUNDING_RATE)
		//fmt.Println(data.NAME, data.FUNDING_RATE)

		// 存储币种名称
		db.InsertDB(db.FundFeeTableName, "name", coin)

		// 存入资金费率
		db.UpdateDB(db.FundFeeTableName, gate.ExchangeName, rate, "name", coin)
	}
}

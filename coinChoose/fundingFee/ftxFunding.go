package fundingFee

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/exchangeBase/ftx"
	"fmt"
	"strings"
)

func PullFtxFundingToDB() {

	allFundRate := ftx.GetFundRate("all").RESULT
	for _, data := range allFundRate {
		coin, rate := strings.Split(fmt.Sprintf("%s", data.FUTURE), ftx.Prefix)[0], fmt.Sprintf("%f", data.RATE)
		//fmt.Println(coin, rate)

		// 存储币种名称
		db.InsertDB(db.FundFeeTableName, "name", coin)

		// 存入资金费率
		db.UpdateDB(db.FundFeeTableName, ftx.ExchangeName, rate, "name", coin)
	}
}

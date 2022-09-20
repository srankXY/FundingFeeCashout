package fundingFee

import (
	"FundingFeeCashout/db"
)

func SaveFundingGap() {

	allFundGap := db.CalcCoinFundFeeGap()
	for _, data := range allFundGap {
		db.UpdateDB(db.FundFeeTableName, "gap", data.Gap, "name", data.Coin)
	}

}

//DB.CheckErr(err)
//DB.UpdateDB(DB.FundFeeTableName, "gap", value, "name", name)

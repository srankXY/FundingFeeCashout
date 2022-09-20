package db

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
CalcCoinFundFeeGap

以列表形式返回计算出来的各个交易所币种的资金费差
*/
func CalcCoinFundFeeGap() []FundingGap {
	var (
		name   string
		value  string
		result []FundingGap
	)
	db := NewDB(DbPath)
	// 关闭db
	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)

	//	拼接sql
	splicingSql := SplicingSQL("SELECT name, '%s' AS exName, MAX(%s) as rate FROM fund_fee_rate group by name UNION ", 2)

	cmd := fmt.Sprintf("select name, max(rate)-min(rate) as gap from (" +
		strings.TrimRight(splicingSql, "UNION ") +
		") GROUP BY name;")

	// 执行查询命令
	rows, err := db.Query(cmd)
	CheckErr(err)

	// 关闭查询结果
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 循环获取查询结果
	for rows.Next() {
		err = rows.Scan(&name, &value)
		result = append(result, FundingGap{Coin: name, Gap: value})
	}

	return result
}

/*
GetMaxCoinAndExchange

返回币种，最大交易所，最小交易锁，和资金费差等信息
*/

func GetMaxCoinAndExchange() CoinExchange {

	// 存储每一行sql查询出来的结果
	var (
		coin string
		// exName, exRate 临时存储sql读出来的结果，用于最终存入result
		exName string
		exRate float64
		gap    string
		// 存储sql 查询出来的结果
		exNameRateSlice []exNameRate
		lastExNameRate  []exNameRate
	)

	db := NewDB(DbPath)
	// 关闭db
	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)
	//	拼接sql

	splicingSql := SplicingSQL("SELECT name, '%s' as exName, %s as rate, max(gap) FROM fund_fee_rate UNION ", 2)
	cmd := strings.TrimRight(splicingSql, "UNION ")

	// 执行查询命令
	rows, err := db.Query(cmd)
	CheckErr(err)

	// 关闭查询结果
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 循环获取查询结果
	for rows.Next() {
		exRate = 0
		err = rows.Scan(&coin, &exName, &exRate, &gap)
		exNameRateSlice = append(exNameRateSlice, exNameRate{name: exName, rate: exRate})
	}

	fmt.Println(exNameRateSlice)

	// 通过循环获取最高资金费及最低资金费交易所
	for i := 0; i < len(exNameRateSlice); i++ {
		for j := 0; j < len(exNameRateSlice)-1-i; j++ {
			if exNameRateSlice[j].rate > exNameRateSlice[j+1].rate {
				tmpData := exNameRateSlice[j]
				exNameRateSlice[j] = exNameRateSlice[j+1]
				exNameRateSlice[j+1] = tmpData
			}
		}
		// 判断剔除最大rate=0的情况
		if exNameRateSlice[len(exNameRateSlice)-1-i].rate != 0 {
			lastExNameRate = append(lastExNameRate, exNameRateSlice[len(exNameRateSlice)-1-i])
		}
	}

	return CoinExchange{
		Coin:        coin,
		MaxExchange: lastExNameRate[0],
		MinExchange: lastExNameRate[len(lastExNameRate)-1],
		Gap:         gap,
	}
}

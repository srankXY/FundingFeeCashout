package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const DbPath = "./ex.db"

// FundFeeTableName 存储资金费率的数据库表名称
const FundFeeTableName = "fund_fee_rate"

// ConfTableName 数据库配置表名称
const ConfTableName = "conf"

//type exDB struct {
//	client *sql.DB
//}

func NewDB(path string) *sql.DB {

	db, _ := sql.Open("sqlite3", path)

	return db
}

func CheckErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func InitDB() {

	// 配置map
	dbconf := map[string]string{

		"PROXY":                          "请输入代理地址，如不使用请输入null， 如: http://127.0.0.1:7890",
		"DEBUG":                          "是否开启debug模式查看详细日志，verbose:打印所有日志, warning:打印除json响应之外的所有日志，info:打印一般信息，留空日志最少",
		"LEVERAGE":                       "请输入合约杠杆倍数， 如: 10",
		"SPILT_COUNT":                    "请输入拆单次数，如：10",
		"PRICE_RATIO":                    "请输入下单价格对于当前实际标记价格的调整比例，只能大于 0.999， 如: 0.9995",
		"OPEN_MODIFY_ORDER_PRICE_OFFSET": "修改挂单时是否也需调整价格，true:是 / false: 否",
		"OKX_KEY":                        "请输入okx key",
		"OKX_SECERT":                     "请输入okx secert",
		"OKX_PASS":                       "请输入okx passhras",
		"FTX_KEY":                        "请输入ftx key",
		"FTX_SECERT":                     "请输入ftx secert",
		"GATE_KEY":                       "请输入gate key",
		"GATE_SECERT":                    "请输入gate secert",
		"PEND_TIMEOUT":                   "未成交挂单超时时间, 以秒为单位，不能低于3s且考虑配置为3的倍数",
		"BALANCE_USED_RATIO":             "【重要】请输入使用空闲资金比例，如：0.9(表示90%的资金用于操作，其余留存保底)，如不留存资金会到底滑点增大，亏损概率增加，请知悉！",
	}
	db := NewDB(DbPath)

	//创建配置表
	confCmd, confErr := db.Prepare("create table ?(name char primary key , value char)")
	// 判断表，存在则不执行sql
	if !CheckErr(confErr) {
		_, err := confCmd.Exec(ConfTableName)
		CheckErr(err)
	}

	// 创建资金费存储表
	fundfCmd, fundErr := db.Prepare("create table ?(name char primary key, " + SplicingSQL("%s float, ", 1) + "gap float)")
	// 判断表，存在则不执行sql
	if !CheckErr(fundErr) {
		_, err := fundfCmd.Exec(FundFeeTableName)
		CheckErr(err)
	}

	//sql
	insert, _ := db.Prepare("INSERT INTO conf(name, value) values(?,?)")
	update, _ := db.Prepare("update conf set value=? where name=?")

	// 循环获取数据，并写入数据库
	for k, v := range dbconf {
		var data string

		fmt.Printf(v + "\n>")
		_, _ = fmt.Scanln(&data)
		//dbconf[k] = data

		//// 设置默认总资产可操作比例
		//if k == "BALANCE_USED_RATIO" && data == "" {
		//	_, _ = insert.Exec(k, 0.9)
		//}

		// 写入数据
		_, err := insert.Exec(k, data)

		// 如果数据存在，则更新数据
		if CheckErr(err) {

			// 如果没有键入值，则不执行update
			if data == "" {
				continue
			}
			_, err := update.Exec(data, k)
			if CheckErr(err) {
				fmt.Println(err)
			} else {
				fmt.Printf("【%s】 更新成功，值：%s \n\n", k, data)
			}
		} else {
			fmt.Printf("【%s】 配置成功，值：%s \n\n", k, data)
		}
	}

	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)
}

/*
QueryDB

查询一张表某一列的具体数据
col: 需要查询的列，tab: 需要查询的表, where*: 查询条件
*/
func QueryDB(col, tab, whereCol, whereVal string) string {

	var value string
	db := NewDB(DbPath)

	// 拼接查询语句
	cmd := fmt.Sprintf("select %s from %s where %s='%s'", col, tab, whereCol, whereVal)

	// 执行查询命令
	rows, err := db.Query(cmd)
	CheckErr(err)

	// 关闭查询结果
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 获取查询结果
	rows.Next()
	err = rows.Scan(&value)
	CheckErr(err)

	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)
	return value
}

/*
UpdateDB

更新某一张表的某一列数据
*/
func UpdateDB(tab, col, val, whereCol, whereVal string) {

	db := NewDB(DbPath)
	// 关闭db
	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)

	// 拼接sql语句
	cmd := fmt.Sprintf("UPDATE %s SET %s=%s WHERE %s='%s'", tab, col, val, whereCol, whereVal)
	//fmt.Println(cmd)
	update, _ := db.Prepare(cmd)

	// 更新数据
	_, err := update.Exec()
	if CheckErr(err) {
	} else {
		fmt.Printf("【%s】 更新成功，列：%s, 币种：%s, 值：%s \n", tab, col, whereVal, val)
	}

}

/*
InsertDB

插入主键列
*/
func InsertDB(tab, col, val string) {

	db := NewDB(DbPath)
	// 关闭db
	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)

	// 拼接sql语句
	//cmd := fmt.Sprintf("REPLACE INTO %s('%s') VALUES('%s')", tab, col, val)
	cmd := fmt.Sprintf("REPLACE INTO %s('%s') SELECT ('%s') WHERE NOT EXISTS (SELECT %s FROM %s WHERE %s='%s')", tab, col, val, col, tab, col, val)
	//fmt.Println(cmd)

	insert, _ := db.Prepare(cmd)

	// 写入数据
	_, err := insert.Exec()

	// 检查错误
	if CheckErr(err) {
	} else {
		fmt.Printf("【%s】 操作成功，列：%s, 值：%s \n", tab, col, val)
	}

}

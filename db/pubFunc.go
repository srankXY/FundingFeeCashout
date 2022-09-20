package db

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
SplicingSQL

根据支持的交易所拼接部分sql，避免重复sql语句

initStr: 预格式化类型的字符串，如："%s float, "
*/
func SplicingSQL(formatStr string, formatCount int) string {

	cmd := ""
	for _, i := range QueryCanUseExchange() {

		if formatCount == 1 {
			cmd = cmd + fmt.Sprintf(formatStr, i)
		} else if formatCount == 2 {
			cmd = cmd + fmt.Sprintf(formatStr, i, i)
		}

	}
	return cmd
}

/*
QueryCanUseExchange

用户可能不是所有支持的交易所都使用，所以需要根据他配置的key，去判断使用了那些交易所
*/
func QueryCanUseExchange() []string {
	var (
		name   string
		value  string
		result []string
	)

	db := NewDB(DbPath)
	// 关闭db
	defer func(db *sql.DB) {
		err := db.Close()
		CheckErr(err)
	}(db)

	cmd := "select * FROM " + ConfTableName + " WHERE name like \"%key%\""
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
		value = ""
		err = rows.Scan(&name, &value)
		if value != "" {
			result = append(result, strings.Split(name, "_")[0])
		}
	}

	return result
}

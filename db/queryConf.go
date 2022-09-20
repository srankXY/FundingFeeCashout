package db

import (
	"strconv"
)

const balanceUseRatio = 0.9
const DefaultPendLoopWait = 3

var Conf *conf

func QueryConf() *conf {
	// 获取拆单次数
	splitCount := QueryDB("value", "conf", "name", "SPILT_COUNT")
	intSplitCount, _ := strconv.Atoi(splitCount)

	// 获取debug配置
	strDebug := QueryDB("value", "conf", "name", "DEBUG")

	// 获取合约倍数
	leverage := QueryDB("value", "conf", "name", "LEVERAGE")
	intLeverage, _ := strconv.Atoi(leverage)

	// 获取修改订单的价格调整策略
	modifyOrderPriceOffset := QueryDB("value", "conf", "name", "OPEN_MODIFY_ORDER_PRICE_OFFSET")
	if modifyOrderPriceOffset == "" {
		modifyOrderPriceOffset = "false"
	}

	// 查询资金使用比例
	balanceUsedRatio := QueryDB("value", "conf", "name", "BALANCE_USED_RATIO")
	floatBalanceUsedRatio, _ := strconv.ParseFloat(balanceUsedRatio, 64)
	if balanceUsedRatio == "" {
		floatBalanceUsedRatio = balanceUseRatio
	}

	// 获取开仓价格上调比例
	priceRatio := QueryDB("value", "conf", "name", "PRICE_RATIO")
	floatPriceRatio, _ := strconv.ParseFloat(priceRatio, 64)

	// 获取挂单超时时间
	pendTimeout := QueryDB("value", "conf", "name", "PEND_TIMEOUT")
	intPendTimeout, _ := strconv.Atoi(pendTimeout)
	if pendTimeout == "" || intPendTimeout < 5 {
		intPendTimeout = DefaultPendLoopWait
	}

	return &conf{
		LEVERAGE:                       intLeverage,
		SPILT_COUNT:                    intSplitCount,
		PRICE_RATIO:                    floatPriceRatio,
		BALANCE_USED_RATIO:             floatBalanceUsedRatio,
		PEND_TIMEOUT:                   intPendTimeout,
		DEBUG:                          strDebug,
		OPEN_MODIFY_ORDER_PRICE_OFFSET: modifyOrderPriceOffset,
	}
}

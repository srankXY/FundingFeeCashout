package lib

import (
	"FundingFeeCashout/db"
	"fmt"
	"srkTools"
	"time"
)

var PlaceOrderSyncSlice []string
var SyncStats string

/*
SyncProcess

两个不同进程之间的通信同步，目前主要用于不同交易所订单挂单未成交时的同步等待操作， 可能存在一方已经执行完成，但另外一方还有任务需要执行不需要再等待对方应答的情况

所以在一个进程执行完成之后，必须要配合SyncStats 变量使用，且值只能为ok

例子: SyncProcess(&PlaceOrderSyncSlice, "OKX")
*/
//goland:noinspection GoNilness
func SyncProcess(syncSlice *[]string, processName string) {

	if SyncStats != "ok" {
		// 添加信号到全局slice
		*syncSlice = append(*syncSlice, processName)
	}
	srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】当前通信通道为: %v，同步状态为: [%s] \n", processName, *syncSlice, SyncStats))
	// 临时slice，用于存储除自己信号之外的其他信号
	var tmpOtherSlice []string
	// 判断除自己之外是否还存在其他信号
	for true {
		if SyncStats == "ok" {
			*syncSlice = []string{}
			break
		}
		for _, v := range *syncSlice {
			if v != processName {
				tmpOtherSlice = append(tmpOtherSlice, v)
			}
		}
		if len(tmpOtherSlice) == 1 {
			break
		}
		srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】正在等待对方订单成交\n", processName))
		time.Sleep(db.DefaultPendLoopWait * time.Second)
	}

	// 清除（消费）除自己之外的其他信号
	if len(tmpOtherSlice) != 0 {
		srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】对方为：%s\n", processName, tmpOtherSlice[0]))
		*syncSlice = func(elem string) []string {
			index := 0
			for _, v := range *syncSlice {
				if v != elem {
					(*syncSlice)[index] = v
					index++
				}
			}
			return (*syncSlice)[:index]
		}(tmpOtherSlice[0])
	}

	// 判断自己的信号是否被消除（消费）
	for true {
		if SyncStats == "ok" {
			break
		}
		var tmpSelfSlice []string
		for _, v := range *syncSlice {
			if v == processName {
				tmpSelfSlice = append(tmpSelfSlice, v)
			}
		}
		if len(tmpSelfSlice) == 0 {
			break
		}
		srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】正在等待对方处理通信通道\n", processName))
		time.Sleep(db.DefaultPendLoopWait * time.Second)
	}
}

/*
FindElem

查找slice中是否包含某个元素
*/
func FindElem(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

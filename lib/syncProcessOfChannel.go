package lib

import (
	"FundingFeeCashout/db"
	"fmt"
	"srkTools"
	"time"
)

// PlaceOrderSyncChan channel 为引用类型
var PlaceOrderSyncChan = make(chan string, 2)

/*
SyncProcessChan
!未测试!

使用前应先传入信号到通道中，再调用同步方法 / 也可以在函数内第一步传入processName到通道中，如：
PlaceOrderSyncChan <- "OKX"
SyncProcessChan(PlaceOrderSyncChan, "OKX")

下单全部结束应该关闭通道
close(PlaceOrderSyncChan)
*/
func SyncProcessChan(SyncChan chan string, processName string) {
	for {
		signal, ok := <-SyncChan

		// 如果通道关闭，则退出循环
		if !ok {
			srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】通道已关闭，对方已完全成交，不需等待 \n", processName))
			break
		}

		if signal == processName {
			SyncChan <- processName
			srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】等待对方成交 \n", processName))
		} else {
			srkTools.DebugLog(DebugLevel.INFO, fmt.Sprintf("【%s】对方为： \n", signal))
			break
		}

		time.Sleep(db.DefaultPendLoopWait * time.Second)
	}
}

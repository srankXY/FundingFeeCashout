package appConf

const CurrentVersion = "1.3"

// MinPrice 低于该价格将不会按照策略调整价格
const MinPrice = 0.0001

// db.DefaultPendLoopWait 默认未成交订单等待时间（交易同步和挂单等待使用）

//MaxPendLoopCount 最大支持多少次正常价格调整（+-价格，超过这个次数将获取买一/卖一价）
const MaxPendLoopCount = 2

const SupportExchange = "ftx, okx, gate"

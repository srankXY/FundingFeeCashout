# 交易所资金费套利脚本V1



## 简介



> 一款免费的交易所资金费套现工具！



## 资金费套利原理



> 具体请查看：[资金费套利原理](https://github.com/srankXY/FundingFeeCashout/blob/master/%E8%B5%84%E9%87%91%E8%B4%B9%E5%A5%97%E5%88%A9%E5%8E%9F%E7%90%86.md)



## 测试数据分享



- V1.0 盈利概率为 50%
- V1.1 正在测试中



## 功能介绍



- 设置合约倍数
- 计算并统一各交易所最优开仓数量（币数量 & 合约张数）
- 各交易所币转张，张转币
- 下单拆单，降低滑点`滑点可以简单理解为亏损率`
- 根据下单方向调整最优下单价格（开多/平空 压低价格，开空/平多 提高价格）
- 挂单检测，未成交挂单在一定的`超时时间`之后重新调整价格
- 预留账户资金（如：有1000u，但是可只用500u操作）
- 配置代理（适合大陆用户）



## 目前支持的交易所



- OKX
- FTX
- GATE



## 操作前的准备



- 选择好币种和交易所以及需要操作的方向，选择规则请参考原理
- 在需要操作的两个交易所准备相同的资金，比如：各1000u
- 核实资金是否在永续合约账户
- 创建好对应的api key， api secret
- 核实交易所持仓模式是否为单向持仓模式



## 部署

### 下载app:

```shell
widows:
wget https://github.com/srankXY/FundingFeeCashout/releases/download/V1.0/FundingFeeCashout-V1.1.exe

wget https://github.com/srankXY/FundingFeeCashout/releases/download/V1.0/windowsRun.ps1


linux:
wget https://github.com/srankXY/FundingFeeCashout/releases/download/V1.0/FundingFeeCashout-linux-V1.1
```

### 初始化:

```shell
FundingFeeCashout init
```

配置会写入当前目录名为ex.db的sqlite当中

#### DB配置项解释:

> ```
> "PROXY":              "代理地址，不支持认证，使用http & https 协议开头，不使用可留空",
> 
> "DEBUG":              "debug模式，目前支持4个等级
> 					   verbose:    打印所有日志,最详细，包括api请求的json数据及err数据
>                        warning:    打印除json & err 响应之外的所有日志
>                        info:       打印一般信息（会包含各交易所计算时的一些数据）
>                        留空:        日志最少, 只会显示关键信息",
> 
> "LEVERAGE":           "合约杠杆倍数，单向持仓，逐仓模式",
> 
> "SPILT_COUNT":        "总共需要拆成多少次进行下单，只能为正整数",
> 
> "PRICE_RATIO":        "下单价格的调整比例，适用于波段行情，且调整的比例应该尽量小，单边行情时请配置为1，例子：0.9999
> 					  开多/平空时表示压低价格为:     price * (1-(1-0.9999))  =  price * 0.9999
> 					  开空/平多时表示提高价格为:     price * (1+(1-0.9999))  =  price / 0.9999",
> 
> "OPEN_MODIFY_ORDER_PRICE_OFFSET":     "修改未成交挂单时，是否也按照 PRICE_RATIO 调整价格
> 									  true：调整
> 									  false：不调整(改成最新成交价)",
> "PEND_TIMEOUT":       "未成交挂单的超时时间，秒为单位，超时之后将修改挂单价格为最新成交价，且不做价格调整",
> 
> "BALANCE_USED_RATIO": "账户预留资金，默认必须要预留一部分，作用：
> 					  - 预留资金预防计算下单量 到 实际下单这个过程中，价格发生变化，导致进仓失败的情况
> 					  - 可用于本来就想预留一部分资金，不动的情况",
> ```

### 运行

```shell
windows:
双击 windowsRun.ps1 

linux:
FundingFeeCashout-linux-V1.1
```

根据提示输入 `开仓` 还是 `平仓`， 以及 `对应的` 交易所



## 结果截图



- OKX:

![ftx.png](https://raw.githubusercontent.com/srankXY/FundingFeeCashout/master/ftx.png)

- GATE:

![gate.png](https://raw.githubusercontent.com/srankXY/FundingFeeCashout/master/gate.png)



## 未来计划



- 自动计算资金费差最大的币种及对应的交易所，自动选择 -> 自动开单
- 仓位监测
  - 爆仓监测
  - 变盘监测
  - 收益监测
- 自动配平交易所资金



## 学习交流



- TG：[FundingFeeCashout CHAT](https://t.me/+rMPBL3WAMWY4M2E9)



## 免责提醒



> 投资均有风险，亏损作者概不负责！请谨慎操作！



-OVER-






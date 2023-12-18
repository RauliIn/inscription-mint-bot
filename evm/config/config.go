package config

type Config struct {
	PriKeys []string //私钥列表
	//ethrpc请求配置
	EthRpcConf struct {
		Url           string  // Url:请求rpc地址 修改来替换各链rpc地址
		IntervalTime  int     //请求轮询时间
		GasPriceRatio float64 //gas比例,比例越高花费gas越多,默认1.01
	}
	//minted限制,用来过滤非热门铭文
	MintConf struct {
		ToAddr    string //发送地址
		InputData string //发送数据

	}
}

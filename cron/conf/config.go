package conf

type CronConfig struct {
	CronList []CronItem
}

type CronItem struct {
	TimeFormat string
	Function   func()
}

/*
 *  与 linux crontab 设置规则是一样的
 * 	具体规则   https://zh.wikipedia.org/wiki/Cron
 */

var Crontab = CronConfig{CronList: []CronItem{
	// {
	//	TimeFormat: "@every 10s", // 每10秒跑一次
	//	Function:   task.RepairDataSaleOrder,
	// },
}}

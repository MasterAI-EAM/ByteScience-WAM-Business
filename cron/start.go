package cron

import (
	"ByteScience-WAM-Business/cron/conf"
	"ByteScience-WAM-Business/cron/task"
	"github.com/robfig/cron/v3"
)

// Start 定时任务
func Start() {
	c := cron.New()

	// 读取配置添加定时任务
	for _, data := range conf.Crontab.CronList {
		c.AddFunc(data.TimeFormat, data.Function)
	}

	// 启动执行计划
	c.Start()
}

// StartForever 持续消耗的任务
func StartForever() {
	go task.ParsedExperimentDataForever()
}

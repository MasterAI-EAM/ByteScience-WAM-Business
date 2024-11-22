package redis

import (
	"ByteScience-WAM-Business/pkg/logger" // 引入日志包
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// redisHook 结构体，记录 Redis 命令执行的日志
type redisHook struct {
	log *logger.LogrusGormLogger
}

// BeforeProcess 在 Redis 命令处理之前调用
func (h *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// 记录 Redis 命令操作
	if h.log != nil {
		h.log.Logger.Infof("Executing Redis Command: %s | Args: %v", cmd.Name(), cmd.Args())
	}

	// 记录命令开始的时间
	return context.WithValue(ctx, "startTime", time.Now()), nil
}

// AfterProcess 在 Redis 命令处理之后调用
func (h *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	// 获取开始时间
	startTime, ok := ctx.Value("startTime").(time.Time)
	if !ok {
		startTime = time.Now() // 如果没有找到开始时间，就使用当前时间
	}

	duration := time.Since(startTime) // 获取操作执行的时间
	if h.log != nil {
		if cmd.Err() != nil {
			h.log.Logger.Errorf("Redis Command Error: %s | Args: %v | Duration: %vms | Error: %v", cmd.Name(), cmd.Args(), duration.Milliseconds(), cmd.Err())
		} else {
			// 记录执行成功的命令，带上执行时长
			h.log.Logger.Infof("Redis Command Success: %s | Args: %v | Duration: %vms", cmd.Name(), cmd.Args(), duration.Milliseconds())
		}
	}
	return nil
}

// BeforeProcessPipeline 在 Redis 命令批处理之前调用
func (h *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	// 记录批处理命令的日志
	if h.log != nil {
		h.log.Logger.Infof("Executing Redis Pipeline Commands: %v", cmds)
	}

	// 记录每个命令的开始时间
	for _, cmd := range cmds {
		ctx = context.WithValue(ctx, cmd.Name(), time.Now())
	}
	return ctx, nil
}

// AfterProcessPipeline 在 Redis 命令批处理之后调用
func (h *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	// 记录批处理命令的执行结果
	if h.log != nil {
		for _, cmd := range cmds {
			startTime, ok := ctx.Value(cmd.Name()).(time.Time)
			if !ok {
				startTime = time.Now()
			}
			duration := time.Since(startTime) // 获取操作执行的时间
			if cmd.Err() != nil {
				h.log.Logger.Errorf("Redis Pipeline Command Error: %s | Args: %v | Duration: %vms | Error: %v", cmd.Name(), cmd.Args(), duration.Milliseconds(), cmd.Err())
			} else {
				// 记录执行成功的命令，带上执行时长
				h.log.Logger.Infof("Redis Pipeline Command Success: %s | Args: %v | Duration: %vms", cmd.Name(), cmd.Args(), duration.Milliseconds())
			}
		}
	}
	return nil
}

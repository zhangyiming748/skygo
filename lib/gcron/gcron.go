package gcron

import (
	"github.com/robfig/cron/v3"
)

// DefaultParser 默认情况下支持“* * * * * *”格式的cron参数
var DefaultParser cron.Parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

func NewDefaultCron() *cron.Cron {
	return cron.New(cron.WithParser(DefaultParser), cron.WithChain())
}

package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.L().Error("sentinel.InitDefault(),err:", zap.Error(err))
		//	log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "api",
			Threshold:              100,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		zap.L().Error("flow.LoadRules([]*flow.Rule,err:", zap.Error(err))
		//	log.Fatalf("Unexpected error: %+v", err)
		return
	}

	_, err = hotspot.LoadRules([]*hotspot.Rule{
		{
			Resource:        "limit_ip",
			MetricType:      hotspot.QPS,
			ControlBehavior: hotspot.Reject,
			ParamIndex:      15,
			Threshold:       2,
			BurstCount:      0, // 控制令牌
			DurationInSec:   3,
		},
	})
	if err != nil {
		zap.L().Error("hotspot.LoadRules([]*hotspot.Rule,err:", zap.Error(err))
		//	log.Fatalf("Unexpected error: %+v", err)
		return
	}
}

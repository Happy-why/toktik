package middleware

import (
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// LimitIP 对 ip 进行限流
func LimitIP() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ip := ctx.ClientIP()
		//fmt.Println("ip:", ip)
		entry, err := sentinel.Entry(
			"limit_ip",
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
			sentinel.WithArgs(ip),
		)
		if err != nil {
			ctx.AbortWithStatusJSON(400, utils.H{
				"err":  "too many request; the quota used up",
				"code": 10222,
			})
			return
		}
		defer entry.Exit()
		ctx.Next(c)
	}
}

package prop

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/model/lobby"
	"kxmj.common/web"
	"kxmj.shop/internal/model"
)

// Demo 该代码只供参考，开发时删除掉
// @Description APP
// @Tags APP
// @Summary 获取网关地址
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.DemoInfo} "请求成功"
// @Router	/prop/demo [GET]
func (s *Service) Demo(ctx *gin.Context) {
	reply := &lobby.GetGatewayResp{}
	err := s.XServer().LobbyClient().Call(ctx, "GetGateways", &lobby.GetGatewayReq{}, reply)
	if err != nil {
		web.RespFailed(ctx, codes.ServerNetErr)
		return
	}

	resp := &model.DemoResp{}
	for _, d := range reply.List {
		resp.List = append(resp.List, &model.DemoInfo{
			SvrType: d.SvrType,
			SvrId:   d.SvrId,
			Addr:    d.Addr,
			Port:    d.Port,
		})
	}

	web.RespSuccess(ctx, resp)
}

package lobby

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.core.api/internal/model"
)

// GetWallet 获取用户钱包信息
// @Description LOBBY
// @Tags LOBBY
// @Summary 获取用户钱包信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetWalletResp} "请求成功"
// @Router	/lobby/get-wallet [GET]
func (s *Service) GetWallet(ctx *gin.Context) {
	userId := uint32(web.GetUserId(ctx))
	wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.GetWalletInfoFailed)
		return
	}

	resp := &model.GetWalletResp{
		Diamond:  wallet.Diamond,
		Gold:     wallet.Gold,
		GoldBean: wallet.GoldBean,
	}

	web.RespSuccess(ctx, resp)
}

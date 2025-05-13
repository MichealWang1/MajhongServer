package email

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/entities/kxmj_logger"
	"kxmj.common/entities/kxmj_report"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/mq"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/mail/welfare"
	"kxmj.common/utils"
	"kxmj.common/web"
	"kxmj.email/internal/model"
	"time"
)

// GetUserAllMails 获取用户邮件列表
// @Description Email
// @Tags Email
// @Summary 获取用户邮件列表
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.UserMailListResp} "请求成功"
// @Router	/email-list [GET]
func (s *Service) GetUserAllMails(ctx *gin.Context) {
	resp := &model.UserMailListResp{}
	// 获取玩家ID
	userId := web.GetUserId(ctx)
	// 获取邮件配置
	configMail, err := redis_cache.GetCache().GetMailCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		log.Sugar().Errorf(" GetUserAllMails GetDetailCache().GetAll err %v ", err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	// 玩家所有的邮件
	userSystemMail, err := redis_cache.GetCache().GetMailCache().GetSystemMailCache().GetAll(ctx, uint32(userId))
	if err != nil {
	}
	userWelfareMail, err := redis_cache.GetCache().GetMailCache().GetWelfareMailCache().GetAll(ctx, uint32(userId))
	if err != nil {
	}
	// 当前时间
	now := uint32(time.Now().Unix())
	// 需要增加到玩家redis的邮件
	addUserMail := make(map[uint32]*kxmj_core.ConfigEmail, 0)
	// 循环邮件配置  key 是 mailId value 是 邮件内容
	for _, mailData := range configMail {
		// 先判断过期时间 过期了则不显示 ExpireAt过去时间 如果 等于 0 则表示永远不过期
		if mailData.ExpireAt != 0 && now >= mailData.ExpireAt {
			continue
		}
		// 如果当前时间没有到发送时间 则 continue
		if now < mailData.SendAt {
			continue
		}
		var status = uint8(0)       // 邮件状态
		var isAddEmail = false      // isAddEmail==true 表示玩家的redis中 没有当前这封邮件 需要添加到redis
		sendTime := mailData.SendAt // 发送时间
		if mailData.EmailType == uint8(model.System) {
			status, isAddEmail = checkUserHaveSystemMail(userSystemMail, mailData.EmailId)
		} else {
			status, sendTime, isAddEmail = checkUserHaveWelfareMail(userWelfareMail, mailData.EmailId, now)
		}
		if mailData.IsSingleSend == int8(model.SingSend) && isAddEmail == true { //单独发送邮件 并且 玩家没有这封邮件
			continue
		}
		if status == uint8(model.Delete) { // 邮件已删除
			continue
		}
		var itemList []*model.MailItem // 取出邮件中能领取的物品
		if mailData.IsReward == uint8(model.HaveItem) {
			err = json.Unmarshal([]byte(mailData.ItemList), &itemList)
			if err != nil || len(itemList) <= 0 {
				continue
			}
		} else {
			itemList = nil
		}
		
		resp.List = append(resp.List, &model.MailData{
			EmailId:   mailData.EmailId,
			EmailType: mailData.EmailType,
			Title:     mailData.Title,
			Remark:    mailData.Remark,
			IsReward:  mailData.IsReward,
			ItemList:  itemList,
			Status:    status,
			CreatedAt: sendTime,
		})
		// 是否要添加到玩家的redis列表中
		if isAddEmail == true {
			addUserMail[mailData.EmailId] = mailData
		}
	}
	if len(addUserMail) > 0 {
		for _, mailValue := range addUserMail {
			// 玩家没有此邮件 把该邮件 初始化到 该玩家的 redis中
			err = AddUserEmail(ctx, uint32(userId), mailValue, now)
			if err != nil {
				log.Sugar().Errorf(" GetUserAllMails AddUserEmail userId:%d mailId:%d err %v ", userId, mailValue.EmailId, err)
			}
		}
	}
	web.RespSuccess(ctx, resp)
}

// SetUserMailRead 设置用户邮件为已读
// @Description Email
// @Tags Email
// @Summary 设置用户邮件为已读
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.SetMailDataReq true "JSON"
// @Success 200 {object} web.Response{data=model.Empty} "请求成功"
// @Router	/set-email-read [POST]
func (s *Service) SetUserMailRead(ctx *gin.Context) {
	// 获取玩家ID
	userId := web.GetUserId(ctx)
	// 取出设置邮件状态 结构体
	setMailReq := &model.SetMailDataReq{}
	err := ctx.ShouldBind(setMailReq)
	if err != nil {
		log.Sugar().Errorf("SetUserMailRead setMailReq ShouldBind error user:%d err:%v", userId, err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	// 要领取邮件物品的邮件ID
	mailSystemList := make(map[uint32]uint8, 0)
	mailWelfareList := make(map[uint32]*welfare.EmailUser, 0)
	if setMailReq.EmailId <= 0 {
		getUserAllSystemMail(ctx, uint32(userId), uint8(model.Read), mailSystemList)
		getUserAllWelfareMail(ctx, uint32(userId), uint8(model.Read), mailWelfareList)
	} else {
		// 获取邮件配置 configMail 则是存储在redis中的邮件配置
		configMail, err := redis_cache.GetCache().GetMailCache().GetDetailCache().GetAll(ctx)
		if err != nil {
			// redis 没有邮件配置数据
			log.Sugar().Errorf("SetUserMailRead GetMailCache().GetDetailCache().GetAll error user:%d err:%v ", userId, err)
			web.RespFailed(ctx, codes.DbError)
			return
		}
		mailConfig, has := configMail[setMailReq.EmailId]
		if has == false {
			web.RespFailed(ctx, codes.NotMailID)
			return
		}
		if mailConfig.EmailType == uint8(model.System) {
			respCode := getUserSystemMailByMailId(ctx, uint32(userId), uint8(model.Read), setMailReq.EmailId, mailSystemList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		} else {
			respCode := getUserWelfareMailByMailId(ctx, uint32(userId), uint8(model.Read), setMailReq.EmailId, mailWelfareList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		}
	}
	if len(mailSystemList) <= 0 && len(mailWelfareList) <= 0 {
		web.RespSuccess(ctx, nil)
		return
	}
	changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Read))   // 系统邮件改变状态
	changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Read)) // 福利邮件改变状态
	web.RespSuccess(ctx, nil)
}

// TakeMailItem 领取邮件里面的物品
// @Description Email
// @Tags Email
// @Summary 领取邮件里面的物品
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.SetMailDataReq true "JSON"
// @Success 200 {object} web.Response{data=model.TakeItemResp} "请求成功"
// @Router	/take-email-item [POST]
func (s *Service) TakeMailItem(ctx *gin.Context) {
	resp := &model.TakeItemResp{}
	// 获取玩家ID
	userId := web.GetUserId(ctx)
	// 取出设置邮件状态 结构体
	setMailReq := &model.SetMailDataReq{}
	err := ctx.ShouldBind(setMailReq)
	if err != nil {
		log.Sugar().Errorf("TakeMailItem ShouldBind error user:%d err:%v ", userId, err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	// 获取邮件配置 configMail 则是存储在redis中的邮件配置
	configMail, err := redis_cache.GetCache().GetMailCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		// redis 没有邮件配置数据
		log.Sugar().Errorf("TakeMailItem GetDetailCache().GetAll error user:%d err:%v ", userId, err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	// 要领取邮件物品的邮件ID
	mailSystemList := make(map[uint32]uint8, 0)
	mailWelfareList := make(map[uint32]*welfare.EmailUser, 0)
	// 读取邮件配置 获取除物品
	if setMailReq.EmailId <= 0 {
		getUserAllSystemMail(ctx, uint32(userId), uint8(model.Take), mailSystemList)
		getUserAllWelfareMail(ctx, uint32(userId), uint8(model.Take), mailWelfareList)
	} else {
		mailConfig, has := configMail[setMailReq.EmailId]
		if has == false {
			web.RespFailed(ctx, codes.NotMailID)
			return
		}
		if mailConfig.EmailType == uint8(model.System) {
			respCode := getUserSystemMailByMailId(ctx, uint32(userId), uint8(model.Take), setMailReq.EmailId, mailSystemList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		} else {
			respCode := getUserWelfareMailByMailId(ctx, uint32(userId), uint8(model.Take), setMailReq.EmailId, mailWelfareList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		}
	}
	if len(mailSystemList) <= 0 && len(mailWelfareList) <= 0 {
		web.RespFailed(ctx, codes.NoTakeItemByMail)
		return
	}
	// 获取邮件配置查看是否有物品领取 key 是 邮件ID  value是 ItemId集合
	needTakeItemMails := analysisItemList(mailWelfareList, mailSystemList, configMail)
	// 判断领取物品
	if needTakeItemMails == nil {
		// 没有可以领取的物品 把邮件状态设置成 已读或者已领取
		// 邮件改变状态
		changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Take))
		changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Take))
		web.RespFailed(ctx, codes.NoTakeItemByMail)
		return
	}
	// 获取游戏物品配置
	configItem, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		// redis 没有 游戏物品配置
		log.Sugar().Errorf("TakeMailItem GetItemCache().GetDetailCache().GetAll error user:%d err:%v ", userId, err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	//  把需要领取的物品提取出来 物品集合 ValueItem 集合
	valueItemList := analysisValueItem(needTakeItemMails, configItem)
	if valueItemList == nil || len(valueItemList) <= 0 {
		// 邮件改变状态
		changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Take))
		changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Take))
		web.RespFailed(ctx, codes.NoTakeItemByMail)
		return
	}
	itemMap := getMailItemGroupMap(needTakeItemMails)
	for itemId, itemCount := range itemMap {
		resp.ItemList = append(resp.ItemList, &model.MailItem{
			Id:    itemId,
			Count: itemCount,
		})
	}
	// 玩家领取邮件物品 把物品列表valueItemList 加载到玩家身上
	s.UserTakePropByEmail(ctx, uint32(userId), valueItemList)
	// 设置邮件的状态
	changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Take))
	changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Take))
	web.RespSuccess(ctx, resp)
}

// DelUserMail 用户删除邮件
// @Description Email
// @Tags Email
// @Summary 用户删除邮件
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.SetMailDataReq true "JSON"
// @Success 200 {object} web.Response{data=model.TakeItemResp} "请求成功"
// @Router	/del-email [POST]
func (s *Service) DelUserMail(ctx *gin.Context) {
	resp := &model.TakeItemResp{}
	// 获取玩家ID
	userId := web.GetUserId(ctx)
	// 取出设置邮件状态 结构体
	setMailReq := &model.SetMailDataReq{}
	err := ctx.ShouldBind(setMailReq)
	// 这里判断 客户端发来的消息判断 是否正确 4 为删除
	if err != nil {
		log.Sugar().Errorf("DelUserMail ShouldBind ParamError error user:%d err:%v", userId, err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	// 所有邮件配置
	mailAllConfig, err := redis_cache.GetCache().GetMailCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		log.Sugar().Errorf("DelUserMail GetMailCache().GetDetailCache().GetAll error user:%d err:%v", userId, err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	// 要删除的邮件ID
	mailSystemList := make(map[uint32]uint8, 0)
	mailWelfareList := make(map[uint32]*welfare.EmailUser, 0)
	if setMailReq.EmailId <= 0 {
		getUserAllSystemMail(ctx, uint32(userId), uint8(model.Delete), mailSystemList)
		getUserAllWelfareMail(ctx, uint32(userId), uint8(model.Delete), mailWelfareList)
	} else {
		mailConfig, has := mailAllConfig[setMailReq.EmailId]
		if has == false {
			web.RespFailed(ctx, codes.NotMailID)
			return
		}
		if mailConfig.EmailType == uint8(model.System) {
			respCode := getUserSystemMailByMailId(ctx, uint32(userId), uint8(model.Delete), setMailReq.EmailId, mailSystemList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		} else {
			respCode := getUserWelfareMailByMailId(ctx, uint32(userId), uint8(model.Delete), setMailReq.EmailId, mailWelfareList)
			if respCode != 0 {
				web.RespFailed(ctx, respCode)
				return
			}
		}
	}
	if len(mailSystemList) <= 0 && len(mailWelfareList) <= 0 {
		web.RespSuccess(ctx, resp)
		return
	}
	
	// 获取邮件配置查看是否有物品领取 key 是 邮件ID  value是 ItemId集合
	needTakeItemMails := analysisItemList(mailWelfareList, mailSystemList, mailAllConfig)
	// 判断领取物品
	if needTakeItemMails == nil {
		// 没有邮件可以领取物品  把需要删除的邮件改变状态
		changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Delete))
		changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Delete))
		web.RespSuccess(ctx, resp)
		return
	}
	// 获取所有物品
	configItem, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		// redis 没有 游戏物品配置
		log.Sugar().Errorf("DelUserMail GetItemCache().GetDetailCache().GetAll error user:%d err:%v ", userId, err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	//  把需要领取的物品提取出来 物品集合 ValueItem 集合
	valueItemList := analysisValueItem(needTakeItemMails, configItem)
	if valueItemList == nil {
		changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Delete))
		changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Delete))
		web.RespSuccess(ctx, resp)
		return
	}
	itemMap := getMailItemGroupMap(needTakeItemMails)
	for itemId, itemCount := range itemMap {
		resp.ItemList = append(resp.ItemList, &model.MailItem{
			Id:    itemId,
			Count: itemCount,
		})
	}
	// 玩家领取邮件物品 把物品列表valueItemList 加载到玩家身上
	s.UserTakePropByEmail(ctx, uint32(userId), valueItemList)
	changeSystemMailStatus(ctx, mailSystemList, uint32(userId), uint8(model.Delete))
	changeWelfareMailStatus(ctx, mailWelfareList, uint32(userId), uint8(model.Delete))
	web.RespSuccess(ctx, resp)
}

// UserTakePropByEmail 从邮件中取出的道具加载到玩家身上
func (s *Service) UserTakePropByEmail(ctx *gin.Context, userId uint32, valueItemList []*item.ValueItem) {
	// 创建一个 orderId 用于查询
	orderId := utils.Snowflake.Generate().Int64()
	// 领取币类
	goldItems := item.GetGoldItems(valueItemList)
	for _, g := range goldItems {
		diamond := "0"
		gold := "0"
		goldBean := "0"
		if g.ItemType == item.Diamond {
			diamond = g.Count
		} else if g.ItemType == item.Gold {
			gold = g.Count
		} else if g.ItemType == item.GoldBean {
			goldBean = g.Count
		}
		args := &center.AddUserWalletReq{
			UserId:       userId,
			OrderId:      orderId,
			Diamond:      diamond,
			Gold:         gold,
			GoldBean:     goldBean,
			BusinessType: 3,
		}
		reply := &center.AddUserWalletResp{}
		err := s.XServer().CenterClient().Call(ctx, "AddUserWallet", args, reply)
		if err != nil {
			log.Sugar().Errorf("DelUserMail AddUserWallet email user:%d order:%d err:%v", userId, orderId, err)
		}
		
	}
	// 领取物品
	propItems := item.GetPropItems(valueItemList)
	err := item.UpdateUserItems(ctx, propItems, userId)
	if err != nil {
		log.Sugar().Errorf("DelUserMail AddUserItem email user:%d order:%d err:%v", userId, orderId, err)
	}
	// 写物品领取日志
	for _, p := range propItems {
		err = mq.AddLogger(&kxmj_logger.ItemTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      orderId,
			UserId:       userId,
			ItemId:       p.ItemId,
			Count:        p.Count,
			BusinessType: 3,
			Type:         1,
			CreatedAt:    uint32(time.Now().Unix()),
		})
	}
}

// checkUserHaveSystemMail 判断该玩家是否有该系统邮件 返回该邮件状态
func checkUserHaveSystemMail(userMail map[uint32]uint8, mailId uint32) (uint8, bool) {
	if userMail == nil {
		return uint8(model.UnRead), true
	}
	status, has := userMail[mailId]
	if has == false {
		return uint8(model.UnRead), true
	}
	return status, false
}

// checkUserHaveWelfareMail 判断该玩家是否有该福利邮件 返回该邮件状态
func checkUserHaveWelfareMail(userMail map[uint32]*welfare.EmailUser, mailId uint32, now uint32) (uint8, uint32, bool) {
	if userMail == nil {
		return uint8(model.UnRead), now, true
	}
	mailInfo, has := userMail[mailId]
	if has == false {
		return uint8(model.UnRead), now, true
	}
	return mailInfo.Status, mailInfo.SendTime, false
}

// mailId 邮件ID userId 玩家ID status邮件状态 time 更新时间
func setMailStatus(ctx *gin.Context, userId uint32, mailId uint32, mailType uint8, status uint8, orderId int64, sendTime, updateTime, drawTime uint32) error {
	if mailId <= 0 || userId <= 0 || status < uint8(model.UnRead) || status > uint8(model.Delete) {
		return errors.New("ParamError error")
	}
	if mailType == uint8(model.System) {
		// 系统邮件
		err := redis_cache.GetCache().GetMailCache().GetSystemMailCache().Set(ctx, userId, mailId, status)
		if err != nil {
			return err
		}
		return nil
	}
	// 设置 redis
	data := &welfare.EmailUser{
		EmailId:   mailId,
		EmailType: mailType,
		Order:     orderId,
		Status:    status,
		SendTime:  sendTime,
	}
	err := redis_cache.GetCache().GetMailCache().GetWelfareMailCache().Set(ctx, userId, mailId, data)
	if err != nil {
		log.Sugar().Errorf("setMailStatus redis email user:%d order:%d mailId:%d status:%d err:%v", userId, orderId, mailId, status, err)
		return err
	}
	// 使用MQ 同步数据
	userEmail := kxmj_report.OrderEmail{
		OrderId:   orderId,
		EmailId:   mailId,
		EmailType: mailType,
		UserId:    userId,
		Status:    status,
		DrawTime:  drawTime,
		CreatedAt: sendTime,
		UpdatedAt: updateTime,
	}
	err = mq.SyncTable(&userEmail, mq.AddOrUpdate)
	if err != nil {
		log.Sugar().Errorf("setMailStatus mqsql email user:%d order:%d mailId:%d status:%d err:%v", userId, orderId, mailId, status, err)
		return err
	}
	return nil
}

// 这个接口只给服务器内部调用 给玩家添加邮件 添加邮件状态 是未读
func AddUserEmail(ctx *gin.Context, userId uint32, configMail *kxmj_core.ConfigEmail, sendTime uint32) error {
	if configMail == nil {
		return errors.New("AddUserEmail configMail = nil")
	}
	if configMail.EmailType == uint8(model.System) {
		err := setMailStatus(ctx, userId, configMail.EmailId, uint8(model.System), uint8(model.UnRead), 0, sendTime, sendTime, 0)
		if err != nil {
			log.Sugar().Errorf("AddUserEmail SetSystemMail error user:%d EmailId:%d err:%v ", userId, configMail.EmailId, err)
			return err
		}
		return nil
	}
	// 创建一个 orderId 用于查询
	orderId := utils.Snowflake.Generate().Int64()
	err := setMailStatus(ctx, userId, configMail.EmailId, configMail.EmailType, uint8(model.UnRead), orderId, sendTime, sendTime, 0)
	if err != nil {
		return err
	}
	return nil
}

// changeSystemMailStatus 改变玩家系统邮件状态
func changeSystemMailStatus(ctx *gin.Context, mailWelfareList map[uint32]uint8, userId uint32, status uint8) {
	if len(mailWelfareList) <= 0 {
		return
	}
	// 设置邮件状态
	updateTime := uint32(time.Now().Unix())
	drawTime := uint32(0)
	for mailId, mailStatus := range mailWelfareList {
		// 设置领取时间
		if status == uint8(model.Take) || (status == uint8(model.Delete) && mailStatus <= uint8(model.Read)) {
			drawTime = updateTime
		}
		// 设置邮件状态 并且 同步到mysql
		err := setMailStatus(ctx, uint32(userId), mailId, uint8(model.System), status, 0, 0, updateTime, drawTime)
		if err != nil {
			log.Sugar().Errorf("SetUserMailRead setMailStatus error user:%d EmailId:%d err:%v", userId, mailId, err)
		}
	}
}

// changeWelfareMailStatus 改变玩家福利邮件状态
func changeWelfareMailStatus(ctx *gin.Context, mailSystemList map[uint32]*welfare.EmailUser, userId uint32, status uint8) {
	if len(mailSystemList) <= 0 {
		return
	}
	// 设置邮件状态
	updateTime := uint32(time.Now().Unix())
	drawTime := uint32(0)
	for _, mailValue := range mailSystemList {
		// 设置领取时间
		if status == uint8(model.Take) || (status == uint8(model.Delete) && mailValue.Status <= uint8(model.Read)) {
			drawTime = updateTime
		}
		// 设置邮件状态 并且 同步到mysql
		err := setMailStatus(ctx, userId, mailValue.EmailId, mailValue.EmailType, status, mailValue.Order, mailValue.SendTime, updateTime, drawTime)
		if err != nil {
			log.Sugar().Errorf("SetUserMailRead setMailStatus error user:%d EmailId:%d err:%v", userId, mailValue.EmailId, err)
		}
	}
}

// getUserAllSystemMail 获取玩家所有系统邮件 邮件的状态 < status
func getUserAllSystemMail(ctx *gin.Context, userId uint32, status uint8, mailSystemList map[uint32]uint8) {
	userAllSystemMail, err := redis_cache.GetCache().GetMailCache().GetSystemMailCache().GetAll(ctx, userId)
	if err != nil {
	}
	for mailId, mailStatus := range userAllSystemMail {
		if status < uint8(model.Delete) && mailId > 0 && mailStatus < status {
			mailSystemList[mailId] = mailStatus
		} else if status == uint8(model.Delete) && mailId > 0 && mailStatus > uint8(model.UnRead) && mailStatus < status {
			mailSystemList[mailId] = mailStatus
		}
	}
}

// getUserAllSystemMail 获取玩家所有福利邮件 邮件的状态 < status
func getUserAllWelfareMail(ctx *gin.Context, userId uint32, status uint8, mailWelfareList map[uint32]*welfare.EmailUser) {
	userAllWelfareMail, err := redis_cache.GetCache().GetMailCache().GetWelfareMailCache().GetAll(ctx, userId)
	if err != nil {
	}
	for _, mailValue := range userAllWelfareMail {
		// 只有未读状态的需要改变状态 因为 2是重复了  3是已领取 4 是删除 所以只有1 才能在这里改变状态
		if status < uint8(model.Delete) && mailValue.EmailId > 0 && mailValue.Status < status {
			mailWelfareList[mailValue.EmailId] = mailValue
		} else if status == uint8(model.Delete) && mailValue.EmailId > 0 && mailValue.Status > uint8(model.UnRead) && mailValue.Status < status {
			mailWelfareList[mailValue.EmailId] = mailValue
		}
	}
}

// getUserSystemMailByMailId 获取玩家所有系统邮件 邮件的状态 < status
func getUserSystemMailByMailId(ctx *gin.Context, userId uint32, status uint8, mailId uint32, mailSystemList map[uint32]uint8) int {
	userSystemMailStatus, err := redis_cache.GetCache().GetMailCache().GetSystemMailCache().Get(ctx, userId, mailId)
	if err != nil {
		return codes.NoMail
	}
	if userSystemMailStatus >= status {
		return 0
	}
	mailSystemList[mailId] = userSystemMailStatus
	return 0
}

// getUserWelfareMailByMailId 获取玩家所有福利邮件 邮件的状态 < status
func getUserWelfareMailByMailId(ctx *gin.Context, userId uint32, status uint8, mailId uint32, mailWelfareList map[uint32]*welfare.EmailUser) int {
	userWelfareMailData, err := redis_cache.GetCache().GetMailCache().GetWelfareMailCache().Get(ctx, userId, mailId)
	if err != nil {
		// 玩家没有邮件 没有邮件则直接返回成功
		return codes.NoMail
	}
	if userWelfareMailData.Status >= status {
		return 0
	}
	mailWelfareList[userWelfareMailData.EmailId] = userWelfareMailData
	return 0
}

// 解析从邮件里面需要领取的物品
func analysisValueItem(needTakeItemMails map[uint32][]*model.MailItem, configItem map[uint32]*kxmj_core.Item) []*item.ValueItem {
	if needTakeItemMails == nil || configItem == nil || len(configItem) <= 0 || len(needTakeItemMails) <= 0 {
		return nil
	}
	valueItemList := make([]*item.ValueItem, 0)
	// 把领取的物品提取出来  map[uint32][]uint32 ---- map[uint32]*kxmj_core.Item
	for mailId, itemList := range needTakeItemMails {
		if mailId <= 0 {
			continue
		}
		for _, itemId := range itemList {
			if itemId.Id <= 0 {
				continue
			}
			dataItem, has := configItem[itemId.Id]
			if has == false {
				continue
			}
			// 获取 物品对应的 valueItem
			valueItem := item.GetValueItem(dataItem)
			// 解析物品基础类型
			values, err := valueItem.ParseBaseValueItems(itemId.Count)
			if err != nil {
				continue
			}
			for _, goldItem := range values {
				valueItemList = append(valueItemList, goldItem)
			}
		}
	}
	if len(valueItemList) <= 0 {
		return nil
	}
	return valueItemList
}

// 把邮件里面的物品解析出来
func analysisItemList(welfareMailList map[uint32]*welfare.EmailUser, systemMailList map[uint32]uint8, mailConfig map[uint32]*kxmj_core.ConfigEmail) map[uint32][]*model.MailItem {
	if len(mailConfig) <= 0 {
		return nil
	}
	needTakeItemMails := make(map[uint32][]*model.MailItem, 0) // 需要领取物品的邮件ID
	// 福利邮件
	for mailId, mailValue := range welfareMailList {
		// value 则是 玩家的mailId
		mailData, has := mailConfig[mailId]
		if has == false {
			continue
		}
		// 已经是领取的 或者删除的 则 跳过
		if mailValue.Status == uint8(model.Take) || mailValue.Status == uint8(model.Delete) || mailData.IsReward != uint8(model.HaveItem) {
			continue
		}
		// 把邮件ID 与 奖励物品 一一对应
		var itemList []*model.MailItem
		err := json.Unmarshal([]byte(mailData.ItemList), &itemList)
		if err != nil || len(itemList) <= 0 {
			continue
		}
		needTakeItemMails[mailId] = itemList
	}
	// 系统邮件
	for mailId, mailStatus := range systemMailList {
		// value 则是 玩家的mailId
		mailData, has := mailConfig[mailId]
		if has == false {
			continue
		}
		// 已经是领取的 或者删除的 则 跳过 邮件配置没有物品
		if mailStatus == uint8(model.Take) || mailStatus == uint8(model.Delete) || mailData.IsReward != uint8(model.HaveItem) {
			continue
		}
		// 把邮件ID 与 奖励物品 一一对应
		var itemList []*model.MailItem
		err := json.Unmarshal([]byte(mailData.ItemList), &itemList)
		if err != nil || len(itemList) <= 0 {
			continue
		}
		needTakeItemMails[mailId] = itemList
	}
	if len(needTakeItemMails) <= 0 {
		return nil
	}
	return needTakeItemMails
}

// 把邮件领取的物品中数量相加 每一种类型 数量相加
func getMailItemGroupMap(needTakeItemMails map[uint32][]*model.MailItem) map[uint32]string {
	if len(needTakeItemMails) <= 0 {
		return nil
	}
	
	itemMap := make(map[uint32]string, 0)
	for _, itemList := range needTakeItemMails {
		for _, item := range itemList {
			itemData, has := itemMap[item.Id]
			if has == false {
				itemMap[item.Id] = item.Count
			} else {
				newItemCount, ok := utils.AddToString(itemData, item.Count)
				if ok {
					itemMap[item.Id] = newItemCount
				}
			}
			
		}
	}
	return itemMap
}

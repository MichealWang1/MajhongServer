package server

import (
	"context"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/redis_cache"
)

// GetRoomConfig 获取游戏房间配置
func (rs *RpcxServer) GetRoomConfig(ctx context.Context, args *center.RoomConfigReq, reply *center.RoomConfigResp) error {
	info, err := redis_cache.GetCache().GetRoomCache().GetDetailCache().Get(ctx, args.GameId, args.RoomId)
	if err != nil {
		reply.Code = codes.GetRoomConfigFailed
		reply.Msg = codes.GetMessage(codes.GetRoomConfigFailed)
		log.Sugar().Errorf("GetRoomConfig game:%d room:%v err:%v", args.GameId, args.RoomId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.RoomConfig{
		RoomId:      info.RoomId,
		RoomType:    info.RoomType,
		GameId:      info.GameId,
		GameType:    info.GameType,
		Tags:        info.Tags,
		Extra:       info.Extra,
		RoomLevel:   info.RoomLevel,
		MinLimit:    info.MinLimit,
		MaxLimit:    info.MaxLimit,
		BaseScore:   info.BaseScore,
		MaxMultiple: info.MaxMultiple,
		Ticket:      info.Ticket,
		MatchTime:   info.MatchTime,
		MatchRobot:  info.MatchRobot,
	}

	return nil
}

// GetRoomConfigList 获取游戏房间配置列表
func (rs *RpcxServer) GetRoomConfigList(ctx context.Context, args *center.RoomConfigListReq, reply *center.RoomConfigListResp) error {
	infos, err := redis_cache.GetCache().GetRoomCache().GetDetailCache().GetAll(ctx, args.GameId)
	if err != nil {
		reply.Code = codes.GetRoomConfigFailed
		reply.Msg = codes.GetMessage(codes.GetRoomConfigFailed)
		log.Sugar().Errorf("GetRoomConfigList game:%d err:%v", args.GameId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	for _, info := range infos {
		reply.Data = append(reply.Data, &center.RoomConfig{
			RoomId:      info.RoomId,
			RoomType:    info.RoomType,
			GameId:      info.GameId,
			GameType:    info.GameType,
			Tags:        info.Tags,
			Extra:       info.Extra,
			RoomLevel:   info.RoomLevel,
			MinLimit:    info.MinLimit,
			MaxLimit:    info.MaxLimit,
			BaseScore:   info.BaseScore,
			MaxMultiple: info.MaxMultiple,
			Ticket:      info.Ticket,
			MatchTime:   info.MatchTime,
			MatchRobot:  info.MatchRobot,
		})
	}
	return nil
}

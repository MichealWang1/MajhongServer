package server

import (
	"context"
	"kxmj.common/codes"
	"kxmj.common/model/lobby"
)

func (rs *RpcxServer) UserLocation(ctx context.Context, args *lobby.LocationReq, reply *lobby.LocationResp) error {
	reply.Code = codes.Success
	if args.UserId <= 0 {
		reply.Code = codes.UserNotExist
	}

	reply.Msg = codes.GetMessage(reply.Code)
	reply.Data = rs.lobby.getLocation(args.UserId)
	return nil
}

func (rs *RpcxServer) GetGateways(ctx context.Context, args *lobby.GetGatewayReq, reply *lobby.GetGatewayResp) error {
	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(reply.Code)
	list := rs.lobby.getGateways()
	for _, endpoint := range list {
		reply.List = append(reply.List, &lobby.GetGatewayInfo{
			SvrType: endpoint.SvrType,
			SvrId:   endpoint.SvrId,
			Addr:    endpoint.Addr,
			Port:    int(endpoint.Port),
		})
	}

	return nil
}

func (rs *RpcxServer) ShoppingSuccess(ctx context.Context, args *lobby.ShoppingSuccessReq, reply *lobby.ShoppingSuccessResp) error {
	rs.lobby.shoppingSuccess(args)
	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(reply.Code)
	return nil
}

func (rs *RpcxServer) GetRoomsOnline(ctx context.Context, args *lobby.GetRoomsOnlineReq, reply *lobby.GetRoomsOnlineResp) error {
	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(reply.Code)
	onlineList := rs.lobby.getRoomsOnline(args.GameId, args.RoomIds)
	for k, v := range onlineList {
		reply.OnlineUsers = append(reply.OnlineUsers, &lobby.RoomOnlineData{
			RoomId: k,
			Users:  v,
		})
	}
	return nil
}

func (rs *RpcxServer) PauseUserGame(ctx context.Context, args *lobby.PauseUserGameReq, reply *lobby.PauseUserGameResp) error {
	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(reply.Code)
	rs.lobby.pauseUserGame(args.UserId)
	return nil
}

func (rs *RpcxServer) ContinueGame(ctx context.Context, args *lobby.ContinueGameReq, reply *lobby.ContinueGameResp) error {
	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(reply.Code)
	rs.lobby.continueUserGame(args.UserId)
	return nil
}

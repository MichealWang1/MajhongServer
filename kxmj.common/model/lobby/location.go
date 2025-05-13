package lobby

type LocationReq struct {
	UserId uint32
}

type LocationResp struct {
	Code int
	Msg  string
	Data *LocationInfo
}

type LocationInfo struct {
	UserId  uint32
	SvrType uint16
	SvrId   uint16
	RoomId  uint32
	DeskId  uint32
}

// GetRoomsOnline

type GetRoomsOnlineReq struct {
	GameId  uint16
	RoomIds []uint32
}

type RoomOnlineData struct {
	RoomId uint32
	Users  uint32
}

type GetRoomsOnlineResp struct {
	Code        int
	Msg         string
	OnlineUsers []*RoomOnlineData
}

// PauseUserGame

type PauseUserGameReq struct {
	UserId uint32
}

type PauseUserGameResp struct {
	Code int
	Msg  string
}

// ContinueGame

type ContinueGameReq struct {
	UserId uint32
}

type ContinueGameResp struct {
	Code int
	Msg  string
}

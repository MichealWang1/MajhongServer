package lobby

type GetGatewayReq struct {
	Args int
}

type GetGatewayResp struct {
	Code int
	Msg  string
	List []*GetGatewayInfo
}

type GetGatewayInfo struct {
	SvrType uint16
	SvrId   uint16
	Addr    string
	Port    int
}

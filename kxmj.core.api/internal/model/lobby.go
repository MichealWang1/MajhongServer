package model

type GetWalletResp struct {
	Diamond  string `json:"diamond"`  // 钻石数
	Gold     string `json:"gold"`     // 金币数
	GoldBean string `json:"goldBean"` // 金豆数
}

package dto

type AddVipBpParameter struct {
	UserId uint32
	BP     uint32
}

type AddVipBpResult struct {
	Level uint32
}

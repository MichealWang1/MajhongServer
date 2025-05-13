package net

type ServerConfig struct {
	Path string `json:"path" yaml:"path"` // rpcx注册服务路径
	Type uint16 `json:"type" yaml:"type"` // 服务器类型
	Id   uint16 `json:"id" yaml:"id"`     // 服务Id
	Ip   string `json:"ip" yaml:"ip"`     // Ip地址
	Port int    `json:"port" yaml:"port"` // 端口号
}

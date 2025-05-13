package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kxmj.common/log"
)

type DbConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Pwd         string `yaml:"pwd"`
	Schema      string `yaml:"schema"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	Name        string `yaml:"name"`
}

var (
	coreMaster *gorm.DB // 业务主库
	coreSlave  *gorm.DB // 业务从库

	loggerMaster *gorm.DB // 业务日志主库
	loggerSlave  *gorm.DB // 业务日志从库

	reportMaster *gorm.DB // 报表主库(business owner)
	reportSlave  *gorm.DB // 报表从库(business owner)

	gameMaster *gorm.DB // 游戏日志主库
	gameSlave  *gorm.DB // 游戏日志从库
)

func create(config *DbConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Pwd, config.Host, config.Port, config.Schema)
	db, err := newDB(dsn, config.MaxIdleConn, config.MaxOpenConn)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("初始化%s失败:", config.Name), zap.Any("db", config.Schema), zap.Any("err", err))
		panic(err)
	}
	log.Sugar().Info(fmt.Sprintf("初始化%s成功:", config.Name), zap.Any("db", config.Schema))
	return db
}

func newDB(dsn string, maxIdleConn int, maxOpenConn int) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger: NewLogger(),
		},
	)

	if err != nil {
		return nil, err
	}

	schema, err := db.DB()
	if err != nil {
		return nil, err
	}

	schema.SetMaxIdleConns(maxIdleConn)
	schema.SetMaxOpenConns(maxOpenConn)
	return db, nil
}

func InitCoreMaster(config *DbConfig) {
	coreMaster = create(config)
}

func InitCoreSlave(config *DbConfig) {
	coreSlave = create(config)
}

func InitLoggerMaster(config *DbConfig) {
	loggerMaster = create(config)
}

func InitLoggerSlave(config *DbConfig) {
	loggerSlave = create(config)
}

func InitReportMaster(config *DbConfig) {
	reportMaster = create(config)
}

func InitReportSlave(config *DbConfig) {
	reportSlave = create(config)
}

func InitGameMaster(config *DbConfig) {
	gameMaster = create(config)
}

func InitGameSlave(config *DbConfig) {
	gameSlave = create(config)
}

func CoreMaster() *gorm.DB {
	return coreMaster
}

func CoreSlave() *gorm.DB {
	return coreSlave
}

func LoggerMaster() *gorm.DB {
	return loggerMaster
}

func LoggerSlave() *gorm.DB {
	return loggerSlave
}

func ReportMaster() *gorm.DB {
	return reportMaster
}

func ReportSlave() *gorm.DB {
	return reportSlave
}

func GameMaster() *gorm.DB {
	return gameMaster
}

func GameSlave() *gorm.DB {
	return gameSlave
}

package database

import (
	"fmt"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/global"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
var(
	db *gorm.DB
	IsConnected bool = true
)
//InitDB init settings of specified db.
func InitDB()error{
	if d, err := gorm.Open(sqlite.Open(global.Sqlite3DBPosition), &gorm.Config{});err!=nil{
		IsConnected = false
		return fmt.Errorf("connect database failed. Some functions are disabled")
	}else{
		db=d
		if db.AutoMigrate(&config.BuildInfo{},&config.TriggerInfo{})!=nil{
			return fmt.Errorf("failed to generate table. Some functions are disabled")
		}
		return nil
	}
}
//Getconn returns db obj. Add debug if specified on start.
func GetConn()*gorm.DB{
	if !global.SqlDebug{
		return db
	}
	return db.Debug()
}
package model

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	Agent struct {
		ID    string `gorm:"column:id"`
		Alias string `gorm:"column:alias"`
		//IP      []string  `gorm:"column:ip;type:blob"`
		Offline bool  `gorm:"column:offline;type:tinyint(1);default:1"`
		Deleted bool  `gorm:"column:deleted;type:tinyint(1);default:0"`
		Enabled bool  `gorm:"column:enabled;type:tinyint(1);default:1"`
		CTime   int64 `gorm:"column:ctime;type:bigint(13)"`
		MTime   int64 `gorm:"column:mtime;type:bigint(13)"`
	}

	Service struct {
		ID     string    `gorm:"column:id"`
		Name   string    `gorm:"column:name;size:45"`
		Remark string    `gorm:"column:remark;size:128"`
		CTime  time.Time `gorm:"column:ctime;type:bigint(13)"`
		MTime  time.Time `gorm:"column:mtime;type:bigint(13)"`
	}
)

func (a *Agent) TableName() string {
	return "agent"
}

func (s *Service) TableName() string {
	return "service"
}

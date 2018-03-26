package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

func updateUnixTimestampOnCreate(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now()

		if ctime, ok := scope.FieldByName("CTime"); ok {
			if ctime.IsBlank {
				ctime.Set(now.UnixNano() / 1e6)
			}
		}

		if mtime, ok := scope.FieldByName("MTime"); ok {
			if mtime.IsBlank {
				mtime.Set(now.UnixNano() / 1e6)
			}
		}
	}
}

func updateUnixTimestampOnUpdate(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("mtime", time.Now().UnixNano()/1e6)
	}
}

func init() {
	gorm.DefaultCallback.Create().Replace("gorm:update_time_stamp", updateUnixTimestampOnCreate)
	gorm.DefaultCallback.Update().Replace("gorm:update_time_stamp", updateUnixTimestampOnUpdate)
}

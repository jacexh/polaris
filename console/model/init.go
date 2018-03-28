package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

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

// softDelete 如果struct包含了Deleted字段，将不会物理删除，而是将Deleted字段置为空
func softDelete(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	var extraOpt string
	if str, ok := scope.Get("gorm:delete_option"); ok {
		extraOpt = fmt.Sprint(str)
	}

	deleted, isSoftDelete := scope.FieldByName("Deleted")

	if !scope.Search.Unscoped && isSoftDelete {
		sql := fmt.Sprintf(
			"UPDATE %v SET `mtime`=%v, %v=%v%v%v",
			scope.QuotedTableName(),
			time.Now().UnixNano()/1e6,
			scope.Quote(deleted.DBName),
			1,
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOpt))
		scope.Raw(sql).Exec()
	} else {
		scope.Raw(fmt.Sprintf(
			"DELETE FROM %v%v%v",
			scope.QuotedTableName(),
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOpt),
		)).Exec()
	}
}

func init() {
	gorm.DefaultCallback.Create().Replace("gorm:update_time_stamp", updateUnixTimestampOnCreate)
	gorm.DefaultCallback.Update().Replace("gorm:update_time_stamp", updateUnixTimestampOnUpdate)

	gorm.DefaultCallback.Delete().Replace("gorm:delete", softDelete)
}

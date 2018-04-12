package wsalemodels

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
)

type Fans struct {
	gorm.Model
	wsalelibs.Fans
	FansExt
}

type FansExt struct {
	UnfollowDate time.Time `json:"unfollow_date"` //取消好友时间（非时时更新，程序逻辑对比结果）
}

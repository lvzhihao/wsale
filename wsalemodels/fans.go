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

//初始化一个实例
func (c *Fans) Ensure(db *gorm.DB, merchantNo, robotWxId, fansWxId string) error {
	unique := &wsalelibs.Fans{
		MerchantNo: merchantNo,
		RobotWxId:  robotWxId,
		FansWxId:   fansWxId,
	}
	return db.Where(Fans{Fans: *unique}).FirstOrInit(c).Error
}

func GetFansByRobot(db *gorm.DB, merchantNo, robotWxId string, finder *Finder) (count int64, rst []*Fans, err error) {
	rst = make([]*Fans, 0)
	if finder != nil {
		db = finder.AddWhere("merchant_no = ?", merchantNo).AddWhere("robot_wx_id = ?", robotWxId).Use(db)
	}
	err = db.Model(&Fans{}).Offset(-1).Limit(-1).Count(&count).Error
	if err == nil {
		err = db.Find(&rst).Error
	}
	return
}

type FansExt struct {
	UnfollowDate time.Time `json:"unfollow_date"` //取消好友时间（非时时更新，程序逻辑对比结果）
}

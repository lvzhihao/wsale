package wsalemodels

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
)

type ChatRoom struct {
	gorm.Model
	wsalelibs.ChatRoom
	ChatRoomExt
}

//初始化一个实例
func (c *ChatRoom) Ensure(db *gorm.DB, merchantNo, robotWxId, chatRoomId string) error {
	unique := &wsalelibs.ChatRoom{
		MerchantNo: merchantNo,
		RobotWxId:  robotWxId,
		ChatRoomId: chatRoomId,
	}
	return db.Where(ChatRoom{ChatRoom: *unique}).FirstOrInit(c).Error
}

func GetChatRoomsByRobot(db *gorm.DB, merchantNo, robotWxId string, finder *Finder) (count int64, rst []*ChatRoom, err error) {
	rst = make([]*ChatRoom, 0)
	db = finder.AddWhere("merchant_no = ?", merchantNo).AddWhere("robot_wx_id = ?", robotWxId).Use(db)
	err = db.Model(&ChatRoom{}).Offset(-1).Limit(-1).Count(&count).Error
	if err == nil {
		err = db.Find(&rst).Error
	}
	return
}

type ChatRoomExt struct {
	//todo
}

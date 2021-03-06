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
	if finder != nil {
		db = finder.AddWhere("merchant_no = ?", merchantNo).AddWhere("robot_wx_id = ?", robotWxId).Use(db)
	}
	err = db.Model(&ChatRoom{}).Offset(-1).Limit(-1).Count(&count).Error
	if err == nil {
		err = db.Find(&rst).Error
	}
	return
}

func UpdateChatRoomMembersCount(db *gorm.DB, chatRoomId string, count int) error {
	return db.Model(&ChatRoom{}).Where("chat_room_id = ?", chatRoomId).Update("members_count", int32(count)).Error
}

type ChatRoomExt struct {
	RobotInStatus bool  `gorm:"default:true" json:"robot_in_status"` // 设备是否在群内
	MembersCount  int32 `gorm:"default:0" json:"members_count"`      // 群成员总数
}

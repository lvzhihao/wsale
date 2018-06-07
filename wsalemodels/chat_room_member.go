package wsalemodels

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
)

type ChatRoomMember struct {
	gorm.Model
	wsalelibs.ChatRoomMember
	ChatRoomMemberExt
}

//初始化一个实例
func (c *ChatRoomMember) Ensure(db *gorm.DB, merchantNo, chatRoomId, fansWxId string) error {
	unique := &wsalelibs.ChatRoomMember{
		MerchantNo: merchantNo,
		ChatRoomId: chatRoomId,
		FansWxId:   fansWxId,
	}
	return db.Where(ChatRoomMember{ChatRoomMember: *unique}).FirstOrInit(c).Error
}

func GetChatRoomMembersByChatRoomId(db *gorm.DB, merchantNo, chatRoomId string, finder *Finder) (count int64, rst []*ChatRoomMember, err error) {
	rst = make([]*ChatRoomMember, 0)
	if finder != nil {
		db = finder.AddWhere("merchant_no = ?", merchantNo).AddWhere("chat_room_id = ?", chatRoomId).Use(db)
	}
	err = db.Model(&ChatRoomMember{}).Offset(-1).Limit(-1).Count(&count).Error
	if err == nil {
		err = db.Find(&rst).Error
	}
	return
}

func EnusreChatRoomMembersCouunt(db *gorm.DB, chatRoomId string) error {
	var count int
	err := db.Model(&ChatRoomMember{}).Where("chat_room_id = ?", chatRoomId).Where("member_in_status = ?", true).Count(&count).Error
	if err != nil {
		return err
	}
	return UpdateChatRoomMembersCount(db, chatRoomId, count)
}

type ChatRoomMemberExt struct {
	MemberInStatus bool      `gorm:"default:true" json:"member_in_status"` // 群用户否在群内
	QuitDate       time.Time `json:"quit_date"`                            // 退群时间
}

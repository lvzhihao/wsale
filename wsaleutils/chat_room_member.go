package wsaleutils

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/goutils"
	"github.com/lvzhihao/wsale/wsalelibs"
	"github.com/lvzhihao/wsale/wsalemodels"
)

func SyncChatRoomMermbers(merchant *wsalelibs.Merchant, chatRoomId string) ([]*wsalelibs.ChatRoomMember, error) {
	list := make([]*wsalelibs.ChatRoomMember, 0)
	var count int32
	var rst []map[string]interface{}
	err := client.M(merchant).SQChatRoomUserList(chatRoomId).ResultKey("nCount", &count).ResultKey("vcList", &rst).Error
	if err != nil {
		return nil, err
	}
	for _, value := range rst {
		member := &wsalelibs.ChatRoomMember{}
		err := member.Unmarshal(merchant.MerchantNo, chatRoomId, value)
		if err != nil {
			return nil, err
		} else {
			list = append(list, member)
		}
	}
	return list, nil
}

func SyncChatRoomMermbersDatabase(db *gorm.DB, merchant *wsalelibs.Merchant, chatRoomId string) ([]*wsalemodels.ChatRoomMember, error) {
	list, err := SyncChatRoomMermbers(merchant, chatRoomId)
	if err != nil {
		return nil, err
	}
	var olds []wsalemodels.ChatRoomMember
	err = db.Where("merchant_no = ?", merchant.MerchantNo).Where("chat_room_id = ?", chatRoomId).Where("member_in_status = ?", true).Find(&olds).Error
	if err != nil {
		return nil, err
	}
	mdls := make([]*wsalemodels.ChatRoomMember, 0)
	ids := make([]string, 0)
	for _, member := range list {
		mdl := &wsalemodels.ChatRoomMember{}
		err := mdl.Ensure(db, member.MerchantNo, member.ChatRoomId, member.FansWxId) //确认数据库记录
		if err != nil {
			return mdls, err
		}
		mdl.ChatRoomMember = *member //更新数据
		mdl.MemberInStatus = true    //设备状态为在群内
		err = db.Save(mdl).Error     //更新记录
		if err != nil {
			return mdls, err
		}
		mdls = append(mdls, mdl)
		ids = append(ids, mdl.FansWxId)
	}
	for _, old := range olds {
		if goutils.InStringSlice(ids, old.FansWxId) {
			continue
		}
		old.MemberInStatus = false
		old.QuitDate = time.Now()
		db.Save(old) //设置状态为不在群内
	}
	go wsalemodels.UpdateChatRoomMembersCount(db, chatRoomId, len(ids)) // 同步人数
	return mdls, nil

}

package wsaleutils

import (
	"encoding/base64"
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
		old.QuitDate.Set(time.Now())
		db.Save(old) //设置状态为不在群内
	}
	go wsalemodels.UpdateChatRoomMembersCount(db, chatRoomId, len(ids)) // 同步人数
	return mdls, nil

}

/*
 * 群好友入群回调及更新数据库
 */
func SyncChatRoomMemberJoinCallback(db *gorm.DB, member *wsalelibs.ChatRoomMember) error {
	// 确认商户
	var merchant wsalemodels.Merchant
	err := db.Where("merchant_no = ?", member.MerchantNo).First(&merchant).Error
	if err != nil {
		return err
	}
	// 初始化群信息
	obj := new(wsalemodels.ChatRoomMember)
	err = obj.Ensure(db, merchant.MerchantNo, member.ChatRoomId, member.FansWxId)
	if err != nil {
		return err
	}
	// update
	obj.NickName = member.NickName
	obj.NickNameBase64 = member.NickNameBase64
	decoded, err := base64.StdEncoding.DecodeString(member.NickNameBase64)
	if err == nil {
		obj.NickName = goutils.ToString(decoded)
	}
	obj.HeadImage = member.HeadImage
	obj.InvitedWxId = member.InvitedWxId
	obj.ChatNickName = member.ChatNickName
	obj.ChatNickNameBase64 = member.ChatNickNameBase64
	decoded, err = base64.StdEncoding.DecodeString(member.ChatNickNameBase64)
	if err == nil {
		obj.ChatNickName = goutils.ToString(decoded)
	}
	obj.MemberInStatus = true    // 成员在群内状态
	obj.JoinDate.Set(time.Now()) // 使用收到入群消息的时间，是一个估值，并不可靠
	return db.Save(obj).Error
}

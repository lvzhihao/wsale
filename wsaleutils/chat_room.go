package wsaleutils

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
	"github.com/lvzhihao/wsale/wsalemodels"
)

/*
 * 同步接口数据
 */
func SyncRobotChatRooms(merchant *wsalelibs.Merchant, robotWxId string) ([]*wsalelibs.ChatRoom, error) {
	page := 1
	size := 20
	list := make([]*wsalelibs.ChatRoom, 0)
	for {
		var count int32
		var rst []map[string]interface{}
		err := client.M(merchant).SQChatRoomList(robotWxId, page, size).ResultKey("nCount", &count).ResultKey("vcList", &rst).Error
		if err != nil {
			return nil, err
		}
		for _, value := range rst {
			chatRoom := &wsalelibs.ChatRoom{}
			err := chatRoom.Unmarshal(value)
			if err != nil {
				return nil, err
			} else {
				//fix merchant no & robot_wx_id
				chatRoom.MerchantNo = merchant.MerchantNo
				chatRoom.RobotWxId = robotWxId
				list = append(list, chatRoom)
			}
		}
		if count <= int32(page*size) {
			break
		}
		page++
	}
	return list, nil
}

/*
 * 同步接口数据入库
 */
func SyncRobotChatRoomsDatabase(db *gorm.DB, merchant *wsalelibs.Merchant, robotWxId string) ([]*wsalemodels.ChatRoom, error) {
	list, err := SyncRobotChatRooms(merchant, robotWxId)
	if err != nil {
		return nil, err
	}
	// todo check robot in status
	mdls := make([]*wsalemodels.ChatRoom, 0)
	for _, chatRoom := range list {
		mdl := &wsalemodels.ChatRoom{}
		err := mdl.Ensure(db, chatRoom.MerchantNo, chatRoom.RobotWxId, chatRoom.ChatRoomId) //确认数据库记录
		if err != nil {
			return mdls, err
		}
		mdl.ChatRoom = *chatRoom //更新数据
		err = db.Save(mdl).Error //更新记录
		if err != nil {
			return mdls, err
		}
		mdls = append(mdls, mdl)
	}
	return mdls, nil
}

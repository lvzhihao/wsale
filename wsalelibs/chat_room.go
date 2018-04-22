package wsalelibs

import (
	"encoding/json"
	"fmt"

	"github.com/lvzhihao/goutils"
)

type ChatRoom struct {
	MerchantNo         string `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_chat_room_id" json:"merchant_no"`  //商户ID
	RobotWxId          string `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_chat_room_id" json:"robot_wx_id"`  //机器人微信ID
	ChatRoomId         string `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_chat_room_id" json:"chat_room_id"` //群号
	ChatRoomName       string `gorm:"type:varchar(100)" json:"chat_room_name"`                                                             //群名称
	ChatRoomNameBase64 string `gorm:"type:varchar(255)" json:"chat_room_name_base64"`                                                      //群名称base64
	RemarkName         string `gorm:"type:varchar(100)" json:"remark_name"`                                                                //备注名
	ServiceWxId        string `gorm:"type:varchar(80)" json:"service_wx_id"`                                                               //消息号ID
	Status             int32  `gorm:"index" json:"status"`                                                                                 //状态
	HeadImage          string `gorm:"type:varchar(500)" json:"head_image"`                                                                 //群头像
	CodeImage          string `gorm:"type:varchar(500)" json:"code_image"`                                                                 //群二维码
}

func (c *ChatRoom) Unmarshal(iter interface{}) error {
	return ChatRoomUnmarshal(iter, c)
}

func ChatRoomUnmarshal(iter interface{}, chatRoom *ChatRoom) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	chatRoomId, ok := m.GetString("vcChatRoomId")
	if !ok {
		return fmt.Errorf("vcChatRoomWxId empty")
	}
	chatRoom.ChatRoomId = chatRoomId
	chatRoom.ChatRoomName, _ = m.GetString("vcName")
	chatRoom.ChatRoomNameBase64, _ = m.GetString("vcBase64Name")
	chatRoom.RemarkName, _ = m.GetString("vcRemarkName")
	chatRoom.ServiceWxId, _ = m.GetString("vcServiceWxId")
	chatRoom.Status, _ = m.GetInt32("nOpenStaus")
	chatRoom.HeadImage, _ = m.GetString("vcHeadImg")
	chatRoom.CodeImage, _ = m.GetString("vcCodeImage")
	// todo
	//chatRoom.FollowDate, _ = xxx
	return nil
}

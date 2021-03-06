package wsalelibs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lvzhihao/goutils"
)

/*
 * 群信息
 */
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
	AdminWxId          string `gorm:"not null;type:varchar(80);index" json:"admin_wx_id"`                                                  //群主微信ID
}

/*
 * 群信息解析
 */
func (c *ChatRoom) Unmarshal(merchant_no, robot_wx_id string, iter interface{}) error {
	return ChatRoomUnmarshal(merchant_no, robot_wx_id, iter, c)
}

/*
 * 群信息解析
 */
func ChatRoomUnmarshal(merchant_no, robot_wx_id string, iter interface{}, chatRoom *ChatRoom) error {
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
	chatRoom.MerchantNo = merchant_no
	chatRoom.RobotWxId = robot_wx_id
	chatRoom.ChatRoomId = chatRoomId
	chatRoom.ChatRoomName, _ = m.GetString("vcName")
	chatRoom.ChatRoomNameBase64, ok = m.GetString("vcBase64Name")
	if ok {
		// 如果有base64并且解码成功，则使用
		dec, err := base64.StdEncoding.DecodeString(chatRoom.ChatRoomNameBase64)
		if err != nil {
			chatRoom.ChatRoomName = goutils.ToString(dec)
		}
	}
	chatRoom.RemarkName, _ = m.GetString("vcRemarkName")
	chatRoom.ServiceWxId, _ = m.GetString("vcServiceWxId")
	chatRoom.Status, _ = m.GetInt32("nOpenStaus")
	chatRoom.HeadImage, _ = m.GetString("vcHeadImg")
	chatRoom.CodeImage, _ = m.GetString("vcCodeImage")
	chatRoom.AdminWxId, _ = m.GetString("vcAdminWxId")
	return nil
}

/*
 * 群变动信息
 */
type ChatRoomModify struct {
	MerchantNo         string `json:"merchant_no"`           //商户ID
	RobotWxId          string `json:"robot_wx_id"`           //机器人微信ID
	ChatRoomId         string `json:"chat_room_id"`          //群号
	ChatRoomName       string `json:"chat_room_name"`        //群名称
	ChatRoomNameBase64 string `json:"chat_room_name_base64"` //群名称base64
	HeadImage          string `json:"head_image"`            //群头像
	AdminWxId          string `json:"admin_wx_id"`           //群主微信ID
}

/*
 * 群变动信息Unmarshal
 */
func (c *ChatRoomModify) Unmarshal(merchant_no, robot_wx_id string, iter interface{}) error {
	return ChatRoomModifyUnmarshal(merchant_no, robot_wx_id, iter, c)
}

/*
 * 群异常信息转换，详见接口文档:群信息异动通知
 */
func ChatRoomModifyUnmarshal(merchant_no, robot_wx_id string, iter interface{}, modify *ChatRoomModify) error {
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
	modify.MerchantNo = merchant_no
	modify.RobotWxId = robot_wx_id
	modify.ChatRoomId = chatRoomId
	modify.ChatRoomName, _ = m.GetString("vcName")
	modify.ChatRoomNameBase64, ok = m.GetString("vcBase64Name")
	if ok {
		// 如果有base64并且解码成功，则使用
		dec, err := base64.StdEncoding.DecodeString(modify.ChatRoomNameBase64)
		if err != nil {
			modify.ChatRoomName = goutils.ToString(dec)
		}
	}
	modify.HeadImage, _ = m.GetString("vcHeadImg")
	modify.AdminWxId, _ = m.GetString("vcAdminWxId")
	return nil
}

/*
 * 群异动回调DataStruct
 * {
 *   "vcMerChantNo": "xxx"
 *   "Data": [
 *     ChatRoomModifyCallbackStruct
 *     ...
 *   ]
 * }
 */
type ChatRoomModifyCallbackStruct struct {
	MerchantNo string        `json:"vcMerChantNo"`
	RobotWxId  string        `json:"vcRobotWxId"`
	Data       []interface{} `json:"vcChatRooms"`
}

/*
 * 群异动回调处理
 */
func ChatRoomModifyCallback(iter interface{}) (ret []ChatRoomModify, err error) {
	ret = make([]ChatRoomModify, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var callback ChatRoomModifyCallbackStruct
		err = json.Unmarshal([]byte(goutils.ToString(data)), &callback)
		if err != nil {
			return
		}
		for _, iter := range callback.Data {
			var obj ChatRoomModify
			err = ChatRoomModifyUnmarshal(callback.MerchantNo, callback.RobotWxId, iter, &obj)
			if err != nil {
				return
			}
			ret = append(ret, obj)
		}
	}
	return
}

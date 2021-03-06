package wsalelibs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lvzhihao/goutils"
)

type Robot struct {
	MerchantNo        string `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id" json:"merchant_no"` //商户ID
	RobotWxId         string `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id" json:"robot_wx_id"` //机器人微信ID
	NickName          string `gorm:"not null;type:varchar(100)" json:"nick_name"`                                           //机器人昵称
	HeadImage         string `gorm:"type:varchar(500)" json:"head_image"`                                                   //机器人头像URL
	WxAlias           string `gorm:"type:varchar(200)" json:"wx_alias"`                                                     //机器备注名
	CodeImage         string `gorm:"type:varchar(500)" json:"code_image"`                                                   //机器人二维码
	WhatsUp           string `gorm:"type:varchar(200)" json:"whats_up"`                                                     //个性签名
	Sex               int32  `gorm:"default:0;index" json:"sex"`                                                            //性别: 0:未定义 1:男 2:女
	Area              string `gorm:"type:varchar(100)" json:"area"`                                                         //地区
	Status            int32  `gorm:"default:12;index" json:"status"`                                                        //状态: 10:在线 12:离线 14:注销
	AutoAllowFan      bool   `gorm:"default:false" json:"auto_allow_fan"`                                                   //是否自动通过好友申请
	AutoAllowChatRoom bool   `gorm:"default:false" json:"auto_allow_chat_room"`                                             //是否自动通过群聊邀请
}

func (c *Robot) Unmarshal(merchantNo string, iter interface{}) error {
	return RobotUnmarshal(merchantNo, iter, c)
}

func RobotUnmarshal(merchantNo string, iter interface{}, robot *Robot) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	id, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	robot.MerchantNo = merchantNo
	robot.RobotWxId = id
	robot.NickName, _ = m.GetString("vcNickName")
	robot.HeadImage, _ = m.GetString("vcHeadImages")
	robot.WxAlias, _ = m.GetString("vcWxAlias")
	robot.CodeImage, _ = m.GetString("vcCodeImages")
	robot.WhatsUp, _ = m.GetString("vcSign")
	robot.Sex, _ = m.GetInt32("nSex")
	robot.Area, _ = m.GetString("vcArea")
	robot.Status, _ = m.GetInt32("nType")
	robot.AutoAllowFan, _ = m.GetBool("nIsAllow")
	robot.AutoAllowChatRoom, _ = m.GetBool("nIsChatRoom")
	return nil
}

func RobotInfoCallback(iter interface{}) (ret []Robot, err error) {
	ret = make([]Robot, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var obj Robot
		err = RobotUnmarshal(rst.MerchantNo, data, &obj)
		if err != nil {
			return
		}
		ret = append(ret, obj)
	}
	return
}

type RobotModifyResult struct {
	MerchantNo string `json:"merchant_no"` //商户ID
	RobotWxId  string `json:"robot_wx_id"` //机器人微信ID
	NickName   string `json:"nick_name"`   //机器人昵称
	HeadImage  string `json:"head_image"`  //机器人头像URL
	WhatsUp    string `json:"whats_up"`    //个性签名
	Sex        int32  `json:"sex"`         //性别: 0:未定义 1:男 2:女
	Area       string `json:"area"`        //地区
}

func (c *RobotModifyResult) Unmarshal(iter interface{}) error {
	return RobotModifyResultUnmarshal(iter, c)
}

func RobotModifyResultUnmarshal(iter interface{}, result *RobotModifyResult) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	id, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	result.RobotWxId = id
	result.NickName, _ = m.GetString("vcNickName")
	base64NickName, ok := m.GetString("vcBase64NickName")
	if ok {
		nickName, err := base64.StdEncoding.DecodeString(base64NickName)
		if err != nil {
			result.NickName = goutils.ToString(nickName)
		}
	}
	result.HeadImage, _ = m.GetString("vcHeadImages")
	result.WhatsUp, _ = m.GetString("vcSign")
	result.Sex, _ = m.GetInt32("nSex")
	result.Area, _ = m.GetString("vcArea")
	return nil
}

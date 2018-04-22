package wsalelibs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lvzhihao/goutils"
)

type Fans struct {
	MerchantNo     string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"merchant_no"` //商户ID
	RobotWxId      string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"robot_wx_id"` //机器人微信ID
	FansWxId       string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"fans_wx_id"`  //好友微信ID
	UserName       string    `gorm:"not null;type:varchar(100)" json:"user_name"`                                                      //好友微信号
	NickName       string    `gorm:"not null;type:varchar(100)" json:"nick_name`                                                       //昵称
	NickNameBase64 string    `gorm:"not null;type:varchar(255)" json:"nick_name_base64"`                                               //昵称base64
	HeadImages     string    `gorm:"type:varchar(500)" json:"head_images"`                                                             //好友头像URL
	WxAlias        string    `gorm:"type:varchar(200)" json:"wx_alias"`                                                                //好友备注名
	WhatsUp        string    `gorm:"type:varchar(200)" json:"whats_up"`                                                                //个性签名
	Sex            int32     `gorm:"default:0;index" json:"sex"`                                                                       //性别: 0:未定义 1:男 2:女
	Proinvice      string    `gorm:"type:varchar(50)" json:"proinvice"`                                                                //所属地省份
	City           string    `gorm:"type:varchar(50)" json:"city"`                                                                     //所属地市区
	WsaleTags      string    `gorm:"type:varchar(200)" json:"wsale_tags"`                                                              //好友标签(后台标签，非微信标签)，多个以,隔开
	FollowDate     time.Time `json:"follow_date"`                                                                                      //添加好友时间（账号离线期间的添加时间无法统计，都归为账号上线时间）
}

func (c *Fans) Unmarshal(iter interface{}) error {
	return FansUnmarshal(iter, c)
}

func FansUnmarshal(iter interface{}, fans *Fans) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	robotWxId, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		return fmt.Errorf("vcFansWxId empty")
	}
	fans.RobotWxId = robotWxId
	fans.FansWxId = fansWxId
	fans.UserName, _ = m.GetString("vcUserName")
	fans.NickName, _ = m.GetString("vcNickName")
	fans.NickNameBase64, _ = m.GetString("vcBase64NickName")
	fans.HeadImages, _ = m.GetString("vcHeadImages")
	fans.WxAlias, _ = m.GetString("vcWxAlias")
	fans.WhatsUp, _ = m.GetString("vcPslSignature")
	fans.Sex, _ = m.GetInt32("nSex")
	fans.Proinvice, _ = m.GetString("vcProinvice")
	fans.City, _ = m.GetString("vcCity")
	fans.WsaleTags, _ = m.GetString("vcTags")
	// todo
	//fans.FollowDate, _ = xxx
	return nil
}

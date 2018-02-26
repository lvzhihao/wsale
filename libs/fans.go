package wsalelibs

import "time"

type Fans struct {
	MerchantNo     string    `gorm:"not null;type:varchar(200);index:idx_merchant_no_robot_wx_id_fans_wx_id" json:"merchant_no"` //商户ID
	RobotWxId      string    `gorm:"not null;type:varchar(200);index:idx_merchant_no_robot_wx_id_fans_wx_id" json:"robot_wx_id"` //机器人微信ID
	FansWxId       string    `gorm:"not null;type:varchar(200);index:idx_merchant_no_robot_wx_id_fans_wx_id" json:"fans_wx_id"`  //好友微信ID
	UserName       string    `gorm:"not null;type:varchar(200)" json:"user_name"`                                                //好友微信号
	NickName       string    `gorm:"not null;type:varchar(200)" json:"nick_name`                                                 //昵称
	NickNameBase64 string    `gorm:"no null;type:varchar(500)" json:"nick_name_base64"`                                          //昵称base64
	HeadImages     string    `gorm:"type:varchar(500)" json:"head_images"`                                                       //好友头像URL
	WxAlias        string    `gorm:"type:varchar(200)" json:"wx_alias"`                                                          //好友备注名
	WhatsUp        string    `gorm:"type:varchar(255)" json:"whats_up"`                                                          //个性签名
	Sex            int32     `gorm:"index:idx_sex" json:"sex"`                                                                   //性别: 0:未定义 1:男 2:女
	Proinvice      string    `gorm:"type:varchar(50)" json:"proinvice"`                                                          //所属地省份
	City           string    `gorm:"type:varchar(50)" json:"city"`                                                               //所属地市区
	WsaleTags      string    `gorm:"type:varchar(200)" json:"wsale_tags"`                                                        //好友标签(后台标签，非微信标签)，多个以,隔开
	FollowDate     time.Time `json:"follow_date"`                                                                                //添加好友时间（账号离线期间的添加时间无法统计，都归为账号上线时间）
}

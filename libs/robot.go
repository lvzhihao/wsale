package wsalelibs

type Robot struct {
	MerchantNo string `gorm:"not null;type:varchar(200);unique_index:uix_merchant_no_robot_wx_id" json:"merchant_no"` //商户ID
	RobotWxId  string `gorm:"not null;type:varchar(200);unique_index:uix_merchant_no_robot_wx_id" json:"robot_wx_id"` //机器人微信ID
	NickName   string `gorm:"not null;type:varchar(200)" json:"nick_name"`                                            //机器人昵称
	HeadImage  string `gorm:"type:varchar(500)" json:"head_image"`                                                    //机器人头像URL
	WxAlias    string `gorm:"type:varchar(200)" json:"wx_alias"`                                                      //机器备注名
	CodeImage  string `gorm:"type:varhcar(500)" json:"code_image"`                                                    //机器人二维码
	WhatsUp    string `gorm:"type:varchar(255)" json:"whats_up"`                                                      //个性签名
	Sex        int32  `gorm:"index:idx_sex" json:"sex"`                                                               //性别: 0:未定义 1:男 2:女
	Status     int32  `gorm:"index:idx_status" json:"status"`                                                         //状态: 10:在线 12:离线 14:注销
}

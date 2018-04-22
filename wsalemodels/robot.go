package wsalemodels

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
)

type Robot struct {
	gorm.Model
	wsalelibs.Robot
	RobotExt
}

//初始化一个实例
func (c *Robot) Ensure(db *gorm.DB, merchantNo, robotWxId string) error {
	unique := &wsalelibs.Robot{
		MerchantNo: merchantNo,
		RobotWxId:  robotWxId,
	}
	return db.Where(Robot{Robot: *unique}).FirstOrInit(c).Error
}

func GetRobotsByMerchant(db *gorm.DB, merchantNo string) (rst []*Robot, err error) {
	err = db.Where("merchant_no = ?", merchantNo).Find(&rst).Error
	return
}

type RobotExt struct {
	FansTotal     int32  `gorm:"default:0;index" json:"fans_total"`      //粉丝数量
	ChatRoomTotal int32  `gorm:"default:0;index" json:"chat_room_total"` //聊天群数量
	CircleImage   string `gorm:"type:varchar(500)" json:"circle_image"`  //朋友圈封面图
}

func (c *RobotExt) ToggleAutoAllow() error {
	//todo
	return nil
}

func (c *RobotExt) OpenAutoAllow() error {
	//todo
	return nil
}

func (c *RobotExt) CloseAutoAllow() error {
	//todo
	return nil
}

func (c *RobotExt) UpdateFansTotal() error {
	//todo
	return nil
}

func (c *RobotExt) UpdateChatRoomTotal() error {
	//todo
	return nil
}

func (c *RobotExt) SetCircleImage(imgUrl string) error {
	//todo
	return nil
}

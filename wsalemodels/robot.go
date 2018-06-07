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

func GetRobotsByMerchant(db *gorm.DB, merchantNo string, finder *Finder) (rst []*Robot, err error) {
	if finder != nil {
		err = finder.Use(db).Where("merchant_no = ?", merchantNo).Find(&rst).Error
	} else {
		err = db.Where("merchant_no = ?", merchantNo).Find(&rst).Error
	}
	return
}

func UpdateRobotFansTotal(db *gorm.DB, robotWxId string, total int) error {
	return db.Model(&Robot{}).Where("robot_wx_id = ?", robotWxId).Update("fans_total", int32(total)).Error
}

func UpdateRobotChatRoomTotal(db *gorm.DB, robotWxId string, total int) error {
	return db.Model(&Robot{}).Where("robot_wx_id = ?", robotWxId).Update("chat_room_total", int32(total)).Error
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

func (c *RobotExt) UpdateFansTotal(db *gorm.DB, total int) error {
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

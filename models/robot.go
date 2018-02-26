package wsalemodels

import (
	"github.com/jinzhu/gorm"
	wsalelibs "github.com/lvzhihao/wsale/libs"
)

type Robot struct {
	gorm.Model
	wsalelibs.Robot
	RbotoExt
}

//初始化一个实例
func (c *Robot) Ensure() error {
	//todo
	return nil
}

type RobotExt struct {
	AutoAllow     bool   `json:"auto_allow"`                                       //是否自动通过好友申请
	FansTotal     int32  `gorm:"index:idx_fans_total" json:"fans_total"`           //粉丝数量
	ChatRoomTotal int32  `gorm:"index:idx_chat_room_total" json:"chat_room_total"` //聊天群数量
	CircleImage   string `gorm:"type:varchar(500)" json:"circle_image"`            //朋友圈封面图
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

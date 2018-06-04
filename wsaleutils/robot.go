package wsaleutils

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
	"github.com/lvzhihao/wsale/wsalemodels"
)

/*
 * 同步个人号上下号及注销状态回调并更新数据库，可重复调用
 */
func SyncRobotInfoCallback(db *gorm.DB, robot *wsalelibs.Robot) error {
	// find exists record
	rst := new(wsalemodels.Robot)
	err := db.Where("merchant_no = ?", robot.MerchantNo).Where("robot_wx_id = ?", robot.RobotWxId).First(rst).Error
	if err != nil {
		return err
	}
	// 这个回调其它数据不可靠，只更新上下号状态
	rst.Status = robot.Status
	// update
	return db.Save(rst).Error
}

/*
 * 同步个人号修改信息回并更新数据库，可重复调用
 */
func SyncRobotModifyResultCallback(db *gorm.DB, result *wsalelibs.RobotModifyResult) error {
	// find exists record
	rst := new(wsalemodels.Robot)
	err := db.Where("merchant_no = ?", result.MerchantNo).Where("robot_wx_id = ?", result.RobotWxId).First(rst).Error
	if err != nil {
		return err
	}
	// 这个回调更新 昵称 头像 性别 签名 地区
	rst.NickName = result.NickName
	rst.HeadImage = result.HeadImage
	rst.Sex = result.Sex
	rst.WhatsUp = result.WhatsUp
	rst.Area = result.Area
	// update
	return db.Save(rst).Error
}

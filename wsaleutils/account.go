package wsaleutils

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
	"github.com/lvzhihao/wsale/wsalemodels"
)

var (
	client *wsalelibs.Client
)

func init() {
	client = wsalelibs.NewClient()
}

/*
 * 同步接口数据
 */
func SyncAccountRobots(iter interface{}) ([]*wsalelibs.Robot, error) {
	var rst []map[string]interface{}
	var merchant wsalelibs.Merchant
	err := client.M(iter).GetAccountList().CurrentM(&merchant).Result(&rst).Error
	if err != nil {
		return nil, err
	}
	robots := make([]*wsalelibs.Robot, 0)
	for _, value := range rst {
		robot, err := wsalelibs.RobotConvert(value)
		if err != nil {
			return nil, err
		} else {
			robot.MerchantNo = merchant.MerchantNo
			robots = append(robots, robot)
		}
	}
	return robots, nil
}

/*
 * 同步接口数据入库
 */
func SyncAccountRobotsDatabase(iter interface{}, db *gorm.DB) ([]*wsalemodels.Robot, error) {
	robots, err := SyncAccountRobots(iter)
	if err != nil {
		return nil, err
	}
	mdls := make([]*wsalemodels.Robot, 0)
	for _, robot := range robots {
		mdl := &wsalemodels.Robot{}
		err := mdl.Ensure(db, robot) //确认数据库记录
		if err != nil {
			return mdls, err
		}
		err = db.Save(&mdl).Error //更新记录
		if err != nil {
			return mdls, err
		}
		mdls = append(mdls, mdl)
	}
	return mdls, nil
}

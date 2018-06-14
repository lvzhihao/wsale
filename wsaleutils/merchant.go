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
func SyncMerchantRobots(merchant *wsalelibs.Merchant) ([]*wsalelibs.Robot, error) {
	var rst []map[string]interface{}
	err := client.M(merchant).GetAccountList().Result(&rst).Error
	if err != nil {
		return nil, err
	}
	robots := make([]*wsalelibs.Robot, 0)
	for _, value := range rst {
		robot := &wsalelibs.Robot{}
		err := robot.Unmarshal(merchant.MerchantNo, value)
		if err != nil {
			return nil, err
		} else {
			robots = append(robots, robot)
		}
	}
	return robots, nil
}

/*
 * 同步接口数据入库
 */
func SyncMerchantRobotsDatabase(db *gorm.DB, merchant *wsalelibs.Merchant) ([]*wsalemodels.Robot, error) {
	robots, err := SyncMerchantRobots(merchant)
	if err != nil {
		return nil, err
	}
	mdls := make([]*wsalemodels.Robot, 0)
	for _, robot := range robots {
		mdl := &wsalemodels.Robot{}
		err := mdl.Ensure(db, robot.MerchantNo, robot.RobotWxId) //确认数据库记录
		if err != nil {
			return mdls, err
		}
		mdl.Robot = *robot
		err = db.Save(mdl).Error //更新记录
		if err != nil {
			return mdls, err
		}
		mdls = append(mdls, mdl)
	}
	return mdls, nil
}

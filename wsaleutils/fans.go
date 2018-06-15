package wsaleutils

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
	"github.com/lvzhihao/wsale/wsalemodels"
)

/*
 * 同步接口数据
 */
func SyncRobotFans(merchant *wsalelibs.Merchant, robotWxId string) ([]*wsalelibs.Fans, error) {
	page := 1
	size := 50
	list := make([]*wsalelibs.Fans, 0)
	for {
		var count int32
		var rst []map[string]interface{}
		err := client.M(merchant).GetFriendsListByRobot(robotWxId, page, size).ResultKey("nCount", &count).ResultKey("vcList", &rst).Error
		if err != nil {
			return nil, err
		}
		for _, value := range rst {
			fans := &wsalelibs.Fans{}
			err := fans.Unmarshal(merchant.MerchantNo, value)
			if err != nil {
				return nil, err
			} else {
				//fix merchant no
				list = append(list, fans)
			}
		}
		if count <= int32(page*size) {
			break
		}
		page++
	}
	return list, nil
}

/*
 * 同步接口数据入库
 */
func SyncRobotFansDatabase(db *gorm.DB, merchant *wsalelibs.Merchant, robotWxId string) ([]*wsalemodels.Fans, error) {
	list, err := SyncRobotFans(merchant, robotWxId)
	if err != nil {
		return nil, err
	}
	mdls := make([]*wsalemodels.Fans, 0)
	for _, fans := range list {
		mdl := &wsalemodels.Fans{}
		err := mdl.Ensure(db, fans.MerchantNo, fans.RobotWxId, fans.FansWxId) //确认数据库记录
		if err != nil {
			return mdls, err
		}
		mdl.Fans = *fans
		err = db.Save(mdl).Error //更新记录
		if err != nil {
			return mdls, err
		}
		mdls = append(mdls, mdl)
	}
	go wsalemodels.UpdateRobotFansTotal(db, robotWxId, len(mdls))
	return mdls, nil
}

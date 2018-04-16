package wsalemodels

import (
	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/wsale/wsalelibs"
)

type Merchant struct {
	gorm.Model
	wsalelibs.Merchant
	MerchantExt
}

type MerchantExt struct {
	//todo
}

func LoadEnabledMerchants(db *gorm.DB) (ret []Merchant, err error) {
	err = db.Where("is_enabled = ?", true).Find(&ret).Error
	return
}

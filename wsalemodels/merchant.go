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

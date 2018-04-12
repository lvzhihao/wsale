package wsalemodels

import (
	"github.com/jinzhu/gorm"
	wsalelibs "github.com/lvzhihao/wsale/libs"
)

type Merchant struct {
	gorm.Model
	wsalelibs.Merchant
	MerchantExt
}

type MerchantExt struct {
	Name      string `gorm:"type:varchar(100)" json:"name"`   //商户名称
	IsEnabled bool   `gorm:"default:false" json:"is_enabled"` //是否可用
}

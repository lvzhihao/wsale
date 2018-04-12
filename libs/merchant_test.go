package wsalelibs

import (
	"testing"
)

func Test_001_Merchant_001_maps(t *testing.T) {
	EnableMerchant(testMerchant)
	v, err := LoadMerchant(testMerchant.MerchantNo)
	if err != nil {
		t.Error(err)
	} else {
		DisableMerchant(testMerchant)
		_, err := LoadMerchant(testMerchant.MerchantNo)
		if err == nil {
			t.Error("merchant map disable error")
		} else {
			t.Logf("merchant map test success: %+v", v)
		}
	}
}

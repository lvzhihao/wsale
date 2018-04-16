package wsaleutils

import (
	"flag"
	"testing"

	"github.com/lvzhihao/wsale/wsalelibs"
)

var (
	testMerchantNo     string
	testMerchantSecret string
	testMerchant       *wsalelibs.Merchant
)

func init() {
	flag.StringVar(&testMerchantNo, "no", "", "merchantNo")
	flag.StringVar(&testMerchantSecret, "secret", "", "merchantSecret")
	flag.Parse()
	testMerchant = &wsalelibs.Merchant{
		MerchantNo:     testMerchantNo,
		MerchantSecret: testMerchantSecret,
	}
}

func Test_001_Account_SyncAccountRobots(t *testing.T) {
	rst, err := SyncAccountRobots(testMerchant)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
}

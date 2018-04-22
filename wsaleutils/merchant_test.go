package wsaleutils

import (
	"flag"
	"testing"

	"github.com/lvzhihao/wsale/wsalelibs"
)

var (
	testMerchantNo     string
	testMerchantSecret string
	testRobotWxId      string
	testMerchant       *wsalelibs.Merchant
)

func init() {
	flag.StringVar(&testMerchantNo, "no", "", "merchantNo")
	flag.StringVar(&testMerchantSecret, "secret", "", "merchantSecret")
	flag.StringVar(&testRobotWxId, "robot", "", "robotWxId")
	flag.Parse()
	testMerchant = &wsalelibs.Merchant{
		MerchantNo:     testMerchantNo,
		MerchantSecret: testMerchantSecret,
	}
}

func Test_001_Merchant_SyncMerchantRobots(t *testing.T) {
	rst, err := SyncMerchantRobots(testMerchant)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
}

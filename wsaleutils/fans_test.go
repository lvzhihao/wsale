package wsaleutils

import (
	"testing"
)

func Test_002_Fans_SyncRobotsFans(t *testing.T) {
	rst, err := SyncRobotFans(testMerchant, testRobotWxId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
}

package wsaleutils

import (
	"testing"

	"github.com/lvzhihao/goutils"
)

func Test_002_Fans_SyncRobotsChatRooms(t *testing.T) {
	rst, err := SyncRobotChatRooms(testMerchant, testRobotWxId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(goutils.ToString(rst))
}

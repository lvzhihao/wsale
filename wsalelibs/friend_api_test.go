package wsalelibs

import (
	"testing"

	"github.com/lvzhihao/goutils"
)

func Test_004_Friend_001_SQRobotFriendsCirclePush(t *testing.T) {
	err := testClient.M(testMerchant).SQRobotFriendsCirclePush(goutils.ToString(testAccount["vcRobotWxId"])).Error
	if err != nil {
		t.Fatal(err)
	}
}

func Test_004_Friend_002_SQRobotFriendsCircleGet(t *testing.T) {
	var rst []map[string]interface{}
	err := testClient.M(testMerchant).SQRobotFriendsCircleGet(goutils.ToString(testAccount["vcRobotWxId"])).Result(&rst).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
}

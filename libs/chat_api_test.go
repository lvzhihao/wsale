package wsalelibs

import (
	"testing"

	"github.com/lvzhihao/goutils"
)

func Test_006_api_001_SQChatRoomList(t *testing.T) {
	var list []map[string]interface{}
	var pages int32
	err := testClient.M(testMerchant).SQChatRoomList(goutils.ToString(testAccount["vcRobotWxId"])).ResultKey("vcList", &list).ResultKey("nCount", &pages).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list, pages)
}

func Test_006_api_001_SQChatRoomUserList(t *testing.T) {
	var list []map[string]interface{}
	var count int32
	//err := testClient.M(testMerchant).SQChatRoomUserList("5396087388@chatroom").ResultKey("vcList", &list).ResultKey("nCount", &count).Error
	err := testClient.M(testMerchant).SQChatRoomUserList("5551031436@chatroom").ResultKey("vcList", &list).ResultKey("nCount", &count).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list, count)
}

package wsalelibs

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/lvzhihao/goutils"
)

var (
	testImgMsgId string
)

func Test_003_api_001_SendMessage(t *testing.T) {
	return
	send := &SendMessage{}
	send.Sender(goutils.ToString(testAccount["vcRobotWxId"]))
	send.Fans(testAccountFansWxId)
	send.AddText("hello")
	send.AddImage("http://www.m555.com/mb_pic/2007/09/20070917093919_6a0709.jpg")
	send.AddLink("title", "http://www.baidu.com", "http://www.m555.com/mb_pic/2007/09/20070917093919_6a0709.jpg", "baidu")
	err := testClient.M(testMerchant).SendMessage(send).Error
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(send)
	t.Logf("%s\n", b)
	send.ChatRoom(testChatRoomId)
	err = testClient.M(testMerchant).SendMessage(send).Error
	if err != nil {
		t.Fatal(err)
	}
	b, _ = json.Marshal(send)
	t.Logf("%s\n", b)
	//send.ChatRoomAtFans(testChatRoomId, []string{testAccountFansWxId})
	send.ChatRoomAtAll(testChatRoomId)
	err = testClient.M(testMerchant).SendMessage(send).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range send.Data {
		if data.MsgType == MSG_TYPE_IMAGE {
			testImgMsgId = data.MsgId
			break
		}
	}
	b, _ = json.Marshal(send)
	t.Logf("%s\n", b)
}

func Test_003_api_002_GetOriginalImg(t *testing.T) {
	var bigImg string
	err := testClient.M(testMerchant).GetOriginalImg(testImgMsgId).ResultKey("vcMessage", &bigImg).Error
	//err := testClient.M(testMerchant).GetOriginalImg("394015433146314452").ResultKey("vcMessage", &bigImg).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bigImg)
}

func Test_003_api_003_GetMiniAppsList(t *testing.T) {
	var count int32
	var rst []map[string]interface{}
	err := testClient.M(testMerchant).GetMiniAppsList().ResultKey("vcList", &rst).ResultKey("nCount", &count).Error
	//err := testClient.M(testMerchant).GetOriginalImg("394015433146314452").ResultKey("vcMessage", &bigImg).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst, count)
}

func Test_003_api_004_SendMiniApps(t *testing.T) {
	send := &SendMessage{}
	send.
		Sender(goutils.ToString(testAccount["vcRobotWxId"])).
		Fans(testAccountFansWxId).
		AddText("world").
		AddMini("wx7c8d593b2c3a7703")
	err := testClient.M(testMerchant).SendMiniApps(send).Error
	log.Fatal(send.Format())
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(send)
	t.Logf("%s\n", b)
}

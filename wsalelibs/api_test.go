package wsalelibs

import (
	"flag"
	"strings"
	"testing"

	"github.com/lvzhihao/goutils"
)

var (
	testMerchantNo      string
	testMerchantSecret  string
	testAccountId       string
	testMerchant        Merchant
	testClient          *Client
	testAccount         map[string]interface{}
	testAccountFans     []map[string]interface{}
	testAccountFansWxId string
	testChatRoomId      string
)

func init() {
	flag.StringVar(&testMerchantNo, "no", "", "merchantNo")
	flag.StringVar(&testMerchantSecret, "secret", "", "merchantSecret")
	flag.StringVar(&testAccountId, "id", "", "accountId")
	flag.StringVar(&testAccountFansWxId, "fansid", "", "fansAccountWxId")
	flag.StringVar(&testChatRoomId, "room", "", "chatroomId")
	flag.Parse()
	testMerchant = Merchant{
		MerchantNo:     testMerchantNo,
		MerchantSecret: testMerchantSecret,
	}
	testClient = NewClient()
}

func Test_002_api_001_GetAccountList(t *testing.T) {
	var rst []map[string]interface{}
	err := testClient.M(testMerchant).GetAccountList().Result(&rst).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
	for _, account := range rst {
		if strings.Compare(goutils.ToString(account["vcRobotWxId"]), testAccountId) == 0 {
			testAccount = account
			break
		}
	}
}

func Test_002_api_002_GetAccount(t *testing.T) {
	var rst map[string]interface{}
	err := testClient.M(testMerchant).GetAccount(goutils.ToString(testAccount["vcRobotWxId"])).First(&rst).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rst)
}

func Test_002_api_003_GetFriendsListByRobot(t *testing.T) {
	var list []map[string]interface{}
	var pages int32
	err := testClient.M(testMerchant).GetFriendsListByRobot(goutils.ToString(testAccount["vcRobotWxId"])).ResultKey("vcList", &list).ResultKey("nCount", &pages).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list, pages)
	testAccountFans = list
}

func Test_002_api_004_GetRobotFansTags(t *testing.T) {
	var fansWxId []string
	for _, fan := range testAccountFans {
		fansWxId = append(fansWxId, goutils.ToString(fan["vcFansWxId"]))
	}
	err := testClient.M(testMerchant).GetRobotFansTags(goutils.ToString(testAccount["vcRobotWxId"]), fansWxId).Error
	if err != nil {
		t.Fatal(err)
	}
}

func Test_002_api_005_SQScanWeiXinLogin(t *testing.T) {
	err := testClient.M(testMerchant).SQScanWeiXinLogin().Error
	if err != nil {
		t.Fatal(err)
	}
}

func Test_002_api_006_SQAccountWeiXinAfreshLogin(t *testing.T) {
	err := testClient.M(testMerchant).SQAccountWeiXinAfreshLogin(goutils.ToString(testAccount["vcRobotWxId"])).Error
	if err != nil {
		t.Fatal(err)
	}
}

func Test_002_api_007_SQScanWeiXinAfreshLogin(t *testing.T) {
	var serialNo string
	err := testClient.M(testMerchant).SQScanWeiXinAfreshLogin(goutils.ToString(testAccount["vcRobotWxId"])).ResultKey("vcMessage", &serialNo).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(serialNo)
}

func Test_002_api_008_SQRobotOffline(t *testing.T) {
	//todo 登录测试用例
	return
	err := testClient.M(testMerchant).SQRobotOffline(goutils.ToString(testAccount["vcRobotWxId"])).Error
	if err != nil {
		t.Fatal(err)
	}
}

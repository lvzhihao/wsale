package wsalelibs

import (
	"strings"

	"github.com/lvzhihao/goutils"
)

/*
 获取个人号好友列表接口
 robotWxId string 个人号ID
 input[0] int 前当页，默认1
 input[1] int 每页大小，默认20
 input[2] string 模糊查询关键字（微信群名称）
*/
func (c *Client) SQChatRoomList(robotWxId string, input ...interface{}) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["nPageIndex"] = 1
	params["nPageSize"] = 20
	if len(input) >= 1 && goutils.ToInt32(input[0]) > 0 {
		params["nPageIndex"] = goutils.ToInt32(input[0])
	}
	if len(input) >= 2 && goutils.ToInt32(input[1]) > 0 {
		params["nPageSize"] = goutils.ToInt32(input[1])
	}
	if len(input) >= 3 && goutils.ToString(input[2]) != "" {
		params["vcName"] = goutils.ToString(input[2])
	}
	return c.Chat("SQChatRoomList", params)
}

/*
 群好友列表查询接口
 chatRoomId string 群ID
*/
func (c *Client) SQChatRoomUserList(chatRoomId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcChatRoomId"] = chatRoomId
	return c.Chat("SQChatRoomUserList", params)
}

/*
 好友入群邀请接口
 ChatRoomId string 群ID
 RobotWxId string 设备ID
 FansWxIdList []string 好友ID
*/
func (c *Client) SQChatRoomRobotJoinRoom(chatRoomId, robotWxId string, fansWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcChatRoomId"] = chatRoomId
	params["vcRobotWxId"] = robotWxId
	params["vcFans"] = strings.Join(fansWxIdList, ",")
	return c.Chat("SQChatRoomRobotJoinRoom", params)
}

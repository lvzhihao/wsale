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

/*
 群开启或关闭状态设置接口
 robotWxId string 设备号
 chatRoomIdList []string 群ID
 status int32 状态
*/
func (c *Client) SQChatRoomSet(robotWxId string, chatRoomIdList []string, status int32) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["aryId"] = chatRoomIdList
	params["nOpenStaus"] = goutils.ToString(status)
	return c.Chat("SQChatRoomSet", params)

}

/*
 群消息号操作接口
 robotWxId string 设备号
 chatRoomIdList []string 群ID
 ntype int32 状态
*/
func (c *Client) SQChatRoomServiceSet(robotWxId string, chatRoomIdList []string, ntype int32) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["aryId"] = chatRoomIdList
	params["nType"] = goutils.ToString(ntype)
	return c.Chat("SQChatRoomServiceSet", params)
}

/*
 请求修改群名称接口
 robotWxId string 设备号
 chatRoomId string 群ID
 name string 群名
*/
func (c *Client) SQChatRoomUpdateName(robotWxId, chatRoomId, name string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["vcChatRoomId"] = chatRoomId
	params["vcName"] = name
	return c.Chat("SQChatRoomUpdateName", params)
}

/*
 请求修改群名称接口
 robotWxId string 设备号
 chatRoomId string 群ID
 name string 群名
*/
func (c *Client) SQUpdateChatRoomRemarkName(robotWxId, chatRoomId, remark string) (client *Client) {
	// 文档的接口没有指定哪个设备下的群备注，此接口暂时不要用
	return c
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["vcChatRoomId"] = chatRoomId
	params["RemarkNamei"] = remark
	return c.Chat("SQUpdateChatRoomRemarkName", params)
}

/*
 设置是否自动同意邀请入群接口
 robotWxIdList []string 设备号
 allow bool 是否自动同意
*/

func (c *Client) SQRobotUpdateIsChatroom(robotWxIdList []string, allow bool) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["tbRobotList"] = robotWxIdList
	if allow {
		params["isChatroom"] = "1"
	} else {
		params["isChatroom"] = "0"
	}
	return c.Chat("SQRobotUpdateIsChatroom", params)
}

/*
 请求退群接口
 robotWxId string 设备号
 chatRoomIdList []string 群ID
*/
func (c *Client) SQChatRoomSetExit(robotWxId string, chatRoomIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["aryChatRoomId"] = chatRoomIdList
	return c.Chat("SQChatRoomSetExit", params)
}

package wsalelibs

import (
	"github.com/lvzhihao/goutils"
)

func (c *Client) SQFriendsDynamicInsert(send *SendMoment) (client *Client) {
	return c
}

func (c *Client) SQRobotFriendsCirclePush(robotWxId string, input ...[]string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	if len(input) >= 1 {
		params["vcFriendsNo"] = goutils.ToString(input[0])
	}
	return c.Friend("SQRobotFriendsCirclePush", params)
}

func (c *Client) SQRobotFriendsCircleGet(robotWxId string, input ...[]string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	if len(input) >= 1 {
		params["vcFriendsNo"] = goutils.ToString(input[0])
	}
	return c.Friend("SQRobotFriendsCircleGet", params)
}

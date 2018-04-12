package wsalelibs

/*
 发送消息接口
 sendMessage *SendMessage 消息体
*/
func (c *Client) SendMessage(sendMessage *SendMessage) (client *Client) {
	client = c.clone()
	msg, err := sendMessage.M(c.Merchant).Format()
	if err != nil {
		client.AddError(err)
	} else {
		client = c.Msg("SendMessage", msg)
	}
	return
}

/*
 获取原图接口
 msgId string 消息唯一ID
*/
func (c *Client) GetOriginalImg(msgId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcMsgId"] = msgId
	return c.Msg("GetOriginalImg", params)
}

/*
 获取小程序列表接口
 robotWxId string 个人号ID
 input[0] int 前当页，默认1
 input[1] int 每页大小，默认20
*/
func (c *Client) GetMiniAppsList(input ...int) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["nPageIndex"] = 1
	params["nPageSize"] = 20
	if len(input) >= 1 && input[0] > 0 {
		params["nPageIndex"] = input[0]
	}
	if len(input) >= 2 && input[1] > 0 {
		params["nPageSize"] = input[1]
	}
	return c.Msg("GetMiniAppsList", params)
}

/*
 发送小程序接口
 sendMessage *SendMessage 消息体
*/
func (c *Client) SendMiniApps(sendMessage *SendMessage) (client *Client) {
	client = c.clone()
	msg, err := sendMessage.M(c.Merchant).Format()
	if err != nil {
		client.AddError(err)
	} else {
		client = c.Msg("SendMiniApps", msg)
	}
	return
}

package wsalelibs

/*
 获取个人号列表接口
 robotWxId string 个人号ID，不传时查询商家下所有个人号
*/
func (c *Client) getAccountList(robotWxId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	if robotWxId != "" {
		params["vcRobotWxId"] = robotWxId // 如果批定微信号
	}
	return c.Robot("GetAccountList", params)
}

/*
 获取个人号列表接口
*/
func (c *Client) GetAccountList() (client *Client) {
	return c.getAccountList("")
}

/*
 获取个人号信息
 robotWxId string 个人号ID
*/
func (c *Client) GetAccount(robotWxId string) (client *Client) {
	return c.getAccountList(robotWxId)
}

/*
 发起个人号扫描登录接口
*/
func (c *Client) SQScanWeiXinLogin() (client *Client) {
	return c.Robot("SQScanWeiXinLogin", nil)
}

/*
 扫码登陆下线后请求重新登陆接口
 robotWxId string 个人号ID
*/
func (c *Client) SQScanWeiXinAfreshLogin(robotWxId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	return c.Robot("SQScanWeiXinAfreshLogin", params)
}

/*
 帐号密码登陆下线后请求重新登录接口
 robotWxId string 个人号ID
*/
func (c *Client) SQAccountWeiXinAfreshLogin(robotWxId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	return c.Robot("SQAccountWeiXinAfreshLogin", params)
}

/*
 个人号下线提交接口
 robotWxId string 个人号ID
*/
func (c *Client) SQRobotOffline(robotWxId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	return c.Robot("SQRobotOffline", params)
}

/*
 个人号修改昵称提交接口
 nickName string 个人号昵称
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateNickName(nickName string, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcNickName"] = nickName
	params["vcJson"] = robotWxIdList
	return c.Robot("SQRobotUpdateNickName", params)
}

/*
 个人号修改头像提交接
 headImage string 个人号昵称
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateHeadImage(headImage string, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcHeadImage"] = headImage
	params["vcJson"] = robotWxIdList
	return c.Robot("SQRobotUpdateHeadImage", params)
}

/*
 个人号修改性别提交接口
 sex int 性别（1男，2女）
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateSex(sex int, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["nSex"] = sex
	params["vcJson"] = robotWxIdList
	return c.Robot("SQRobotUpdateSex", params)
}

/*
 个人号修改个性签名提交接口
 sign string 个性签名
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateSign(sign string, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcSign"] = sign
	params["vcJson"] = robotWxIdList
	return c.Robot("SQRobotUpdateSign", params)
}

/*
 个人号修改地区提交接口
 area string 所在地区（例如：CN_Hunan_Changsha 湖南-长沙）
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateArea(area string, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcArea"] = area
	params["vcJson"] = robotWxIdList
	return c.Robot("SQRobotUpdateArea", params)
}

/*
 设置是否自动通过好友申请接口
 isAllow bool 是否自动同意邀请入群，1是（开启），0否（关闭
 robotWxIdList []string 个号人ID列表
*/
func (c *Client) SQRobotUpdateIsAllow(isAllow bool, robotWxIdList []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	if isAllow {
		params["isAllow"] = "1"
	} else {
		params["isAllow"] = "0"
	}
	params["tbRobotList"] = robotWxIdList
	return c.Robot("SQRobotUpdateIsAllow", params)
}

/*
 获取线下标签
 可通过该接口获取手机上个人号的标签。(通过回调地址返回数据)
 robotWxId string 个人号ID
 fansData []string 好友微信ID
*/
func (c *Client) GetRobotFansTag(robotWxId string, FansWxId []string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["vcReturnUrl"] = "apiRequest"
	params["fansData"] = FansWxId
	return c.Robot("GetRobotFansTags", params)
}

/*
 获取个人号好友列表接口
 robotWxId string 个人号ID
 input[0] int 前当页，默认1
 input[1] int 每页大小，默认20
*/
func (c *Client) GetFriendsListByRobot(robotWxId string, input ...int) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["nPageIndex"] = 1
	params["nPageSize"] = 20
	if len(input) >= 1 && input[0] > 0 {
		params["nPageIndex"] = input[0]
	}
	if len(input) >= 2 && input[1] > 0 {
		params["nPageSize"] = input[1]
	}
	return c.Robot("GetFriendsListByRobot", params)
}

/*
 好友打标签接口
 robotWxId string 个人号ID
 fansWxId string 好友微信ID
 tag string 标签名称
*/
func (c *Client) InsertFriendsTag(robotWxId, fansWxId, tag string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["vcFansWxId"] = fansWxId
	params["vcTag"] = tag
	return c.Robot("InsertFriendsTag", params)
}

/*
 发送同意陌生人加我的请求
 robotWxId string 个人号ID
 fansWxId string 好友微信ID
*/
func (c *Client) PushAgreeFriendsRequest(robotWxId, fansWxId string) (client *Client) {
	params := make(map[string]interface{}, 0)
	params["vcRobotWxId"] = robotWxId
	params["vcFansWxId"] = fansWxId
	return c.Robot("PushAgreeFriendsRequest", params)
}

package wsalelibs

import (
	"fmt"
	"strings"

	"github.com/lvzhihao/goutils"
)

const (
	MSG_TYPE_TEXT  string = "2001" //文字
	MSG_TYPE_IMAGE string = "2002" //图片
	MSG_TYPE_VOICE string = "2003" //语音
	MSG_TYPE_LINK  string = "2005" //链接
	MSG_TYPE_CARD  string = "2006" //名片
	MSG_TYPE_MINI  string = "2013" //小程序
)

type SendMessage struct {
	Merchant   *Merchant      `json:"merchant"`     //商家编号
	RobotWxId  string         `json:"robot_wx_id"`  //个人号微信ID
	NickName   string         `json:"nick_name"`    //个人号昵称
	Source     int8           `json:"source"`       //10：私聊 11：群聊
	Data       []*MessageData `json:"data"`         //发送消息内容
	ChatRoomId string         `json:"chat_room_id"` //群ID
	FansWxId   []string       `json:"fans_wx_id"`   //好友ID，群聊时为@好友列表
	AgreeAtAll bool           `json:"agree_at_all"` //是否@所有人
}

func (c *SendMessage) M(m interface{}) *SendMessage {
	switch m.(type) {
	case *Merchant:
		c.Merchant = m.(*Merchant)
	case string:
		c.Merchant, _ = LoadMerchant(m.(string))
	default:
		c.Merchant = nil
	}
	return c
}

func (c *SendMessage) Sender(robotWxId string, input ...[]string) *SendMessage {
	c.RobotWxId = robotWxId
	if len(input) == 1 {
		c.NickName = goutils.ToString(input[0])
	}
	return c
}

func (c *SendMessage) ChatRoom(chatRoomId string) *SendMessage {
	c.Source = 11
	c.ChatRoomId = chatRoomId
	c.FansWxId = make([]string, 0)
	return c
}

func (c *SendMessage) ChatRoomAtAll(chatRoomId string) *SendMessage {
	c.ChatRoom(chatRoomId)
	c.AgreeAtAll = true
	return c
}

func (c *SendMessage) ChatRoomAtFans(chatRoomId string, fansWxId []string) *SendMessage {
	c.ChatRoom(chatRoomId)
	c.AgreeAtAll = false
	c.FansWxId = fansWxId
	return c
}

func (c *SendMessage) Fans(fansWxId string) *SendMessage {
	c.Source = 10
	c.FansWxId = []string{fansWxId}
	c.ChatRoomId = ""
	return c
}

func (c *SendMessage) AddText(text string) (send *SendMessage) {
	data := &MessageData{
		MsgType: MSG_TYPE_TEXT,
		Content: text,
	}
	return c.AddData(data)
}

func (c *SendMessage) AddImage(imgUrl string) *SendMessage {
	data := &MessageData{
		MsgType: MSG_TYPE_IMAGE,
		Image:   imgUrl,
	}
	return c.AddData(data)
}

func (c *SendMessage) AddVoice(voiceUrl string, duration int) *SendMessage {
	data := &MessageData{
		MsgType:       MSG_TYPE_VOICE,
		Link:          voiceUrl,
		MediaDuration: int8(duration),
	}
	return c.AddData(data)
}

func (c *SendMessage) AddLink(link, img, desc string) *SendMessage {
	data := &MessageData{
		MsgType: MSG_TYPE_LINK,
		Link:    link,
		Image:   img,
		Desc:    desc,
	}
	return c.AddData(data)
}

func (c *SendMessage) AddCard() *SendMessage {
	//todo
	return c
}

func (c *SendMessage) AddMini(appId string) *SendMessage {
	data := &MessageData{
		MsgType: MSG_TYPE_MINI,
		Link:    appId,
	}
	return c.AddData(data)
}

func (c *SendMessage) AddData(data *MessageData) *SendMessage {
	data.MsgId = goutils.RandStr(20)
	c.Data = append(c.Data, data)
	return c
}

// 小程序只可以单独发送
func (c *SendMessage) Format() (rst map[string]interface{}, err error) {
	if c.Merchant == nil {
		err = fmt.Errorf("Merchant don't exists")
		return
	}
	if c.RobotWxId == "" {
		err = fmt.Errorf("empty sender with robot")
		return
	}
	if !(c.Source == 10 || c.Source == 11) {
		err = fmt.Errorf("error sender source")
		return
	}
	if c.Data == nil {
		err = fmt.Errorf("empty data")
		return
	}
	err = nil
	rst = make(map[string]interface{}, 0)
	rst["vcMerChantNo"] = c.Merchant.MerchantNo
	rst["vcNickName"] = c.NickName
	rst["vcRobotWxID"] = c.RobotWxId
	rst["inSource"] = goutils.ToString(c.Source)
	if c.Source == 10 {
		rst["vcCustomerWxID"] = c.FansWxId[0]
	} else {
		rst["vcChatRoomID"] = c.ChatRoomId
		if c.AgreeAtAll == true {
			rst["at"] = "@"
		} else if len(c.FansWxId) > 0 {
			rst["at"] = strings.Join(c.FansWxId, ",")
		}
	}
	rst["data"] = make([]interface{}, 0)
	for _, data := range c.Data {
		if data.MsgType == MSG_TYPE_MINI {
			// 如果有小程序，则只发送小程序，兼容接口业务方包装
			delete(rst, "data")
			rst["vcAppId"] = data.Link
			break
		} else {
			rst["data"] = append(rst["data"].([]interface{}), data.Format())
		}
	}
	return
}

type ReceiveMessage struct {
	Merchant string                `json:"merchant_no"`
	Data     []*ReceiveMessageData `json:"data"`
}

type ReceiveMessageData struct {
	MessageData
	SendType   string `json:"send_type"`    //消息类型（10:用户消息）
	Source     string `json:"source"`       //10：私聊 11：群聊
	RobotWxId  string `json:"robot_wx_id"`  //个人号微信ID
	ChatRoomId string `json:"chat_room_id"` //群ID
	FansWxId   string `json:"fans_wx_id"`   //好友微信ID
	FansType   string `json:"fans_type"`    //好友微信类型（10:用户消息）
}

type MessageData struct {
	MsgId         string `json:"msg_id"`         //消息唯一ID
	MsgType       string `json:"msg_type"`       //2001:文字 2002:图片 2003:语音 2005:链接 2006:名片
	Content       string `json:"content"`        //消息内容、链接标题
	Image         string `json:"image"`          //图片地址，链接图片
	Link          string `json:"link"`           //链接地址，语音地址
	Desc          string `json:"desc"`           //链接描述
	MediaDuration int8   `json:"meida_duration"` //资源长度，单位秒
}

func (c *MessageData) Format() (rst map[string]interface{}) {
	rst = make(map[string]interface{}, 0)
	rst["vcMsgId"] = c.MsgId
	rst["snMsgType"] = c.MsgType
	rst["vcContent"] = c.Content
	rst["vcUrl"] = c.Link
	rst["vcImg"] = c.Image
	rst["vcLinkDesc"] = c.Desc
	rst["inDuration"] = int(c.MediaDuration)
	return
}

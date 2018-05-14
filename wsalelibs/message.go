package wsalelibs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lvzhihao/goutils"
)

const (
	MSG_TYPE_TEXT  int32 = 2001 //文字
	MSG_TYPE_IMAGE int32 = 2002 //图片
	MSG_TYPE_VOICE int32 = 2003 //语音
	MSG_TYPE_LINK  int32 = 2005 //链接
	MSG_TYPE_CARD  int32 = 2006 //名片
	MSG_TYPE_MINI  int32 = 2013 //小程序
)

type SendMessage struct {
	Merchant   *Merchant      `json:"merchant"`     //商家编号
	RobotWxId  string         `json:"robot_wx_id"`  //个人号微信ID
	NickName   string         `json:"nick_name"`    //个人号昵称
	Source     int32          `json:"source"`       //10：私聊 11：群聊
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
		MsgType:       MSG_TYPE_TEXT,
		Content:       text,
		ContentBase64: base64.StdEncoding.EncodeToString([]byte(text)),
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

func (c *SendMessage) AddCard(cardInfo string) *SendMessage {
	data := &MessageData{
		MsgType: MSG_TYPE_CARD,
		Content: cardInfo,
	}
	return c.AddData(data)
}

// 发送小程序改为了单独接口，此接口不推荐使用
func (c *SendMessage) AddMini(appId string) *SendMessage {
	data := &MessageData{
		MsgType: MSG_TYPE_MINI,
		Content: appId,
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
			rst["vcAppId"] = data.Content
			break
		} else {
			rst["data"] = append(rst["data"].([]interface{}), data.Format())
		}
	}
	return
}

type MessageData struct {
	MsgId         string `json:"msg_id"`         //消息唯一ID
	MsgType       int32  `json:"msg_type"`       //2001:文字 2002:图片 2003:语音 2005:链接 2006:名片
	Content       string `json:"content"`        //消息内容、链接标题
	ContentBase64 string `json:"content_base64"` //消息内容Base64
	Image         string `json:"image"`          //图片地址，链接图片
	Link          string `json:"link"`           //链接地址，语音地址
	Desc          string `json:"desc"`           //链接描述
	MediaDuration int8   `json:"media_duration"` //资源长度，单位秒
}

func (c *MessageData) Format() (rst map[string]interface{}) {
	rst = make(map[string]interface{}, 0)
	rst["vcMsgId"] = c.MsgId
	rst["snMsgType"] = c.MsgType
	rst["vcContent"] = c.Content
	if c.MsgType == MSG_TYPE_TEXT {
		rst["vcBase64Content"] = c.ContentBase64
	} // 消息内容，base64编码（仅用于snMsgType=2001）详见文档
	rst["vcUrl"] = c.Link
	rst["vcImg"] = c.Image
	rst["vcLinkDesc"] = c.Desc
	rst["inDuration"] = int(c.MediaDuration)
	return
}

func (c *MessageData) Unmarshal(iter interface{}) error {
	return MessageDataUnmarshal(iter, c)
}

func MessageDataUnmarshal(iter interface{}, msg *MessageData) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	id, ok := m.GetString("vcMsgId")
	if !ok {
		return fmt.Errorf("vcMsgId empty")
	}
	msg.MsgId = id
	msg.MsgType, _ = m.GetInt32("snMsgType")
	msg.ContentBase64, _ = m.GetString("vcContent")
	if content, err := base64.StdEncoding.DecodeString(msg.ContentBase64); err == nil {
		msg.Content = goutils.ToString(content)
	} // save origin text
	msg.Image, _ = m.GetString("vcImg")
	msg.Link, _ = m.GetString("vcUrl")
	msg.Desc, _ = m.GetString("vcLinkDesc")
	duration, _ := m.GetInt32("inDuration")
	msg.MediaDuration = int8(duration)
	return nil
}

type ReceiveMessage struct {
	MerchantNo string    `json:"merchant_no"`  //商户号
	SendType   int32     `json:"send_type"`    //消息类型（10:用户消息）
	MsgDate    time.Time `json:"msg_date"`     //消息时间
	Source     int32     `json:"source"`       //10：私聊 11：群聊
	RobotWxId  string    `json:"robot_wx_id"`  //个人号微信ID
	ChatRoomId string    `json:"chat_room_id"` //群ID
	FansWxId   string    `json:"fans_wx_id"`   //好友微信ID
	FansType   int32     `json:"fans_type"`    //好友微信类型（10:用户消息）
	MessageData
}

func (c *ReceiveMessage) Unmarshal(iter interface{}) error {
	return ReceiveMessageUnmarshal(iter, c)
}

func ReceiveMessageUnmarshal(iter interface{}, msg *ReceiveMessage) error {
	err := msg.MessageData.Unmarshal(iter)
	if err != nil {
		return err
	}
	var input map[string]interface{}
	err = json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	msg.SendType, _ = m.GetInt32("snSendType")
	msg.Source, _ = m.GetInt32("inSource")
	msg.RobotWxId, _ = m.GetString("vcRobotWxID")
	msg.ChatRoomId, _ = m.GetString("vcChatRoomID")
	msg.FansWxId, _ = m.GetString("vcCustomerWxID")
	msg.FansType, _ = m.GetInt32("nCustomerType")
	if t, ok := m.GetString("dtMsgTime"); ok {
		msg.MsgDate, _ = time.ParseInLocation("2006/1/2 15:04:05", t, TimeLocation)
	}
	return nil
}

type SendMessageResult struct {
	MerchantNo string    `gorm:"not null;type:varchar(80);index" json:"merchant_no"` //商户ID
	MsgId      string    `gorm:"not null;type:varchar(100);index" json:"msg_id"`     //消息唯一ID
	WxMsgId    string    `gorm:"not null;type:varchar(100);index" json:"wx_msg_id"`  //微信消息ID
	Type       int32     `json:"type"`                                               //结果 10:成功 11:失败
	MsgDate    time.Time `json:"msg_date"`                                           //消息发送时间
}

func (c *SendMessageResult) Unmarshal(iter interface{}) error {
	return SendMessageResultUnmarshal(iter, c)
}

func SendMessageResultUnmarshal(iter interface{}, rst *SendMessageResult) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	rst.MsgId, _ = m.GetString("vcMsgId")
	rst.WxMsgId, _ = m.GetString("vcWxMsgId")
	rst.Type, _ = m.GetInt32("inType")
	if t, ok := m.GetInt64("dtMsgTime"); ok {
		rst.MsgDate = time.Unix(t, 0).In(TimeLocation)
	}
	return nil
}

type SendMiniApp struct {
	SendMessage
	MsgId      string `json:"msg_id"`      // 消息唯一ID
	AppId      string `json:"app_id"`      // 小程序appID
	WxGZH      string `json:"wx_gzh"`      // 小程序公众号
	AppTitle   string `json:"app_title"`   // 小程序标题
	AppName    string `json:"app_name"`    // 小程序名称
	SmallImage string `json:"small_image"` // 小程序小图
	BigImage   string `json:"big_image"`   // 小程序大图
	PagePath   string `json:"page_path"`   // 小程序页面路径
}

// todo other set
func (c *SendMiniApp) SetAppId(appId string) *SendMiniApp {
	c.MsgId = goutils.RandStr(20)
	c.AppId = appId
	return c
}

func (c *SendMiniApp) FormatMiniApp() (rst map[string]interface{}, err error) {
	c.Data = make([]*MessageData, 0)
	rst, err = c.SendMessage.Format()
	if err != nil {
		return
	}
	delete(rst, "data")
	rst["vcMsgId"] = c.MsgId
	rst["vcAppId"] = c.AppId
	rst["vcWxGzh"] = c.WxGZH
	rst["vcAppTitle"] = c.AppTitle
	rst["vcAppName"] = c.AppName
	rst["vcSmallImg"] = c.SmallImage
	rst["vcBigImg"] = c.BigImage
	rst["vcPagePath"] = c.PagePath
	return rst, nil
}

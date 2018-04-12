package wsalelibs

import (
	"time"
)

const (
	MMT_TYPE_TEXT  string = "2001" //文字
	MMT_TYPE_IMAGE string = "2002" //图片
	MMT_TYPE_LINK  string = "2005" //链接
)

type SendMoment struct {
	Merchant  *Merchant     `json:"merchant"` //商家编号
	Display   bool          `json:"display"`
	TaskTime  time.Time     `json:"task_time"`
	Data      []*MomentData `json:"content"`
	RobotList []string      `json:"robot_list"`
	FansList  []string      `json:"fans_list"`
}

func (c *SendMoment) M(m interface{}) *SendMoment {
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

func (c *SendMoment) Fans(fansList []string) *SendMoment {
	c.FansList = fansList
	return c
}

func (c *SendMoment) AddText(text string) (send *SendMoment) {
	data := &MomentData{
		MsgType: MMT_TYPE_TEXT,
		Content: text,
	}
	return c.AddData(data)
}

func (c *SendMoment) AddImage(imgUrl string) *SendMoment {
	data := &MomentData{
		MsgType: MMT_TYPE_IMAGE,
		Content: imgUrl,
	}
	return c.AddData(data)
}

func (c *SendMoment) AddLink(link, img, desc string) *SendMoment {
	data := &MomentData{
		MsgType: MMT_TYPE_LINK,
		Link:    link,
		Content: img,
		Desc:    desc,
	}
	return c.AddData(data)
}

func (c *SendMoment) AddData(data *MomentData) *SendMoment {
	c.Data = append(c.Data, data)
	return c
}

// 小程序只可以单独发送
func (c *SendMoment) Format() (rst map[string]interface{}, err error) {
	/*
		if c.Merchant == nil {
			err = fmt.Errorf("Merchant don't exists")
			return
		}
		if len(c.RobotList) == 0 {
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
	*/
	return
}

type MomentData struct {
	MsgType string `json:"msg_type"` //2001:文字 2002:图片 2005:链接
	Content string `json:"content"`  //内容
	Title   string `json:"title"`    //标题
	Desc    string `json:"desc"`     //链接描述
	Link    string `json:"link"`     //链接地址，语音地址
}

func (c *MomentData) Format() (rst map[string]interface{}) {
	rst = make(map[string]interface{}, 0)
	rst["nMsgType"] = c.MsgType
	rst["vcContent"] = c.Content
	rst["vcShareTitle"] = c.Title
	rst["vcShareDesc"] = c.Desc
	rst["vcShareHref"] = c.Link
	return
}

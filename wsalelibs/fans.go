package wsalelibs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lvzhihao/goutils"
)

type Fans struct {
	MerchantNo     string   `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"merchant_no"` //商户ID
	RobotWxId      string   `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"robot_wx_id"` //机器人微信ID
	FansWxId       string   `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_robot_wx_id_fans_wx_id" json:"fans_wx_id"`  //好友微信ID
	UserName       string   `gorm:"not null;type:varchar(100)" json:"user_name"`                                                      //好友微信号
	NickName       string   `gorm:"not null;type:varchar(100)" json:"nick_name`                                                       //昵称
	NickNameBase64 string   `gorm:"not null;type:varchar(255)" json:"nick_name_base64"`                                               //昵称base64
	HeadImages     string   `gorm:"type:varchar(500)" json:"head_images"`                                                             //好友头像URL
	WxAlias        string   `gorm:"type:varchar(200)" json:"wx_alias"`                                                                //好友备注名
	WhatsUp        string   `gorm:"type:varchar(200)" json:"whats_up"`                                                                //个性签名
	Sex            int32    `gorm:"default:0;index" json:"sex"`                                                                       //性别: 0:未定义 1:男 2:女
	Proinvice      string   `gorm:"type:varchar(50)" json:"proinvice"`                                                                //所属地省份
	City           string   `gorm:"type:varchar(50)" json:"city"`                                                                     //所属地市区
	WsaleTags      string   `gorm:"type:varchar(200)" json:"wsale_tags"`                                                              //好友标签(后台标签，非微信标签)，多个以,隔开
	FollowDate     NullTime `gorm:"default:NULL" json:"follow_date"`                                                                  //添加好友时间（账号离线期间的添加时间无法统计，都归为账号上线时间）
}

func (c *Fans) Unmarshal(iter interface{}) error {
	return FansUnmarshal(iter, c)
}

func FansUnmarshal(iter interface{}, fans *Fans) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	robotWxId, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		return fmt.Errorf("vcFansWxId empty")
	}
	fans.RobotWxId = robotWxId
	fans.FansWxId = fansWxId
	fans.UserName, _ = m.GetString("vcUserName")
	fans.NickName, _ = m.GetString("vcNickName")
	if fans.NickNameBase64, ok = m.GetString("vcBase64NickName"); ok {
		// preferred to use base64decode
		decode, err := base64.StdEncoding.DecodeString(fans.NickNameBase64)
		if err == nil {
			fans.NickName = goutils.ToString(decode)
		}
	}
	fans.HeadImages, _ = m.GetString("vcHeadImages")
	fans.WxAlias, _ = m.GetString("vcWxAlias")
	fans.WhatsUp, _ = m.GetString("vcPslSignature")
	fans.Sex, _ = m.GetInt32("nSex")
	fans.Proinvice, _ = m.GetString("vcProinvice")
	fans.City, _ = m.GetString("vcCity")
	fans.WsaleTags, _ = m.GetString("vcTags")
	if t, ok := m.GetString("dtCreateDate"); ok {
		if followDate, err := time.ParseInLocation("2006-01-02T15:04:05.999", t, TimeLocation); err == nil {
			fans.FollowDate.Time = followDate
			fans.FollowDate.Valid = true
		}
	}

	return nil
}

type FansTags struct {
	MerchantNo string   `json:"merchant_no"` //商户ID
	RobotWxId  string   `json:"robot_wx_id"` //机器人微信ID
	FansWxId   string   `json:"fans_wx_id"`  //好友微信ID
	Tags       []string `json:"tags"`        //线下标签
}

func FansTagsMapUnmarshal(iter interface{}) (map[string]*FansTags, error) {
	fansTagsList := make(map[string]*FansTags, 0)
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return nil, err
	}
	m := goutils.NewMap(input)
	robotWxId, ok := m.GetString("vcRobotWxId")
	if !ok {
		return nil, fmt.Errorf("vcRobotWxId empty")
	}
	data, ok := m.GetSlice("tagData")
	if !ok {
		return nil, fmt.Errorf("tagData empty")
	}
	for _, tag := range data {
		var info map[string]interface{}
		err := json.Unmarshal([]byte(goutils.ToString(tag)), &info)
		if err != nil {
			continue
		}
		m := goutils.NewMap(info)
		fansWxId, _ := m.GetString("vcFansWxId")
		tagName, _ := m.GetString("vcTagName")
		index := robotWxId + ">>" + fansWxId
		obj, ok := fansTagsList[index]
		if ok {
			obj.Tags = append(obj.Tags, tagName)
		} else {
			fansTagsList[index] = &FansTags{
				RobotWxId: robotWxId,
				FansWxId:  fansWxId,
				Tags:      []string{tagName},
			}
		}
	}
	return fansTagsList, nil
}

type FansInvite struct {
	MerchantNo     string    `gorm:"not null;type:varchar(80);index" json:"merchant_no"` //商户ID
	RobotWxId      string    `gorm:"not null;type:varchar(80);index" json:"robot_wx_id"` //机器人微信ID
	FansWxId       string    `gorm:"not null;type:varchar(80);index "json:"fans_wx_id"`  //好友微信ID
	NickName       string    `gorm:"not null;type:varchar(100)" json:"nick_name`         //昵称
	HeadImage      string    `gorm:"type:varchar(500); json:"head_image"`                //关像
	RequestText    string    `gorm:"type:varchar(200)" json:"request_text"`              //验证语
	RequestDate    time.Time `gorm:"default:NULL" json:"request_date"`                   //请求时间
	Source         string    `gorm:"type:varchar(50)" json:"source"`                     //来源
	RequestPackage string    `gorm:"type:text" json:"request_package"`                   //请求包，同意加好友使用
	FansStatus     int32     `json:"fans_status"`                                        // 1 已经添加 2 未同意 3 同意失败
}

func (c *FansInvite) Unmarshal(iter interface{}) error {
	return FansInviteUnmarshal(iter, c)
}

func FansInviteUnmarshal(iter interface{}, invite *FansInvite) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	robotWxId, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		return fmt.Errorf("vcFansWxId empty")
	}
	invite.RobotWxId = robotWxId
	invite.FansWxId = fansWxId
	invite.NickName, _ = m.GetString("vcNickName")
	invite.HeadImage, _ = m.GetString("vcHeadImage")
	invite.RequestText, _ = m.GetString("vcRequestText")
	invite.Source, _ = m.GetString("vcSource")
	invite.RequestPackage, _ = m.GetString("vcRequestPackage")
	invite.FansStatus, _ = m.GetInt32("nIsFans")
	if t, ok := m.GetString("dtSendDate"); ok {
		invite.RequestDate, _ = time.ParseInLocation("2006-01-02T15:04:05", t, TimeLocation)
	}
	return nil
}

type FansAgreeResult struct {
	MerchantNo string `gorm:"not null;type:varchar(80);index" json:"merchant_no"` //商户ID
	RobotWxId  string `gorm:"not null;type:varchar(80);index" json:"robot_wx_id"` //机器人微信ID
	FansWxId   string `gorm:"not null;type:varchar(80);index "json:"fans_wx_id"`  //好友微信ID
	Result     int32  `json:"result"`                                             //请求结果 1 添加成功 0 添加失败
}

func (c *FansAgreeResult) Unmarshal(merchantNo string, iter interface{}) error {
	return FansAgreeResultUnmarshal(merchantNo, iter, c)
}

func FansAgreeResultUnmarshal(merchantNo string, iter interface{}, result *FansAgreeResult) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	robotWxId, ok := m.GetString("vcRobotWxId")
	if !ok {
		return fmt.Errorf("vcRobotWxId empty")
	}
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		return fmt.Errorf("vcFansWxId empty")
	}
	result.MerchantNo = merchantNo
	result.RobotWxId = robotWxId
	result.FansWxId = fansWxId
	result.Result, _ = m.GetInt32("nResult")
	return nil
}

func FansAgreeResultCallback(iter interface{}) (ret []FansAgreeResult, err error) {
	ret = make([]FansAgreeResult, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var obj FansAgreeResult
		err = FansAgreeResultUnmarshal(rst.MerchantNo, data, &obj)
		if err != nil {
			return
		}
		ret = append(ret, obj)
	}
	return
}

/*
 * 好友来源及状态，异步回调获取
 */
type FansSourceStatus struct {
	MerchantNo string `json:"merchant_no"` //商户ID
	RobotWxId  string `json:"robot_wx_id"` //机器人微信ID
	FansWxId   string `json:"fans_wx_id"`  //好友微信ID
	Source     string `json:"source"`      //来源
	IsDeleted  bool   `json:"is_deleted"`  //好友是否已经删除个人号
}

/*
 * 回调解码，单条数据，详见个人号开发文档
 */
func (c *FansSourceStatus) Unmarshal(merchantNo, robotWxId string, iter interface{}) error {
	return FansSourceStatusUnmarshal(merchantNo, robotWxId, iter, c)
}

func FansSourceStatusUnmarshal(merchantNo, robotWxId string, iter interface{}, result *FansSourceStatus) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	result.MerchantNo = merchantNo
	result.RobotWxId = robotWxId
	m := goutils.NewMap(input)
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		return fmt.Errorf("vcFansWxId empty")
	}
	result.FansWxId = fansWxId
	result.Source, _ = m.GetString("vcSource")
	result.IsDeleted, _ = m.GetBool("nIsDelFans")
	return nil
}

type FansSourceStatusCallbackStruct struct {
	RobotWxId string        `json:"vcRobotWxId"`
	Data      []interface{} `json:"vcSource"`
}

/*
 * 回调解析
 */
func FansSourceStatusCallback(iter interface{}) (ret []*FansSourceStatus, err error) {
	ret = make([]*FansSourceStatus, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var callback FansSourceStatusCallbackStruct
		err = json.Unmarshal([]byte(goutils.ToString(data)), &callback)
		if err != nil {
			return
		}
		for _, iter := range callback.Data {
			obj := new(FansSourceStatus)
			err = FansSourceStatusUnmarshal(rst.MerchantNo, callback.RobotWxId, iter, obj)
			if err != nil {
				return
			}
			ret = append(ret, obj)
		}
	}
	return
}

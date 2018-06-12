package wsalelibs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lvzhihao/goutils"
)

type ChatRoomMember struct {
	MerchantNo         string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_chat_room_id_fans_wx_id" json:"merchant_no"`  //商户ID
	ChatRoomId         string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_chat_room_id_fans_wx_id" json:"chat_room_id"` //群号
	FansWxId           string    `gorm:"not null;type:varchar(80);unique_index:uix_merchant_no_chat_room_id_fans_wx_id" json:"fans_wx_id"`   //群用户ID
	NickName           string    `gorm:"not null;type:varchar(100)" json:"nick_name"`                                                        //群用户昵称
	NickNameBase64     string    `gorm:"type:varchar(250)" json:"nick_name_base64"`                                                          //用户昵称base64
	HeadImage          string    `gorm:"type:varchar(500)" json:"head_image"`                                                                //用户头像
	WxAlias            string    `gorm:"type:varchar(100)" json:"wx_alias"`                                                                  //微信号
	InvitedWxId        string    `gorm:"type:varchar(80)" json:"invited_wx_id"`                                                              //邀请人ID
	ChatNickName       string    `gorm:"type:varchar(100)" json:"chat_nick_name"`                                                            //群用户群内昵称
	ChatNickNameBase64 string    `gorm:"type:varchar(250)" json:"chat_nick_name_base64"`                                                     //用户群内昵称base64
	JoinDate           time.Time `gorm:"default:NULL" json:"join_date"`                                                                      //入群时间
	//todo 所有返回数据都没有时间和类型
	// JoinType
}

func (c *ChatRoomMember) Unmarshal(merchantNo, chatRoomId string, iter interface{}) error {
	return ChatRoomMemberUnmarshal(merchantNo, chatRoomId, iter, c)
}

func ChatRoomMemberUnmarshal(merchantNo, chatRoomId string, iter interface{}, member *ChatRoomMember) error {
	var input map[string]interface{}
	err := json.Unmarshal([]byte(goutils.ToString(iter)), &input)
	if err != nil {
		return err
	}
	m := goutils.NewMap(input)
	fansWxId, ok := m.GetString("vcFansWxId")
	if !ok {
		fansWxId, ok = m.GetString("vcWxId") // 入群和退群的回调为vcWxId 兼容
		if !ok || fansWxId == "" {
			return fmt.Errorf("vcFansWxId empty")
		}
	}
	member.MerchantNo = merchantNo
	member.ChatRoomId = chatRoomId
	member.FansWxId = fansWxId

	member.NickName, _ = m.GetString("vcNickName")
	member.NickNameBase64, ok = m.GetString("vcBase64NickName")
	if ok {
		// 如果有base64并且解码成功，则使用
		dec, err := base64.StdEncoding.DecodeString(member.NickNameBase64)
		if err != nil {
			member.NickName = goutils.ToString(dec)
		}
	}
	member.HeadImage, _ = m.GetString("vcHeadImages")
	member.WxAlias, _ = m.GetString("vcWxAlias")
	member.InvitedWxId, _ = m.GetString("vcInvitedWxId")

	member.ChatNickName, _ = m.GetString("vcChatNickName")
	member.ChatNickNameBase64, ok = m.GetString("vcChatBase64NickName")
	if ok {
		// 如果有base64并且解码成功，则使用
		dec, err := base64.StdEncoding.DecodeString(member.ChatNickNameBase64)
		if err != nil {
			member.ChatNickName = goutils.ToString(dec)
		}
	}

	//todo
	return nil
}

type ChatRoomMemberJoinCallbackStruct struct {
	MerchantNo string        `json:"vcMerChantNo"`
	ChatRoomId string        `json:"vcChatRoomId"`
	Data       []interface{} `json:"vcChatRoomUsers"`
}

func ChatRoomMemberJoinCallback(iter interface{}) (ret []*ChatRoomMember, err error) {
	ret = make([]*ChatRoomMember, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var callback ChatRoomMemberJoinCallbackStruct
		err = json.Unmarshal([]byte(goutils.ToString(data)), &callback)
		if err != nil {
			return
		}
		for _, iter := range callback.Data {
			obj := new(ChatRoomMember)
			err = ChatRoomMemberUnmarshal(callback.MerchantNo, callback.ChatRoomId, iter, obj)
			if err != nil {
				return
			}
			obj.JoinDate = time.Now() // 补上入群时间
			ret = append(ret, obj)
		}
	}
	return
}

type ChatRoomMemberQuitCallbackStruct struct {
	MerchantNo string        `json:"vcMerChantNo"`
	ChatRoomId string        `json:"vcChatRoomId"`
	Data       []interface{} `json:"vcChatRoomUsers"`
}

func ChatRoomMemberQuitCallback(iter interface{}) (ret []*ChatRoomMember, err error) {
	ret = make([]*ChatRoomMember, 0)
	rst := new(Callback)
	err = rst.Unmarshal(iter)
	if err != nil {
		return
	}
	for _, data := range rst.Each() {
		var callback ChatRoomMemberQuitCallbackStruct
		err = json.Unmarshal([]byte(goutils.ToString(data)), &callback)
		if err != nil {
			return
		}
		for _, iter := range callback.Data {
			obj := new(ChatRoomMember)
			err = ChatRoomMemberUnmarshal(callback.MerchantNo, callback.ChatRoomId, iter, obj)
			if err != nil {
				return
			}
			ret = append(ret, obj)
		}
	}
	return
}

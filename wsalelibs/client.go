package wsalelibs

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lvzhihao/goutils"
)

const (
	G_ROBOT  = "ApiRobotService"
	G_MSG    = "ApiMessageService"
	G_FRIEND = "ApiFriendsCircleService"
	G_CHAT   = "ApiChatRoomService"
)

var (
	WsaleApiPrefix                     string         = "http://app.wsale.net/Api/InterFace" //接口地址
	TimeZone                           string         = "Asia/Shanghai"                      //时区设置
	TimeLocation                       *time.Location                                        //当前时区
	DefaultTransportInsecureSkipVerify bool           = true
	DefaultTransportDisableCompression bool           = true
)

func init() {
	SetTimeZone(TimeZone)
	if os.Getenv("WSALE_TIMEZONE") != "" {
		loc, err := time.LoadLocation(os.Getenv("WSALE_TIMEZONE"))
		// loc set success
		if err == nil {
			TimeZone = os.Getenv("WSALE_TIMEZONE")
			TimeLocation = loc
		}
	}
}

// 设置时区
func SetTimeZone(zone string) error {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return err
	} else {
		TimeZone = zone
		TimeLocation = loc
		return nil
	}
}

// 客户端
type Client struct {
	Merchant   *Merchant    //商户信息
	Error      error        //最后一次错误
	Value      interface{}  //最后一次结果
	apiPrefix  string       //API
	httpClient *http.Client //
}

// 小U机器接口返回格式
type ApiResult struct {
	Result interface{} `json:"nResult"`   //1为成功, 其余为失败，获取vcMessage为Error信息
	Error  string      `json:"vcMessage"` //文字描述
	Data   interface{} `json:"Data"`      //返回结果
}

//  初始化一个新的实例
//  client := wsalelibs.NewClient()
//  client.DefaultTimeout = 30 * time.Second
func NewClient() *Client {
	client := &Client{
		apiPrefix: WsaleApiPrefix,
	}
	return client.EnsureHttpClient(10 * time.Second)
}

// 确认HttpClient
func (c *Client) EnsureHttpClient(timeout time.Duration) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	URL, err := url.Parse(c.apiPrefix)
	if err != nil && strings.ToLower(URL.Scheme) == "https" {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: DefaultTransportInsecureSkipVerify,
			},
			DisableCompression: DefaultTransportDisableCompression,
		}
	}
	c.httpClient = httpClient
	return c
}

func (c *Client) clone() *Client {
	return &Client{
		Merchant:   c.Merchant,
		Error:      c.Error,
		Value:      c.Value,
		apiPrefix:  c.apiPrefix,
		httpClient: c.httpClient,
	} //clone client
}

func (c *Client) M(m interface{}) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	switch m.(type) {
	case *Merchant:
		client.Merchant = m.(*Merchant)
	case string:
		client.Merchant, client.Error = LoadMerchant(m.(string))
	default:
		client.Merchant = nil
		client.AddError(fmt.Errorf("Merchant don't exists"))
	}
	return
}

func (c *Client) Robot(action string, ctx interface{}) (client *Client) {
	client = c.Action(G_ROBOT, action, ctx)
	return
}

func (c *Client) Msg(action string, ctx interface{}) (client *Client) {
	client = c.Action(G_MSG, action, ctx)
	return
}

func (c *Client) Friend(action string, ctx interface{}) (client *Client) {
	client = c.Action(G_FRIEND, action, ctx)
	return
}

func (c *Client) Chat(action string, ctx interface{}) (client *Client) {
	client = c.Action(G_CHAT, action, ctx)
	return
}

// 发送动作
// var ret []byte
// err := client.M(iter).Action("GetAccountList", nil).Result(&ret).Error
// OR:
// var ret []byte
// ctx := make(map[string]string, 0)
// ctx["vcRobotWxId"] = "xxxxx"
// err := client.M(iter).Action("SQRobotOffline", ctx).Result(&ret).Error
func (c *Client) Action(group, action string, ctx interface{}) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	request, err := client.request(group, action, ctx)
	if err != nil {
		client.AddError(err)
		client.Value = nil
	} else {
		value, err := client.scan(client.do(request))
		client.AddError(err)
		client.Value = value
	}
	return
}

// 当前商家信息
func (c *Client) CurrentM(cur *Merchant) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	*cur = *client.Merchant
	return
}

// Result 结果
func (c *Client) Result(iter interface{}) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	client.AddError(json.Unmarshal([]byte(goutils.ToString(client.Value)), &iter))
	return
}

func (c *Client) First(iter interface{}) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	var list []interface{}
	err := client.AddError(json.Unmarshal([]byte(goutils.ToString(client.Value)), &list))
	if err != nil {
		return
	}
	for _, obj := range list {
		b, _ := json.Marshal(obj)
		client.AddError(json.Unmarshal(b, &iter))
		break
	}
	return
}

func (c *Client) ResultKey(key string, iter interface{}) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	var res map[string]interface{}
	err := client.AddError(json.Unmarshal([]byte(goutils.ToString(client.Value)), &res))
	if err != nil {
		return
	}
	if rst, ok := res[key]; !ok {
		client.AddError(fmt.Errorf("no value"))
	} else {
		b, _ := json.Marshal(rst)
		client.AddError(json.Unmarshal(b, &iter))
	}
	return
}

func (c *Client) CheckSign(strCtx, strSign string) (client *Client) {
	client = c.clone()
	if client.Error != nil {
		return
	}
	if strings.Compare(c.sign(strCtx), strings.ToLower(strings.TrimSpace(strSign))) != 0 {
		client.AddError(fmt.Errorf("checksign error"))
	}
	return
}

// 设置最后一次错误
func (c *Client) AddError(err error) error {
	if err != nil {
		c.Error = err
	}
	return err
}

// create http request
func (c *Client) request(group, action string, ctx interface{}) (*http.Request, error) {
	var b []byte
	var err error
	if c.Merchant == nil {
		return nil, fmt.Errorf("merchant unknow")
	}
	switch ctx.(type) {
	case nil:
		input := make(map[string]interface{}, 0)
		input["vcMerChantNo"] = c.Merchant.MerchantNo
		b, err = json.Marshal(input)
	case map[string]string:
		input := ctx.(map[string]string)
		input["vcMerChantNo"] = c.Merchant.MerchantNo
		b, err = json.Marshal(input)
	case map[string]interface{}:
		input := ctx.(map[string]interface{})
		input["vcMerChantNo"] = c.Merchant.MerchantNo
		b, err = json.Marshal(input)
	default:
		return nil, fmt.Errorf("input params format error")
	}
	if err != nil {
		return nil, err
	}
	p := url.Values{}
	p.Set("strContext", string(b))
	p.Set("strSign", c.sign(string(b))) //todo
	req, err := http.NewRequest(
		"POST",
		strings.TrimRight(c.apiPrefix, "/")+"/"+group+".asmx/"+strings.TrimLeft(action, "/"),
		bytes.NewBufferString(p.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", goutils.ToString(len(p.Encode())))
	//req.Header.Add("X-Version", ")
	return req, nil
}

func (c *Client) fixJson(b []byte) []byte {
	b = bytes.Replace(b, []byte("\n"), []byte(""), -1)
	b = bytes.Replace(b, []byte("\r\n"), []byte(""), -1)
	return bytes.Replace(b, []byte("\\'"), []byte(""), -1)
}

// 审查请求结果
func (c *Client) scan(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//log.Printf("%s\n", b)
	var rst ApiResult
	err = json.Unmarshal(c.fixJson(b), &rst)
	if err != nil {
		return nil, err
	}
	if goutils.ToString(rst.Result) != "1" {
		if rst.Error != "" {
			return b, fmt.Errorf(rst.Error)
		} else {
			return b, fmt.Errorf("未知错误")
		}
	}
	if rst.Data != nil {
		return json.Marshal(rst.Data)
	} else {
		return b, nil
	}
}

// 发送请求
func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// 签名信息
func (c *Client) sign(strCtx string) string {
	if c.Merchant != nil {
		return fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimSpace(strCtx)+c.Merchant.MerchantSecret)))
	} else {
		return ""
	}
}

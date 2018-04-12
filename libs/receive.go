package wsalelibs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/lvzhihao/goutils"
	"go.uber.org/zap"
)

var (
	DefaultApplyChannelSize     int    = 1000
	DefaultWsaleCallbackSuccess string = "SUCCESS"
)

type Receive struct {
	app      *echo.Echo
	client   *Client
	Logger   *zap.Logger
	applyMap *sync.Map
}

type ReceiveApply struct {
	act     string
	handle  func(act, strContext string) error
	sync    bool
	channel chan string // 选择异步操作，提搞网络效率，handle需要确保数据可靠性
}

func NewReceive() *Receive {
	app := goutils.NewEcho()
	var logger *zap.Logger
	if os.Getenv("DEBUG") == "true" {
		logger, _ = zap.NewDevelopment()
		app.Logger.SetLevel(log.DEBUG)
	} else {
		logger, _ = zap.NewProduction()
	}
	return &Receive{
		app:      app,
		client:   NewClient(),
		Logger:   logger,
		applyMap: new(sync.Map),
	}
}

func (c *Receive) Check(ctx echo.Context) error {
	var params map[string]interface{}
	err := json.Unmarshal([]byte(ctx.FormValue("strContext")), &params)
	if err != nil {
		return err
	}
	if merchantNo, ok := params["vcMerChantNo"]; !ok {
		return fmt.Errorf("merchantNo empty")
	} else {
		return c.client.M(merchantNo).CheckSign(ctx.FormValue("strContext"), ctx.FormValue("strSign")).Error
	}
}

func (c *Receive) Sync(act string, handle func(act, strContext string) error) {
	apply := &ReceiveApply{
		act:     act,
		handle:  handle,
		sync:    true,
		channel: make(chan string, 0),
	}
	c.applyMap.Store(act, apply)

}

func (c *Receive) Async(act string, handle func(act, strContext string) error) {
	apply := &ReceiveApply{
		act:     act,
		handle:  handle,
		sync:    false,
		channel: make(chan string, DefaultApplyChannelSize),
	}
	c.applyMap.Store(act, apply)
	go apply.Consumer(c.Logger)
}

func (c *Receive) Start(host string) {
	defer c.Logger.Sync()
	c.app.Any("/*", func(ctx echo.Context) error {
		act := ctx.QueryParam("act")
		strContext := ctx.FormValue("strContext")
		strSign := ctx.FormValue("strSign")
		err := c.Check(ctx)
		if err != nil {
			c.Logger.Error("receive info check error", zap.String("act", act), zap.String("strContext", strContext), zap.String("strSign", strSign))
			return ctx.NoContent(http.StatusForbidden)
		} else {
			if apply, ok := c.applyMap.Load(act); !ok {
				c.Logger.Debug("action don't register", zap.String("act", act), zap.String("strContext", strContext), zap.String("strSign", strSign))
				return ctx.HTML(http.StatusOK, DefaultWsaleCallbackSuccess)
			} else {
				err := apply.(*ReceiveApply).Receive(strContext)
				if err != nil {
					c.Logger.Error("receive info failure", zap.Error(err), zap.String("act", act), zap.String("strContext", strContext), zap.String("strSign", strSign))
					return ctx.HTML(http.StatusInternalServerError, err.Error())
				} else {
					c.Logger.Debug("receive info success", zap.String("act", act), zap.String("strContext", strContext), zap.String("strSign", strSign))
					return ctx.HTML(http.StatusOK, DefaultWsaleCallbackSuccess)
				}
			}
		}
	})
	goutils.EchoStartWithGracefulShutdown(c.app, host)
}

func (c *ReceiveApply) Receive(strContext string) error {
	if c.sync == true {
		return c.handle(c.act, strContext)
	} else {
		c.channel <- strContext
		return nil
	}
}

func (c *ReceiveApply) Consumer(logger *zap.Logger) {
	for {
		select {
		case strContext := <-c.channel:
			err := c.handle(c.act, strContext)
			if err != nil {
				logger.Error("sync receive info failure", zap.Error(err), zap.String("act", c.act), zap.String("strContext", strContext))
			}
		}
	}
}

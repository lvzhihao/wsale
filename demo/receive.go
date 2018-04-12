package main

import (
	"flag"

	"github.com/lvzhihao/wsale/wsalelibs"
	"go.uber.org/zap"
)

var (
	demoReceiveHost    string = ":8279"
	demoMerchantNo     string
	demoMerchantSecret string
)

func init() {
	flag.StringVar(&demoMerchantNo, "no", "", "merchantNo")
	flag.StringVar(&demoMerchantSecret, "secret", "", "merchantSecret")
	flag.Parse()
	wsalelibs.EnableMerchant(&wsalelibs.Merchant{
		MerchantNo:     demoMerchantNo,
		MerchantSecret: demoMerchantSecret,
	})
}

func main() {
	logger, _ := zap.NewProduction()
	receive := wsalelibs.NewReceive()
	receive.Sync("robotinfo", func(act, ctx string) error {
		logger.Info(act, zap.Any("ctx", ctx))
		return nil
	})
	receive.Async("msg", func(act, ctx string) error {
		logger.Info(act, zap.Any("ctx", ctx))
		return nil
	})
	receive.Sync("friend", func(act, ctx string) error {
		logger.Info(act, zap.Any("ctx", ctx))
		return nil
	})
	receive.Sync("msgresult", func(act, ctx string) error {
		logger.Info(act, zap.Any("ctx", ctx))
		return nil
	})
	receive.Start(demoReceiveHost)
}

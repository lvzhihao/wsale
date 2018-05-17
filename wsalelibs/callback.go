package wsalelibs

import (
	"encoding/json"
	"fmt"

	"github.com/lvzhihao/goutils"
)

type Callback struct {
	MerchantNo string      `json:"vcMerChantNo"`
	Data       interface{} `json:"Data"`
}

func (c *Callback) Unmarshal(iter interface{}) error {
	err := json.Unmarshal([]byte(goutils.ToString(iter)), c)
	if err != nil {
		return err
	}
	return c.Check()
}

func (c *Callback) Check() error {
	if c.MerchantNo == "" {
		return fmt.Errorf("merchant empty")
	}
	return nil
}

func (c *Callback) Each() []interface{} {
	rst := make([]interface{}, 0)
	json.Unmarshal([]byte(goutils.ToString(c.Data)), &rst)
	return rst
}

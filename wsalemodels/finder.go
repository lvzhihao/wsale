package wsalemodels

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/lvzhihao/goutils"
)

type Finder struct {
	Columns  []string
	Where    map[string]interface{}
	Order    []string
	Offset   int32
	Limit    int32
	Page     int32
	PageSize int32
}

func (c *Finder) Use(db *gorm.DB) *gorm.DB {
	ret := db
	ret = c.columns(ret)
	ret = c.where(ret)
	ret = c.order(ret)
	ret = c.offset(ret)
	ret = c.limit(ret)
	return ret
}

func FinderUnmarshal(iter interface{}, finder *Finder) error {
	return json.Unmarshal([]byte(goutils.ToString(iter)), &finder)
}

func NewFinder() *Finder {
	return &Finder{
		Columns:  make([]string, 0),
		Where:    make(map[string]interface{}, 0),
		Order:    make([]string, 0),
		Offset:   -1,
		Limit:    -1,
		Page:     -1,
		PageSize: 0,
	}
}

func (c *Finder) Unmarshal(iter interface{}) error {
	return FinderUnmarshal(iter, c)
}

func (c *Finder) AddColumn(column ...string) *Finder {
	c.Columns = append(c.Columns, column...)
	return c
}

func (c *Finder) AddWhere(key string, value interface{}) *Finder {
	c.Where[key] = value
	return c
}

func (c *Finder) AddOrder(order ...string) *Finder {
	c.Order = append(c.Order, order...)
	return c
}

func (c *Finder) SetPage(page int32) *Finder {
	c.Page = page
	return c
}

func (c *Finder) SetPageSize(pageSize int32) *Finder {
	c.PageSize = pageSize
	return c
}

func (c *Finder) columns(db *gorm.DB) *gorm.DB {
	if len(c.Columns) > 0 {
		db = db.Select(c.Columns)
	}
	return db
}

func (c *Finder) where(db *gorm.DB) *gorm.DB {
	for k, v := range c.Where {
		db = db.Where(k, v)
	}
	return db
}

func (c *Finder) order(db *gorm.DB) *gorm.DB {
	for _, v := range c.Order {
		db = db.Order(v)
	}
	return db
}

func (c *Finder) offset(db *gorm.DB) *gorm.DB {
	if c.Page > 0 && c.PageSize > 0 {
		c.Offset = (c.Page - 1) * c.PageSize
	} else if c.Offset == 0 {
		c.Offset = -1
	}
	return db.Offset(c.Offset)
}

func (c *Finder) limit(db *gorm.DB) *gorm.DB {
	if c.PageSize > 0 {
		c.Limit = c.PageSize
	} else if c.Limit == 0 {
		c.Limit = -1
	}
	return db.Limit(c.Limit)
}

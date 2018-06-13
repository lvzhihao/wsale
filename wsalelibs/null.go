package wsalelibs

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/go-sql-driver/mysql"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (c *NullTime) Scan(value interface{}) (err error) {
	var t mysql.NullTime
	if err := t.Scan(value); err != nil {
		c.Valid = false
		return err
	}
	c.Time, c.Valid = t.Time, true
	return nil
}

func (c NullTime) Value() (driver.Value, error) {
	return mysql.NullTime{
		c.Time,
		c.Valid,
	}.Value()
}

func (c *NullTime) MarshalJSON() ([]byte, error) {
	if c.Valid {
		return json.Marshal(c.Time)
	} else {
		return []byte("null"), nil
	}
}

func (c *NullTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	err := json.Unmarshal(data, &t)
	if err != nil {
		c.Valid = false
	} else {
		c.Set(t)
	}
	return err
}

func (c *NullTime) Set(t time.Time) {
	c.Time = t
	/*
		if c.Time.IsZero() {
			c.Valid = false
		} else {
			c.Valid = true
		}
	*/
	c.Valid = true
}

func (c *NullTime) Get() (time.Time, bool) {
	return c.Time, c.Valid
}

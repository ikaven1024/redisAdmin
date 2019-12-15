package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createTime,omitempty"`
	UpdatedAt time.Time  `json:"updateTime,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleteTime,omitempty"`
}

type StringSlice struct {
	Data []string
}

func NewStringSlice(data []string) StringSlice {
	return StringSlice{Data: data}
}

var _ driver.Valuer = (*StringSlice)(nil)
var _ sql.Scanner = (*StringSlice)(nil)

func (c StringSlice) Value() (driver.Value, error) {
	b, err := json.Marshal(c.Data)
	return string(b), err
}

func (c *StringSlice) Scan(input interface{}) error {
	s := input.(string)
	return json.Unmarshal([]byte(s), &c.Data)
}

package models

import (
	"github.com/dailydotdev/platform-go-common/util/uuid"
)

type Table2 struct {
	ID   uuid.UUID `gorm:"column:id;type:uuid" json:"id"`
	Data string    `gorm:"column:data;type:text" json:"data"`
}

func (c Table2) TableName() string {
	return "table2"
}
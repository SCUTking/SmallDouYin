package model

import "time"

type Message struct {
	MessageId  uint64    `gorm:"column:message_id;primary_key;NOT NULL"`
	ToUserId   uint64    `gorm:"column:to_user_id;NOT NULL"`
	FromUserId uint64    `gorm:"column:from_user_id;NOT NULL"`
	Content    string    `gorm:"column:content;NOT NULL"`
	CreatedAt  time.Time `gorm:"column:created_at" redis:"-"`
	UpdatedAt  time.Time `gorm:"column:updated_at" redis:"-"`
}

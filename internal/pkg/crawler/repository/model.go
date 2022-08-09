package repository

import "time"

type Event struct {
	Crawler   string    `gorm:"type:varchar(128)"`
	Job       string    `gorm:"type:varchar(128)"`
	UserName  string    `gorm:"type:varchar(128)"`
	UID       string    `gorm:"type:varchar(128);uniqueIndex:uid_name"`
	Name      string    `gorm:"type:varchar(128);uniqueIndex:uid_name"`
	Message   string    `gorm:"type:varchar(65535)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Event) TableName() string {
	return "event"
}

type Channel struct {
	ID   string `gorm:"type:varchar(128);column:id;primary_key"`
	Name string `gorm:"type:varchar(128);index"`
}

func (Channel) TableName() string {
	return "channel"
}

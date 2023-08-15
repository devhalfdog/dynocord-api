package database

import "gorm.io/gorm"

type Model struct {
	Connector *gorm.DB
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
}

type Chat struct {
	No        uint64 `json:"no" gorm:"column:no"`
	From      string `json:"from" gorm:"column:fromStreamer"`
	To        string `json:"to" gorm:"column:toStreamer"`
	Message   string `json:"message" gorm:"column:message"`
	Emoji     string `json:"emoji" gorm:"column:emoji"`
	Timestamp uint64 `json:"timestamp" gorm:"column:timestamp"`
}

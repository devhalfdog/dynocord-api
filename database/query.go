package database

import (
	e "github.com/devhalfdog/dynocord-api/errors"
)

func (m *Model) CreateChat(c *Chat) error {
	err := m.Connector.Table("chats").Create(&c).Error
	if err != nil {
		return e.ErrCreateChat
	}

	return nil
}

func (m *Model) GetChat(streamer string, before int64) ([]Chat, error) {
	c := new([]Chat)

	query := m.Connector.
		Table("chats").
		Where("fromStreamer = ? AND toStreamer = ?", streamer, streamer)

	if before > 0 {
		query = query.Where("no > ?", before)
	}

	err := query.Limit(10).Find(&c).Error

	if err != nil {
		return *c, e.ErrGetChat
	}

	return *c, nil
}

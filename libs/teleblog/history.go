package teleblog

import (
	"fmt"
	"strconv"

	"gopkg.in/telebot.v3"
)

type HistoryMessageTextEntity struct {
	Type telebot.EntityType `json:"type"`
	Text string             `json:"text"`
}

type HistoryMessage struct {
	Id           int                        `json:"id"`
	Type         string                     `json:"type"` // service | message
	Date         string                     `json:"date"`
	DateUnix     string                     `json:"date_unixtime"`
	Edited       string                     `json:"edited"`
	EditedUnix   string                     `json:"edited_unixtime"`
	From         string                     `json:"from"`
	FromId       string                     `json:"from_id"`
	TextEntities []HistoryMessageTextEntity `json:"text_entities"`
	File         *string                    `json:"file"`
}

type History struct {
	Id       int64            `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"` // "public_channel"
	Messages []HistoryMessage `json:"messages"`
}

func (h *History) AddMessageGetChatTgId() (int64, error) {
	return strconv.ParseInt(fmt.Sprintf("-100%d", h.Id), 10, 64)
}

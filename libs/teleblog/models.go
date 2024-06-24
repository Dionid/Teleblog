package teleblog

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"gopkg.in/telebot.v3"
)

// # User

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel

	Email    string `json:"email" db:"email"`
	Verified bool   `json:"verified" db:"verified"`
	Name     string `json:"name" db:"name"`
	// PasswordHash string `json:"passwordHash" db:"password_hash"`

	TgUserId   int64  `json:"tgUserId" db:"tg_user_id"`
	TgUsername string `json:"tgUsername" db:"tg_username"`
}

func (m *User) TableName() string {
	return "users"
}

func UserQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&User{})
}

// # Verification Token

var _ models.Model = (*TgVerificationToken)(nil)

type TgVerificationToken struct {
	models.BaseModel

	UserId   string `json:"userId" db:"user_id"`
	Value    string `json:"value" db:"value"`
	Verified bool   `json:"verified" db:"verified"`
}

func (m *TgVerificationToken) TableName() string {
	return "tg_verification_token"
}

func TgVerificationTokenQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&TgVerificationToken{})
}

// # Chat

var _ models.Model = (*Chat)(nil)

type Chat struct {
	models.BaseModel

	UserId       string `json:"userId" db:"user_id"`
	LinkedChatId string `json:"linkedChatId" db:"linked_chat_id"`

	TgUsername     string `json:"tgUsername" db:"tg_username"`
	TgChatId       int64  `json:"tgChatId" db:"tg_chat_id"`
	TgType         string `json:"tgType" db:"tg_type"` // channel | group
	TgLinkedChatId int64  `json:"tgLinkedChatId" db:"tg_linked_chat_id"`
}

func (m *Chat) TableName() string {
	return "chat"
}

func ChatQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Chat{})
}

// # Post

var _ models.Model = (*Post)(nil)

type Post struct {
	models.BaseModel

	ChatId      string `json:"chatId" db:"chat_id"`
	IsTgMessage bool   `json:"isTgMessage" db:"is_tg_message"`

	Text string `json:"text" db:"text"`

	TgPostId   int                                    `json:"thPostId" db:"tg_post_id"`
	TgEntities types.JsonArray[telebot.MessageEntity] `json:"tgEntities" db:"tg_entities"`
}

func (m *Post) TableName() string {
	return "post"
}

func PostQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Post{})
}

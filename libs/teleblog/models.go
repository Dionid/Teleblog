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
	TgType         string `json:"tgType" db:"tg_type"` //  "private" | "group" | "supergroup" | "channel" | "privatechannel"
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

	ChatId             string `json:"chatId" db:"chat_id"`
	IsTgMessage        bool   `json:"isTgMessage" db:"is_tg_message"`
	IsTgHistoryMessage bool   `json:"isTgHistoryMessage" db:"is_tg_history_message"`

	Text string `json:"text" db:"text"`

	TgMessageId       int                                       `json:"tgMessageId" db:"tg_post_id"`
	TgGroupMessageId  int                                       `json:"tgGroupMessageId" db:"tg_group_message_id"`
	TgEntities        types.JsonArray[telebot.MessageEntity]    `json:"tgEntities" db:"tg_entities"`
	TgHistoryEntities types.JsonArray[HistoryMessageTextEntity] `json:"tgHistoryEntities" db:"tg_history_entities"`
	TgMessageRaw      types.JsonMap                             `json:"tgMessageRaw" db:"tg_message_raw"`
}

func (m *Post) TableName() string {
	return "post"
}

func PostQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Post{})
}

// # Comment

var _ models.Model = (*Comment)(nil)

type Comment struct {
	models.BaseModel

	ChatId string `json:"chatId" db:"chat_id"`
	PostId string `json:"postId" db:"post_id"`

	Text string `json:"text" db:"text"`

	TgMessageId        int                                    `json:"tgMessageId" db:"tg_comment_id"`
	TgEntities         types.JsonArray[telebot.MessageEntity] `json:"tgEntities" db:"tg_entities"`
	TgMessageRaw       types.JsonMap                          `json:"tgMessageRaw" db:"tg_message_raw"`
	TgReplyToMessageId int                                    `json:"tgReplyToMessageId" db:"tg_reply_to_message_id"`
}

func (m *Comment) TableName() string {
	return "Comment"
}

func CommentQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Comment{})
}

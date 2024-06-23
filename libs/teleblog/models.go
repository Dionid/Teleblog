package teleblog

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// # User

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel

	Email        string `json:"email" db:"email"`
	Verified     bool   `json:"verified" db:"verified"`
	Name         string `json:"name" db:"name"`
	PasswordHash string `json:"passwordHash" db:"password_hash"`
	Avatar       string `json:"avatar" db:"avatar"`

	TelegramUserId int64 `json:"telegramUserId" db:"telegram_user_id"`
}

func (m *User) TableName() string {
	return "users"
}

func UserQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&User{})
}

// # Verification Token

var _ models.Model = (*VerificationToken)(nil)

type VerificationToken struct {
	models.BaseModel

	UserId   string `json:"userId" db:"user_id"`
	Value    string `json:"value" db:"value"`
	Verified bool   `json:"verified" db:"verified"`
}

func (m *VerificationToken) TableName() string {
	return "verification_token"
}

func VerificationTokenQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&VerificationToken{})
}

// # Source

var _ models.Model = (*Source)(nil)

type Source struct {
	models.BaseModel

	UserId         string `json:"userId" db:"user_id"`
	LinkedSourceId string `json:"linkedSourceId" db:"linked_source_id"`

	Username     string `json:"username" db:"username"`
	ChatId       int64  `json:"chatId" db:"chat_id"`
	Type         string `json:"type" db:"type"` // channel | group
	LinkedChatId int64  `json:"linkedChatId" db:"linked_chat_id"`
}

func (m *Source) TableName() string {
	return "source"
}

func SourceQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Source{})
}

// # Post

var _ models.Model = (*Post)(nil)

type Post struct {
	models.BaseModel

	SourceId          string `json:"sourceId" db:"source_id"`
	TelegramPostId    int64  `json:"postId" db:"post_id"`
	IsTelegramMessage bool   `json:"isTelegramMessage" db:"is_telegram_message"`
	Message           string `json:"message" db:"message"`
}

func (m *Post) TableName() string {
	return "post"
}

func PostQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Post{})
}

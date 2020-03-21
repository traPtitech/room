package repository

import (
	"time"

	"github.com/gofrs/uuid"
)

// Model is defalut
type Model struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updateAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// StartEndTime has start and end time
type StartEndTime struct {
	TimeStart time.Time `json:"timeStart" gorm:"type:TIME;"`
	TimeEnd   time.Time `json:"timeEnd" gorm:"type:TIME;"`
}

// User traQユーザー情報構造体
type User struct {
	// ID traQID
	ID uuid.UUID `gorm:"type:char(36); primary_key"`
	// Admin アプリの管理者かどうか
	Admin       bool   `gorm:"not null"`
	Name        string `gorm:"-"`
	DisplayName string `gorm:"-"`
}

// UserSession has user session
type UserSession struct {
	Token         string    `gorm:"primary_key; type:char(32);"`
	UserID        uuid.UUID `gorm:"type:char(36);"`
	Authorization string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

// Tag Room Group Event have tags
type Tag struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name     string    `json:"name" gorm:"unique; type:varchar(16)"`
	Official bool      `json:"official"`
	Locked   bool      `json:"locked" gorm:"-"`
	Model
}

// EventTag is many to many table
type EventTag struct {
	TagID   uuid.UUID `gorm:"type:char(36); primary_key"`
	EventID uuid.UUID `gorm:"type:char(36); primary_key"`
	Locked  bool
}

// GroupUsers is many to many table
type GroupUsers struct {
	GroupID uuid.UUID `gorm:"type:char(36); primary_key"`
	UserID  uuid.UUID `gorm:"type:char(36); primary_key"`
}

// Room 部屋情報
type Room struct {
	ID            uuid.UUID      `json:"id" gorm:"type:char(36);primary_key"`
	Place         string         `json:"place" gorm:"type:varchar(16);unique_index:idx_room_unique"`
	TimeStart     time.Time      `json:"timeStart" gorm:"type:TIME; unique_index:idx_room_unique"`
	TimeEnd       time.Time      `json:"timeEnd" gorm:"type:TIME; unique_index:idx_room_unique"`
	AvailableTime []StartEndTime `json:"availableTime" gorm:"-"`
	Model
}

// Group グループ情報
// Group is not user JSON
type Group struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"type:varchar(32);not null"`
	Description string    `gorm:"type:varchar(1024)"`
	JoinFreely  bool
	Members     []User    `gorm:"many2many:group_users; association_autoupdate:false;association_autocreate:false"`
	CreatedBy   uuid.UUID `gorm:"type:char(36);"`
	Model
}

// Event 予約情報
type Event struct {
	ID            uuid.UUID `json:"eventId" gorm:"type:char(36);primary_key"`
	Name          string    `json:"name" gorm:"type:varchar(32); not null"`
	Description   string    `json:"description" gorm:"type:varchar(1024)"`
	GroupID       uuid.UUID `json:"groupId" gorm:"type:char(36);not null"`
	Group         Group     `json:"-" gorm:"foreignkey:group_id; save_associations:false"`
	RoomID        uuid.UUID `json:"roomId" gorm:"type:char(36);not null"`
	Room          Room      `json:"-" gorm:"foreignkey:room_id; save_associations:false"`
	TimeStart     string    `json:"timeStart" gorm:"type:TIME"`
	TimeEnd       string    `json:"timeEnd" gorm:"type:TIME"`
	CreatedBy     uuid.UUID `json:"createdBy" gorm:"type:char(36);"`
	AllowTogether bool      `json:"sharedRoom"`
	Tags          []Tag     `json:"tags" gorm:"many2many:event_tags; association_autoupdate:false;association_autocreate:false"`
	Model
}

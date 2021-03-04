// Code generated by gotypeconverter; DO NOT EDIT.
package infra

import (
	"time"

	"github.com/traPtitech/knoQ/domain"
	"gorm.io/gorm"
)

func ConvertEventTodomainEvent(src Event) (dst domain.Event) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Room = ConvertRoomTodomainRoom(src.Room)
	dst.Group = ConvertGroupTodomainGroup(src.Group)
	dst.TimeStart = src.TimeStart
	dst.TimeEnd = src.TimeEnd
	dst.CreatedBy = ConvertUserMetaTodomainUser(src.CreatedBy)
	dst.Tags = make([]domain.EventTag, len(src.Tags))
	for i := range src.Tags {
		dst.Tags[i] = ConvertTagTodomainEventTag(src.Tags[i])
	}
	dst.AllowTogether = src.AllowTogether
	dst.Model.CreatedAt = src.Model.CreatedAt
	dst.Model.UpdatedAt = src.Model.UpdatedAt
	(*dst.Model.DeletedAt) = ConvertgormDeletedAtTotimeTime(src.Model.DeletedAt)
	return
}

func ConvertGroupTodomainGroup(src Group) (dst domain.Group) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.JoinFreely = src.JoinFreely
	dst.Members = make([]domain.User, len(src.Members))
	for i := range src.Members {
		dst.Members[i] = ConvertUserMetaTodomainUser(src.Members[i])
	}
	dst.CreatedBy = ConvertUserMetaTodomainUser(src.CreatedBy)
	return
}
func ConvertRoomTodomainRoom(src Room) (dst domain.Room) {
	dst.ID = src.ID
	dst.Place = src.Place
	dst.Verified = src.Verified
	dst.TimeStart = src.TimeStart
	dst.TimeEnd = src.TimeEnd
	dst.Events = make([]domain.Event, len(src.Events))
	for i := range src.Events {
		dst.Events[i] = ConvertEventTodomainEvent(src.Events[i])
	}
	dst.CreatedBy = ConvertUserMetaTodomainUser(src.CreatedBy)
	return
}

func ConvertTagTodomainEventTag(src Tag) (dst domain.EventTag) {
	dst.Tag.ID = src.ID
	dst.Tag.Name = src.Name
	dst.Locked = src.Locked
	return
}
func ConvertUserMetaTodomainUser(src UserMeta) (dst domain.User) {
	dst.ID = src.ID
	return
}

func ConvertdomainWriteEventParamsToEvent(src domain.WriteEventParams) (dst Event) {
	dst.Name = src.Name
	dst.Description = src.Description
	dst.GroupID = src.GroupID
	dst.RoomID = src.RoomID
	dst.TimeStart = src.TimeStart
	dst.TimeEnd = src.TimeEnd
	dst.AllowTogether = src.AllowTogether
	dst.Tags = make([]Tag, len(src.Tags))
	for i := range src.Tags {
		dst.Tags[i].Name = src.Tags[i].Name
		dst.Tags[i].Locked = src.Tags[i].Locked
		dst.Tags[i].Model.DeletedAt.Valid = src.Tags[i].Locked // bug
	}
	dst.Model.CreatedAt = src.TimeStart // bug
	return
}
func ConvertgormDeletedAtTotimeTime(src gorm.DeletedAt) (dst time.Time) {
	dst = src.Time
	return
}

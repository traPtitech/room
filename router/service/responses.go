package service

import (
	"time"

	repo "github.com/traPtitech/knoQ/repository"

	"github.com/gofrs/uuid"
)

type GroupRes struct {
	ID uuid.UUID `json:"groupId"`
	GroupReq
	IsTraQGroup bool      `json:"isTraQGroup"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// EventRes is event response
type EventRes struct {
	ID uuid.UUID `json:"eventId"`
	EventReq
	Tags      []TagRelationRes `json:"tags"`
	CreatedBy uuid.UUID        `json:"createdBy"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	DeletedAt *time.Time       `json:"deletedAt,omitempty"`
}

// TagRelationRes show relation one to tag
type TagRelationRes struct {
	ID uuid.UUID `json:"tagId"`
	TagRelationReq
}

type UserRes struct {
	ID          uuid.UUID `json:"userId"`
	Admin       bool      `json:"admin"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
}

type RoomRes struct {
	ID        uuid.UUID `json:"roomId"`
	Place     string    `json:"place"`
	Public    bool      `json:"public"`
	TimeStart string    `json:"timeStart"`
	TimeEnd   string    `json:"timeEnd"`
	CreatedBy uuid.UUID `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TagRes struct {
	ID        uuid.UUID `json:"tagId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FormatGroupRes(g *repo.Group, IsTraQgroup bool) *GroupRes {
	res := &GroupRes{
		ID: g.ID,
		GroupReq: GroupReq{
			Name:        g.Name,
			Description: g.Description,
			JoinFreely:  g.JoinFreely,
			Members:     formatGroupMembersRes(g.Members),
			Admins:      formatGroupAdminsRes(g.Admins),
		},
		IsTraQGroup: IsTraQgroup,
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	return res
}

func formatGroupMembersRes(ms []repo.GroupUsers) []uuid.UUID {
	ids := make([]uuid.UUID, len(ms))
	for i, m := range ms {
		ids[i] = m.UserID
	}
	return ids
}

func formatGroupAdminsRes(ms []repo.GroupAdmins) []uuid.UUID {
	ids := make([]uuid.UUID, len(ms))
	for i, m := range ms {
		ids[i] = m.UserID
	}
	return ids
}

func FormatGroupsRes(gs []*repo.Group, IsTraQGroup bool) []*GroupRes {
	res := make([]*GroupRes, len(gs))
	for i, g := range gs {
		res[i] = FormatGroupRes(g, IsTraQGroup)
	}
	return res
}

func FormatTagsRes(ts []repo.Tag) []TagRelationRes {
	res := make([]TagRelationRes, len(ts))
	for i, t := range ts {
		res[i] = TagRelationRes{
			ID: t.ID,
			TagRelationReq: TagRelationReq{
				Name:   t.Name,
				Locked: t.Locked,
			},
		}
	}
	return res

}

func FormatEventRes(e *repo.Event) *EventRes {
	return &EventRes{
		ID: e.ID,
		EventReq: EventReq{
			Name:          e.Name,
			Description:   e.Description,
			AllowTogether: e.AllowTogether,
			TimeStart:     e.TimeStart,
			TimeEnd:       e.TimeEnd,
			RoomID:        e.RoomID,
			GroupID:       e.GroupID,
			Place:         e.Room.Place,
			Admins:        FormatEventAdmins(e.Admins),
		},
		Tags:      FormatTagsRes(e.Tags),
		CreatedBy: e.CreatedBy,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

func FormatEventAdmins(ea []repo.EventAdmins) []uuid.UUID {
	ids := make([]uuid.UUID, len(ea))
	for i, m := range ea {
		ids[i] = m.UserID
	}
	return ids

}

func FormatEventsRes(es []*repo.Event) []*EventRes {
	res := make([]*EventRes, len(es))
	for i, e := range es {
		res[i] = FormatEventRes(e)
	}
	return res
}

func FormatUserRes(u *User) *UserRes {
	return &UserRes{
		ID:          u.ID,
		Admin:       u.Admin,
		Name:        u.Name,
		DisplayName: u.DisplayName,
	}
}

func FormatUsersRes(us []*User) []*UserRes {
	res := make([]*UserRes, len(us))
	for i, u := range us {
		res[i] = FormatUserRes(u)
	}
	return res
}

func FormatRoomRes(r *repo.Room) *RoomRes {
	return &RoomRes{
		ID:        r.ID,
		Place:     r.Place,
		Public:    r.Public,
		TimeStart: r.TimeStart.Format(time.RFC3339),
		TimeEnd:   r.TimeEnd.Format(time.RFC3339),
		CreatedBy: r.CreatedBy,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func FormatTagRes(t *repo.Tag) *TagRes {
	return &TagRes{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

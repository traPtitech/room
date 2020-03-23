package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	traQutils "github.com/traPtitech/traQ/utils"
)

func TestGormRepository_CreateRoom(t *testing.T) {
	t.Parallel()
	repo, _, _ := setupGormRepo(t, common)

	params := WriteRoomParams{
		Place:     traQutils.RandAlphabetAndNumberString(10),
		TimeStart: time.Now(),
		TimeEnd:   time.Now().Add(1 * time.Hour),
	}
	if room, err := repo.CreateRoom(params); assert.NoError(t, err) {
		assert.NotNil(t, room)
	}

	t.Run("Time error", func(t *testing.T) {
		params.TimeEnd = params.TimeStart.Add(-1 * time.Hour)
		_, err := repo.CreateRoom(params)
		assert.EqualError(t, err, ErrInvalidArg.Error())
	})
}

func TestGormRepository_UpdateRoom(t *testing.T) {
	repo, _, _ := setupGormRepo(t, common)
	room := mustMakeRoom(t, repo, traQutils.RandAlphabetAndNumberString(10))

	params := WriteRoomParams{
		Place:     traQutils.RandAlphabetAndNumberString(10),
		TimeStart: time.Now(),
		TimeEnd:   time.Now().Add(3 * time.Hour),
	}
	if room, err := repo.UpdateRoom(room.ID, params); assert.NoError(t, err) {
		assert.NotNil(t, room)
	}

	t.Run("Time error", func(t *testing.T) {
		params.TimeEnd = params.TimeStart.Add(-1 * time.Hour)
		_, err := repo.UpdateRoom(room.ID, params)
		assert.EqualError(t, err, ErrInvalidArg.Error())
	})

}

func TestGormRepository_DeleteRoom(t *testing.T) {
	t.Parallel()
	repo, _, _ := setupGormRepo(t, common)
	room := mustMakeRoom(t, repo, traQutils.RandAlphabetAndNumberString(10))

	t.Run("Delete existing room", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteRoom(room.ID)
		assert.NoError(t, err)
	})

	t.Run("Delete not existing room", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteRoom(mustNewUUIDV4(t))
		assert.EqualError(t, err, ErrNotFound.Error())
	})
}

func TestGormRepository_GetRoom(t *testing.T) {
	t.Parallel()
	repo, _, _ := setupGormRepo(t, common)
	room := mustMakeRoom(t, repo, traQutils.RandAlphabetAndNumberString(10))

	if room, err := repo.GetRoom(room.ID); assert.NoError(t, err) {
		assert.NotNil(t, room)
	}
}

func TestRoom_calcAvailableTime(t *testing.T) {
	now := time.Now()
	type fields struct {
		TimeStart time.Time
		TimeEnd   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   []Event
		want   []StartEndTime
	}{
		{
			name: "success",
			fields: fields{
				TimeStart: now,
				TimeEnd:   now.Add(10 * time.Hour),
			},
			args: []Event{
				{
					TimeStart:     now.Add(1 * time.Hour),
					TimeEnd:       now.Add(2 * time.Hour),
					AllowTogether: false,
				},
			},
			want: []StartEndTime{
				{
					TimeStart: now,
					TimeEnd:   now.Add(1 * time.Hour),
				},
				{
					TimeStart: now.Add(2 * time.Hour),
					TimeEnd:   now.Add(10 * time.Hour),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Room{
				TimeStart: tt.fields.TimeStart,
				TimeEnd:   tt.fields.TimeEnd,
			}
			r.calcAvailableTime(tt.args)
			assert.Equal(t, tt.want, r.AvailableTime)
		})
	}
}

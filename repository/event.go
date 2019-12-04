package repository

import (
	"errors"
	"net/url"
	"room/utils"
	"strconv"
)

func FindRvs(values url.Values) ([]Event, error) {
	events := []Event{}
	cmd := DB.Preload("Group").Preload("Group.Members").Preload("Group.CreatedBy").Preload("Room").Preload("Tags")

	if values.Get("id") != "" {
		id, _ := strconv.Atoi(values.Get("id"))
		cmd = cmd.Where("id = ?", id)
	}

	if values.Get("name") != "" {
		cmd = cmd.Where("name LIKE ?", "%"+values.Get("name")+"%")
	}

	if values.Get("traQID") != "" {
		groupsID, err := GetGroupIDsBytraQID(values.Get("traQID"))
		if err != nil {
			return nil, err
		}
		cmd = cmd.Where("group_id in (?)", groupsID)
	}

	if values.Get("groupid") != "" {
		groupid, _ := strconv.Atoi(values.Get("groupid"))
		cmd = cmd.Where("group_id = ?", groupid)
	}

	if values.Get("roomid") != "" {
		roomid, _ := strconv.Atoi(values.Get("roomid"))
		cmd = cmd.Where("room_id = ?", roomid)
	}

	if values.Get("date_begin") != "" {
		cmd = cmd.Where("rooms.date >= ?", values.Get("date_begin"))
	}
	if values.Get("date_end") != "" {
		cmd = cmd.Where("rooms.date <= ?", values.Get("date_end"))
	}

	// room の日付を見たい
	if err := cmd.Select("events.*").Joins("JOIN rooms on rooms.id = room_id").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Event) AfterFind() (err error) {
	e.GroupID = 0
	e.RoomID = 0
	return
}

// 時間が部屋の範囲内か、endがstartの後かどうか確認する
func (rv *Event) TimeConsistency() error {
	timeStart, err := utils.StrToTime(rv.TimeStart)
	if err != nil {
		return err
	}
	timeEnd, err := utils.StrToTime(rv.TimeEnd)
	if err != nil {
		return err
	}
	if !(rv.Room.InTime(timeStart) && rv.Room.InTime(timeEnd)) {
		return errors.New("invalid time")
	}
	if !timeStart.Before(timeEnd) {
		return errors.New("invalid time")
	}
	return nil
}

// GetCreatedBy get who created it
func (rv *Event) GetCreatedBy() (string, error) {
	if err := DB.First(&rv).Error; err != nil {
		return "", err
	}
	return rv.CreatedBy, nil
}

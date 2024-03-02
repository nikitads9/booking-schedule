package convert

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/model"
	"time"

	"github.com/gofrs/uuid"
)

func ToEventInfo(req *api.Event) (*model.EventInfo, error) {
	if req == nil {
		return nil, api.ErrEmptyRequest
	}

	res := &model.EventInfo{
		UserID:    req.UserID,
		SuiteID:   req.SuiteID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	if req.EventID != uuid.Nil {
		res.ID = req.EventID
	}

	if req.NotifyAt.Valid {
		dur, err := time.ParseDuration(req.NotifyAt.String)
		if err != nil {
			return nil, err
		}
		res.NotifyAt = dur
	}

	return res, nil
}

func ToApiEventInfo(mod *model.EventInfo) *api.EventInfo {

	res := &api.EventInfo{
		ID:        mod.ID,
		SuiteID:   mod.SuiteID,
		StartDate: mod.StartDate,
		EndDate:   mod.EndDate,
		CreatedAt: mod.CreatedAt,
		UserID:    mod.UserID,
	}

	if mod.NotifyAt != 0 {
		res.NotifyAt = mod.NotifyAt.String()
	}

	if mod.UpdatedAt.Valid {
		res.UpdatedAt = mod.UpdatedAt.Time
	}

	return res
}

func ToApiEventsInfo(events []*model.EventInfo) []*api.EventInfo {
	if events == nil {
		return nil
	}

	res := make([]*api.EventInfo, 0, len(events))
	for _, elem := range events {
		res = append(res, ToApiEventInfo(elem))
	}

	return res
}

func ToApiSuites(mod []*model.Suite) []*api.Suite {
	var res []*api.Suite
	for _, elem := range mod {
		res = append(res, &api.Suite{
			SuiteID:  elem.SuiteID,
			Capacity: elem.Capacity,
			Name:     elem.Name,
		})
	}

	return res
}

// Эта функция преобразует массив занятых интервалов к виду свободных
func ToFreeIntervals(mod []*model.Interval) []*api.Interval {
	now := time.Now()
	month := now.Add(720 * time.Hour)
	var res []*api.Interval

	if mod == nil {
		res = append(res, &api.Interval{
			StartDate: now,
			EndDate:   month,
		})
		return res
	}

	if now.Before(mod[0].StartDate) {
		res = append(res, &api.Interval{
			StartDate: now,
			EndDate:   mod[0].StartDate,
		})
	}

	if len(mod) == 1 && mod[0].EndDate.After(month) {
		return res
	}

	if len(mod) == 1 {
		res = append(res, &api.Interval{
			StartDate: mod[0].EndDate,
			EndDate:   month,
		})
		return res
	}

	for i := 1; i < len(mod); i++ {
		if mod[i].EndDate.Before(month) {
			res = append(res, &api.Interval{
				StartDate: mod[i-1].EndDate,
				EndDate:   mod[i].StartDate,
			})
		} else {
			res = append(res, &api.Interval{
				StartDate: mod[i-1].EndDate,
				EndDate:   mod[i].StartDate,
			})
			return res
		}

	}

	if mod[len(mod)-1].EndDate.Before(month) {
		res = append(res, &api.Interval{
			StartDate: mod[len(mod)-1].EndDate,
			EndDate:   month,
		})
	}

	return res
}

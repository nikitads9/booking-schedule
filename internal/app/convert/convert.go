package convert

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/model"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

func ToEvent(r *http.Request, req *api.Request) (*model.Event, error) {
	if req == nil {
		return nil, api.ErrEmptyRequest
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		return nil, api.ErrNoUserID
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}

	res := &model.Event{
		UserID:    id,
		SuiteID:   req.SuiteID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if r.Method == "PATCH" {
		eventID := chi.URLParam(r, "event_id")
		if eventID == "" {
			return nil, api.ErrNoEventID
		}

		eventUUID, err := uuid.FromString(eventID)
		if err != nil {
			return nil, api.ErrParse
		}

		if eventUUID == uuid.Nil {
			return nil, api.ErrNoEventID
		}

		res.SetEventID(eventUUID)
	}

	if req.NotificationPeriod.Valid {
		dur, err := time.ParseDuration(req.NotificationPeriod.String)
		if err != nil {
			return nil, err
		}
		res.SetNotifyAt(dur)
	}

	return res, nil
}

/* func ToUpdateEventInfo(r *http.Request, req *api.UpdateEventRequest) (*model.UpdateEventInfo, error) {
	if req == nil {
		return nil, api.ErrEmptyRequest
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		return nil, api.ErrNoUserID
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, api.ErrParse
	}

	eventID := chi.URLParam(r, "event_id")
	if eventID == "" {
		return nil, api.ErrNoEventID
	}

	eventUUID, err := uuid.FromString(eventID)
	if err != nil {
		return nil, api.ErrParse
	}

	if eventUUID == uuid.Nil {
		return nil, api.ErrNoEventID
	}

	if req.NotificationPeriod.Valid {
		dur, err := time.ParseDuration(req.NotificationPeriod.String)
		if err != nil {
			return nil, err
		}

		return &model.UpdateEventInfo{
			EventID:   eventUUID,
			UserID:    id,
			SuiteID:   req.SuiteID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			NotifyAt: null.Time{Time: req.StartDate.Add(-dur),
				Valid: true},
		}, nil
	}

	return &model.UpdateEventInfo{
		EventID:   eventUUID,
		UserID:    id,
		SuiteID:   req.SuiteID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
} */

func ToGetEventsInfo(r *http.Request) (*model.GetEventsInfo, error) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		return nil, api.ErrNoUserID
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, api.ErrParse
	}

	start := r.URL.Query().Get("start")
	if start == "" {
		return nil, api.ErrNoInterval
	}

	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, api.ErrParse
	}

	end := r.URL.Query().Get("end")
	if end == "" {
		return nil, api.ErrNoInterval
	}

	endDate, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, err
	}

	// Проверка, что обе даты еще не прошли
	if (startDate.Before(time.Now())) || (endDate.Before(time.Now())) {
		return nil, api.ErrExpiredDate
	}

	//проверка, что дата окончания не находится перед датой начала и не совпадает с ней
	if endDate.Sub(startDate) <= 0 {
		return nil, api.ErrInvalidInterval
	}

	return &model.GetEventsInfo{
		UserID:    id,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func ToGetRoomsInfo(r *http.Request) (*model.Interval, error) {
	start := r.URL.Query().Get("start")
	if start == "" {
		return nil, api.ErrNoInterval
	}
	end := r.URL.Query().Get("end")
	if end == "" {
		return nil, api.ErrNoInterval
	}

	startDate, err := time.Parse("2006-01-02T15:04:05-07:00", start)
	if err != nil {
		return nil, api.ErrParse
	}
	endDate, err := time.Parse("2006-01-02T15:04:05-07:00", end)
	if err != nil {
		return nil, api.ErrParse
	}

	// Проверка, что обе даты еще не прошли
	if (startDate.Before(time.Now())) || (endDate.Before(time.Now())) {
		return nil, api.ErrExpiredDate
	}

	//проверка, что дата окончания не находится перед датой начала и не совпадает с ней
	if endDate.Sub(startDate) <= 0 {
		return nil, api.ErrInvalidInterval
	}

	return &model.Interval{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// Эта функция преобразует массив занятых интервалов к виду свободных
func ToFreeIntervals(mod []*model.Interval) []*model.Interval {
	now := time.Now()
	month := now.Add(720 * time.Hour)
	var res []*model.Interval

	if mod == nil {
		res = append(res, &model.Interval{
			StartDate: now,
			EndDate:   month,
		})
		return res
	}

	if now.Before(mod[0].StartDate) {
		res = append(res, &model.Interval{
			StartDate: now,
			EndDate:   mod[0].StartDate,
		})
	}

	if len(mod) == 1 && mod[0].EndDate.After(month) {
		return res
	}

	if len(mod) == 1 {
		res = append(res, &model.Interval{
			StartDate: mod[0].EndDate,
			EndDate:   month,
		})
		return res
	}

	for i := 1; i < len(mod); i++ {
		if mod[i].EndDate.Before(month) {
			res = append(res, &model.Interval{
				StartDate: mod[i-1].EndDate,
				EndDate:   mod[i].StartDate,
			})
		} else {
			res = append(res, &model.Interval{
				StartDate: mod[i-1].EndDate,
				EndDate:   mod[i].StartDate,
			})
			return res
		}

	}

	if mod[len(mod)-1].EndDate.Before(month) {
		res = append(res, &model.Interval{
			StartDate: mod[len(mod)-1].EndDate,
			EndDate:   month,
		})
	}

	return res
}

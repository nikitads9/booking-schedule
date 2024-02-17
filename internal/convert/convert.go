package convert

import (
	"database/sql"
	"event-schedule/internal/api"
	"event-schedule/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

func ToEvent(r *http.Request, req *api.AddEventRequest) (*model.Event, error) {
	var dur time.Duration

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		return nil, api.ErrNoUserID
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}

	if req == nil {
		return nil, api.ErrEmptyRequest
	}

	if req.NotificationPeriod != "" {
		dur, err = time.ParseDuration(req.NotificationPeriod)
		if err != nil {
			return nil, err
		}

		return &model.Event{
			UserID:    id,
			SuiteID:   req.SuiteID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			NotifyAt: sql.NullTime{Time: req.StartDate.Add(-dur),
				Valid: true},
		}, nil
	}

	return &model.Event{
		UserID:    id,
		SuiteID:   req.SuiteID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func ToUpdateEventInfo(r *http.Request, req *api.UpdateEventRequest) (*model.UpdateEventInfo, error) {
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

	return &model.UpdateEventInfo{
		EventID:            eventUUID,
		UserID:             id,
		SuiteID:            req.SuiteID,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		NotificationPeriod: req.NotificationPeriod,
	}, nil
}

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

	return &model.Interval{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

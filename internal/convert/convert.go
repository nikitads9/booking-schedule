package convert

import (
	"event-schedule/internal/api"
	"event-schedule/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

func ToEvent(req *api.AddEventRequest, id int64) *model.Event {
	if req == nil {
		return nil
	}

	return &model.Event{
		UserID:             id,
		SuiteID:            req.SuiteID,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		NotificationPeriod: req.NotificationPeriod,
	}
}

func ToUpdateEventInfo(req *api.UpdateEventRequest, eventID uuid.UUID, userID int64) *model.UpdateEventInfo {
	if req == nil {
		return nil
	}

	return &model.UpdateEventInfo{
		EventID:            eventID,
		UserID:             userID,
		SuiteID:            req.SuiteID,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		NotificationPeriod: req.NotificationPeriod,
	}
}

func ToGetEventsInfo(r *http.Request) (*model.GetEventsInfo, error) {
	r.URL.Query().Get("user_id")
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		return nil, api.ErrNoUserID
	}

	/* 	start := chi.URLParam(r, "start")
	   	if start == "" {
	   		return nil, api.ErrNoInterval
	   	}

	   	end := chi.URLParam(r, "end")
	   	if end == "" {
	   		return nil, api.ErrNoInterval
	   	} */

	start := r.URL.Query().Get("start")
	if start == "" {
		return nil, api.ErrNoInterval
	}

	end := r.URL.Query().Get("end")
	if end == "" {
		return nil, api.ErrNoInterval
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, err
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

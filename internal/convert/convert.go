package convert

import (
	"event-schedule/internal/api"
	"event-schedule/internal/model"

	"github.com/gofrs/uuid"
)

func ToEventInfo(req *api.AddEventRequest) *model.Event {
	if req == nil {
		return nil
	}

	return &model.Event{
		SuiteID:            req.SuiteID,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		NotificationPeriod: req.NotificationPeriod,
	}
}

func ToUpdateEventInfo(req *api.UpdateEventRequest, eventID uuid.UUID) *model.UpdateEventInfo {
	if req == nil {
		return nil
	}

	return &model.UpdateEventInfo{
		EventID:            eventID,
		SuiteID:            req.SuiteID,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		NotificationPeriod: req.NotificationPeriod,
	}
}

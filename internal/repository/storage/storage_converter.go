package storage

import (
""
)

func (a *adapterImpl) toTimesheetDto(t *domain.Timesheet) *timesheet {
	if t == nil {
		return nil
	}

	dto := &timesheet{
		Id:        t.Id,
		Owner:     t.Owner,
		DateFrom:  t.DateFrom,
		DateTo:    t.DateTo,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	return dto
}

func (a *adapterImpl) toTimesheetDomain(dto *timesheet) *domain.Timesheet {
	if dto == nil {
		return nil
	}
	return &domain.Timesheet{
		Id:        dto.Id,
		Owner:     dto.Owner,
		DateFrom:  dto.DateFrom,
		DateTo:    dto.DateTo,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func (a *adapterImpl) toTimesheetsDomain(dtos []*timesheet) []*domain.Timesheet {
	var r []*domain.Timesheet
	for _, dto := range dtos {
		r = append(r, a.toTimesheetDomain(dto))
	}
	return r
}

func (a *adapterImpl) toEventDto(t *domain.Event) *event {
	if t == nil {
		return nil
	}

	dto := &event{
		Id:          t.Id,
		TimesheetId: t.TimesheetId,
		Subject:     t.Subject,
		WeekDay:     t.WeekDay,
		TimeStart:   t.TimeStart,
		TimeEnd:     t.TimeEnd,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	return dto
}

func (a *adapterImpl) toEventDomain(dto *event) *domain.Event {
	if dto == nil {
		return nil
	}
	return &domain.Event{
		Id:          dto.Id,
		TimesheetId: dto.TimesheetId,
		Subject:     dto.Subject,
		WeekDay:     dto.WeekDay,
		TimeStart:   dto.TimeStart,
		TimeEnd:     dto.TimeEnd,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func (a *adapterImpl) toEventsDomain(dtos []*event) []*domain.Event {
	var r []*domain.Event
	for _, dto := range dtos {
		r = append(r, a.toEventDomain(dto))
	}
	return r
}

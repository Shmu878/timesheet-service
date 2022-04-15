package domain

import (
	"context"
)

type TimetableStorage interface {
	// CreateTimetable creates a new Timetable
	CreateTimetable(ctx context.Context, Timesheet *Timesheet) (*Timesheet, error)
	// UpdateTimetable updates an existent Timetable
	UpdateTimetable(ctx context.Context, Timesheet *Timesheet) (*Timesheet, error)
	// DeleteTimetable creates a new Timetable
	DeleteTimetable(ctx context.Context, id string) error
	// GetTimetable retrieves a Timetable by id
	GetTimetable(ctx context.Context, id string) (bool, *Timesheet, error)
	// SearchTimetable retrieves a timesheet by owner
	SearchTimetable(ctx context.Context, rq *SearchTimesheetRequest) (bool, *SearchTimesheetsResponse, error)
}

type EventStorage interface {
	// CreateEvent creates a new event
	CreateEvent(ctx context.Context, Event *Event) (*Event, error)
	// UpdateEvent updates an existent event
	UpdateEvent(ctx context.Context, Event *Event) (*Event, error)
	// DeleteEvent creates a new event
	DeleteEvent(ctx context.Context, id string, timesheetId string) error
	// GetEvent retrieves a event by id
	GetEvent(ctx context.Context, id string) (bool, *Event, error)
	// SearchEvents retrieves events by timesheetId
	SearchEvents(ctx context.Context, id string) (bool, *SearchResponse, error)
}

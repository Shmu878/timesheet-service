package domain

import (
	"context"
	"time"
)

type Timesheet struct {
	Id        string    `json:"id"`        // Id - unique identifier
	Owner     string    `json:"owner"`     // Owner - id of user
	DateFrom  time.Time `json:"dateFrom"`  // DateFrom - start of the schedule
	DateTo    time.Time `json:"dateTo"`    // DateTo - end of the schedule
	CreatedAt time.Time `json:"createdAt"` // CreatedAt - date of creation
	UpdatedAt time.Time `json:"updatedAt"` // UpdatedAt - last update
}

type Event struct {
	Id          string    `json:"id"`          // Id - unique identifier
	TimesheetId string    `json:"timesheetId"` // TimesheetId - timesheet id
	Subject     string    `json:"subject"`     // Subject - name
	WeekDay     string    `json:"weekday"`     // WeekDay - doy of the week
	TimeStart   time.Time `json:"timeStart"`   // TimeStart - lesson start time
	TimeEnd     time.Time `json:"timeEnd"`     // TimeEnd - lesson end time
	CreatedAt   time.Time `json:"createdAt"`   // CreatedAt - date of creation
	UpdatedAt   time.Time `json:"updatedAt"`   // UpdatedAt - last update
}

// SearchTimesheetRequest request for timesheet search
type SearchTimesheetRequest struct {
	Owner          string
	DateFromSearch *time.Time
	DateToSearch   *time.Time
}

type SearchResponse struct {
	Events []*Event
}

type SearchTimesheetsResponse struct {
	Timesheets []*Timesheet
}

type TimesheetService interface {
	// Create creates a new Timetable
	Create(ctx context.Context, Timesheet *Timesheet) (*Timesheet, error)
	// Update updates an existent Timetable
	Update(ctx context.Context, Timesheet *Timesheet) (*Timesheet, error)
	// Get retrieves a Timetable by id
	Get(ctx context.Context, id string) (bool, *Timesheet, error)
	// Search retrieves a timesheet by owner
	Search(ctx context.Context, rq *SearchTimesheetRequest) (bool, *SearchTimesheetsResponse, error)
	// Delete deletes a Timetable
	Delete(ctx context.Context, id string) error
	// CreateEvent creates a new even
	CreateEvent(ctx context.Context, Event *Event) (*Event, error)
	// UpdateEvent updates an existent event
	UpdateEvent(ctx context.Context, Event *Event) (*Event, error)
	// GetEvent retrieves a event by id
	GetEvent(ctx context.Context, id string) (bool, *Event, error)
	// DeleteEvent deletes a event
	DeleteEvent(ctx context.Context, id string) error
	// SearchEvents retrieves events by timesheetId
	SearchEvents(ctx context.Context, id string) (bool, *SearchResponse, error)
}

package impl

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// sampleTimesheetImpl - implements SampleService interface
type timesheetImpl struct {
	storageT domain.TimetableStorage
	storageE domain.EventStorage
}

func NewTimesheetService(
	storageT domain.TimetableStorage,
	storageE domain.EventStorage,

) domain.TimesheetService {
	return &timesheetImpl{
		storageT: storageT,
		storageE: storageE,
	}
}

func (s *timesheetImpl) l() log.CLogger {
	return logger.L().Cmp("timesheet")
}

func (s *timesheetImpl) isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (s *timesheetImpl) Create(ctx context.Context, tt *domain.Timesheet) (*domain.Timesheet, error) {
	s.l().C(ctx).Mth("create").Dbg()

	// check owner
	if tt.Owner == "" {
		return nil, errors.ErrTimesheetOwnerIsEmpty(ctx)
	}
	if !s.isValidUUID(tt.Owner) {
		return nil, errors.ErrTimesheetAccessVerification(ctx)
	}

	// check timeFrom
	if tt.DateFrom.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	// check timeTo
	if tt.DateTo.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	now := time.Now().UTC()
	tt.Id = utils.NewId()
	tt.CreatedAt, tt.UpdatedAt = now, now

	// save to store
	tt, err := s.storageT.CreateTimetable(ctx, tt)
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (s *timesheetImpl) Update(ctx context.Context, tt *domain.Timesheet) (*domain.Timesheet, error) {
	s.l().C(ctx).Mth("update").Dbg()

	// validates id
	if tt.Id == "" {
		return nil, errors.ErrTimesheetIdIsEmpty(ctx)
	}

	if tt.DateFrom.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	if tt.DateTo.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	// retrieve stored sample by id
	found, stored, err := s.storageT.GetTimetable(ctx, tt.Id)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.ErrTimesheetIdIsEmpty(ctx)
	}

	// set updated params
	now := time.Now().UTC()
	stored.DateFrom = tt.DateFrom
	stored.DateTo = tt.DateTo
	stored.UpdatedAt = now

	// save to store
	_, err = s.storageT.UpdateTimetable(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *timesheetImpl) Get(ctx context.Context, id string) (bool, *domain.Timesheet, error) {
	s.l().C(ctx).Mth("get").Dbg()
	// validates id
	if id == "" {
		return false, nil, errors.ErrTimesheetIdIsEmpty(ctx)
	}
	return s.storageT.GetTimetable(ctx, id)
}

func (s *timesheetImpl) Search(ctx context.Context, rq *domain.SearchTimesheetRequest) (bool, *domain.SearchTimesheetsResponse, error) {
	s.l().C(ctx).Mth("search").Dbg()

	// validates owner
	if rq.Owner == "" {
		return false, nil, errors.ErrTimesheetOwnerIsEmpty(ctx)
	}

	return s.storageT.SearchTimetable(ctx, rq)
}

func (s *timesheetImpl) Delete(ctx context.Context, id string) error {
	s.l().C(ctx).Mth("delete").Dbg()

	// check id isn't empty
	if id == "" {
		return errors.ErrTimesheetIdIsEmpty(ctx)
	}

	// retrieve stored sample by id
	found, stored, err := s.storageT.GetTimetable(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return errors.ErrTimesheetNotFound(ctx, id)
	}

	// set updated params
	now := time.Now().UTC()
	stored.UpdatedAt = now

	// save to store
	return s.storageT.DeleteTimetable(ctx, id)
}

func (s *timesheetImpl) CreateEvent(ctx context.Context, ev *domain.Event) (*domain.Event, error) {
	s.l().C(ctx).Mth("create").Dbg()

	// check timesheetId
	if ev.TimesheetId == "" {
		return nil, errors.ErrTimesheetIdIsEmpty(ctx)
	}
	if !s.isValidUUID(ev.TimesheetId) {
		return nil, errors.ErrTimesheetAccessVerification(ctx)
	}

	// check timesheet isn't empty
	if ev.Subject == "" {
		return nil, errors.ErrTimesheetSubjectIsEmpty(ctx)
	}
	if ev.WeekDay == "" {
		return nil, errors.ErrTimesheetWeekDayIsEmpty(ctx)
	}
	if ev.TimeStart.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}
	if ev.TimeEnd.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	now := time.Now().UTC()
	ev.Id = utils.NewId()
	ev.CreatedAt, ev.UpdatedAt = now, now

	// save to store
	ev, err := s.storageE.CreateEvent(ctx, ev)
	if err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *timesheetImpl) UpdateEvent(ctx context.Context, ev *domain.Event) (*domain.Event, error) {
	s.l().C(ctx).Mth("update").Dbg()

	// validates id
	if ev.Id == "" {
		return nil, errors.ErrTimesheetIsEmpty(ctx)
	}

	// check timesheet isn't empty
	if ev.Subject == "" {
		return nil, errors.ErrTimesheetSubjectIsEmpty(ctx)
	}
	if ev.WeekDay == "" {
		return nil, errors.ErrTimesheetWeekDayIsEmpty(ctx)
	}
	if ev.TimeStart.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}
	if ev.TimeEnd.IsZero() {
		return nil, errors.ErrTimesheetTimeIsEmpty(ctx)
	}

	// retrieve stored sample by id
	found, stored, err := s.storageE.GetEvent(ctx, ev.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.ErrTimesheetNotFound(ctx, ev.Id)
	}

	// set updated params
	now := time.Now().UTC()
	stored.Subject = ev.Subject
	stored.WeekDay = ev.WeekDay
	stored.TimeStart = ev.TimeStart
	stored.TimeEnd = ev.TimeEnd
	stored.UpdatedAt = now

	// save to store
	_, err = s.storageE.UpdateEvent(ctx, stored)
	if err != nil {
		return nil, err
	}

	return stored, nil
}

func (s *timesheetImpl) GetEvent(ctx context.Context, id string) (bool, *domain.Event, error) {
	s.l().C(ctx).Mth("get").Dbg()
	if id == "" {
		return false, nil, errors.ErrTimesheetIsEmpty(ctx)
	}
	return s.storageE.GetEvent(ctx, id)
}

func (s *timesheetImpl) DeleteEvent(ctx context.Context, id string) error {
	s.l().C(ctx).Mth("delete").Dbg()

	// check id isn't empty
	if id == "" {
		return errors.ErrTimesheetIdIsEmpty(ctx)
	}

	// retrieve stored sample by id
	found, stored, err := s.storageE.GetEvent(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return errors.ErrTimesheetNotFound(ctx, id)
	}

	// set updated params
	now := time.Now().UTC()
	stored.UpdatedAt = now

	// save to store
	return s.storageE.DeleteEvent(ctx, id, stored.TimesheetId)
}

func (s *timesheetImpl) SearchEvents(ctx context.Context, id string) (bool, *domain.SearchResponse, error) {
	s.l().C(ctx).Mth("search").Dbg()

	// check id isn't empty
	if id == "" {
		return false, nil, errors.ErrTimesheetIsEmpty(ctx)
	}

	return s.storageE.SearchEvents(ctx, id)
}

//go:build integration
// +build integration

package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type timesheetStorageTestSuite struct {
	suite.Suite
	storageT domain.TimetableStorage
	storageE domain.EventStorage
	ctx      context.Context
}

// SetupSuite is called once for a suite
func (s *timesheetStorageTestSuite) SetupSuite() {

	// setup context
	s.ctx = kitContext.NewRequestCtx().Test().ToContext(context.Background())

	// load config
	cfg, err := config.Load()
	if err != nil {
		s.T().Fatal(err)
	}

	// disable applying migrations
	cfg.Storages.Database.MigPath = ""

	// initialize adapter
	a := NewAdapter()
	err = a.Init(cfg.Storages)
	if err != nil {
		s.T().Fatal(err)
	}
	s.storageT = a
	s.storageE = a

}

func (s *timesheetStorageTestSuite) getTimesheetRequest() *domain.Timesheet {
	now := time.Now().UTC()
	return &domain.Timesheet{
		Id:        kitUtils.NewId(),
		Owner:     "owner",
		DateFrom:  now,
		DateTo:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *timesheetStorageTestSuite) getEventRequest() *domain.Event {
	nowEvent := time.Now().UTC()
	return &domain.Event{
		Id:          kitUtils.NewId(),
		TimesheetId: kitUtils.NewId(),
		Subject:     "subject",
		WeekDay:     "monday",
		TimeStart:   nowEvent,
		TimeEnd:     nowEvent,
		CreatedAt:   nowEvent,
		UpdatedAt:   nowEvent,
	}
}

func (s *timesheetStorageTestSuite) getEventRequest2() *domain.Event {
	nowEvent := time.Now().UTC()
	return &domain.Event{
		Id:          kitUtils.NewId(),
		TimesheetId: "123456",
		Subject:     "subject",
		WeekDay:     "monday",
		TimeStart:   nowEvent,
		TimeEnd:     nowEvent,
		CreatedAt:   nowEvent,
		UpdatedAt:   nowEvent,
	}
}

func (s *timesheetStorageTestSuite) getEventRequest3() *domain.Event {
	nowEvent := time.Now().UTC()
	return &domain.Event{
		Id:          kitUtils.NewId(),
		TimesheetId: "123456",
		Subject:     "subject2",
		WeekDay:     "monday2",
		TimeStart:   nowEvent,
		TimeEnd:     nowEvent,
		CreatedAt:   nowEvent,
		UpdatedAt:   nowEvent,
	}
}

func (s *timesheetStorageTestSuite) getSearchRequest() *domain.SearchTimesheetRequest {
	nowRequest := time.Now().UTC()
	return &domain.SearchTimesheetRequest{
		Owner:          "owner",
		DateFromSearch: &nowRequest,
		DateToSearch:   &nowRequest,
	}
}

func TestSuiteTimesheet(t *testing.T) {
	suite.Run(t, new(timesheetStorageTestSuite))
}

func (s *timesheetStorageTestSuite) Test_CreateTimesheet_GetFromDbAndCache() {

	// create a task
	expected := s.getTimesheetRequest()
	_, err := s.storageT.CreateTimetable(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Sample by id
	_, actual, err := s.storageT.GetTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.Equal(s.T(), expected.Id, actual.Id)
	assert.Equal(s.T(), expected.Owner, actual.Owner)

	// wait for async caching
	time.Sleep(time.Millisecond * 100)

	// get Task by id again (cache hit)
	_, actual, err = s.storageT.GetTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.Equal(s.T(), expected.Id, actual.Id)
	assert.Equal(s.T(), expected.Owner, actual.Owner)
}

func (s *timesheetStorageTestSuite) TestUpdateTimesheet() {

	// create a sample
	expected := s.getTimesheetRequest()
	_, err := s.storageT.CreateTimetable(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageT.GetTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)

	// update sample
	actual.Owner = "ownerNew"

	_, err = s.storageT.UpdateTimetable(s.ctx, actual)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Sample by id
	_, actual, err = s.storageT.GetTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)
	assert.Equal(s.T(), actual.Owner, "ownerNew")
}

func (s *timesheetStorageTestSuite) TestDeleteTimesheet() {

	// create a sample
	expected := s.getTimesheetRequest()
	_, err := s.storageT.CreateTimetable(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageT.GetTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)

	// delete sample
	err = s.storageT.DeleteTimetable(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *timesheetStorageTestSuite) TestSearchTimesheet() {

	// create a sample
	expected := s.getTimesheetRequest()
	searchExpected := s.getSearchRequest()
	_, err := s.storageT.CreateTimetable(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageT.SearchTimetable(s.ctx, searchExpected)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Timesheets)
}

func (s *timesheetStorageTestSuite) Test_CreateEvent_GetFromDbAndCache() {

	// create a task
	expected := s.getEventRequest()
	_, err := s.storageE.CreateEvent(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Sample by id
	_, actual, err := s.storageE.GetEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.Equal(s.T(), expected.Id, actual.Id)
	assert.Equal(s.T(), expected.TimesheetId, actual.TimesheetId)
	assert.Equal(s.T(), expected.Subject, actual.Subject)
	assert.Equal(s.T(), expected.WeekDay, actual.WeekDay)

	// wait for async caching
	time.Sleep(time.Millisecond * 100)

	// get Task by id again (cache hit)
	_, actual, err = s.storageE.GetEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.Equal(s.T(), expected.Id, actual.Id)
	assert.Equal(s.T(), expected.TimesheetId, actual.TimesheetId)
	assert.Equal(s.T(), expected.Subject, actual.Subject)
	assert.Equal(s.T(), expected.WeekDay, actual.WeekDay)
}

func (s *timesheetStorageTestSuite) TestUpdateEvent() {

	// create a sample
	expected := s.getEventRequest()
	_, err := s.storageE.CreateEvent(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageE.GetEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)

	// update sample
	actual.Subject = "subjectNew"
	actual.WeekDay = "monday"
	_, err = s.storageE.UpdateEvent(s.ctx, actual)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Sample by id
	_, actual, err = s.storageE.GetEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)
	assert.Equal(s.T(), actual.Subject, "subjectNew")
}

func (s *timesheetStorageTestSuite) TestDeleteEvent() {

	// create a sample
	expected := s.getEventRequest()
	_, err := s.storageE.CreateEvent(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageE.GetEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Id)

	// delete sample
	err = s.storageE.DeleteEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *timesheetStorageTestSuite) TestSearchEvents() {

	// create a sample
	expected := s.getEventRequest2()
	expected2 := s.getEventRequest3()
	_, err := s.storageE.CreateEvent(s.ctx, expected)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = s.storageE.CreateEvent(s.ctx, expected2)
	if err != nil {
		s.T().Fatal(err)
	}

	// get Task by id
	_, actual, err := s.storageE.SearchEvents(s.ctx, expected.TimesheetId)
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), actual.Events)

	err = s.storageE.DeleteEvent(s.ctx, expected.Id)
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.storageE.DeleteEvent(s.ctx, expected2.Id)
	if err != nil {
		s.T().Fatal(err)
	}
}

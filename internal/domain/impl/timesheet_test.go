package impl

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type timesheetTestSuite struct {
	suite.Suite
	timetableStorage *mocks.TimetableStorage
	eventStorage     *mocks.EventStorage
	timesheetService domain.TimesheetService
	ctx              context.Context
}

func (s *timesheetTestSuite) SetupSuite() {
	s.timetableStorage = &mocks.TimetableStorage{}
	s.eventStorage = &mocks.EventStorage{}
	s.ctx = kitContext.NewRequestCtx().Test().ToContext(context.Background())
	s.timesheetService = NewTimesheetService(s.timetableStorage, s.eventStorage)
}

func (s *timesheetTestSuite) SetupTest() {
	s.timetableStorage.ExpectedCalls = nil
	s.timetableStorage.On("CreateTimetable",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*domain.Timesheet")).
		Return(&domain.Timesheet{
			Id:        kitUtils.NewId(),
			Owner:     "123",
			DateFrom:  time.Now(),
			DateTo:    time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
	s.timetableStorage.On("UpdateTimetable",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*domain.Timesheet")).
		Return(&domain.Timesheet{
			Id:        kitUtils.NewId(),
			Owner:     "123",
			DateFrom:  time.Now(),
			DateTo:    time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
	s.timetableStorage.On("GetTimetable",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string")).
		Return(true, &domain.Timesheet{
			Id:        kitUtils.NewId(),
			Owner:     "123",
			DateFrom:  time.Now(),
			DateTo:    time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
	s.timetableStorage.On("SearchTimetable",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*domain.SearchTimesheetRequest")).
		Return(true, &domain.SearchTimesheetsResponse{Timesheets: []*domain.Timesheet{{
			Id:        kitUtils.NewId(),
			Owner:     "123",
			DateFrom:  time.Time{},
			DateTo:    time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}}}, nil)
	s.timetableStorage.On("DeleteTimetable",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string")).
		Return(nil)

	s.eventStorage.ExpectedCalls = nil
	s.eventStorage.On("CreateEvent",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*domain.Event")).
		Return(&domain.Event{
			Id:          kitUtils.NewId(),
			TimesheetId: kitUtils.NewId(),
			Subject:     "test subject",
			WeekDay:     "Monday",
			TimeStart:   time.Now(),
			TimeEnd:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)
	s.eventStorage.On("UpdateEvent",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*domain.Event")).
		Return(&domain.Event{
			Id:          kitUtils.NewId(),
			TimesheetId: kitUtils.NewId(),
			Subject:     "test subject 2",
			WeekDay:     "Monday",
			TimeStart:   time.Now(),
			TimeEnd:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)
	s.eventStorage.On("GetEvent",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string")).
		Return(true, &domain.Event{
			Id:          kitUtils.NewId(),
			TimesheetId: kitUtils.NewId(),
			Subject:     "test",
			WeekDay:     "Monday",
			TimeStart:   time.Now(),
			TimeEnd:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)
	s.eventStorage.On("DeleteEvent",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string")).
		Return(nil)
	s.eventStorage.On("SearchEvents",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string")).
		Return(true, &domain.SearchResponse{Events: []*domain.Event{{
			Id:          kitUtils.NewId(),
			TimesheetId: "123",
			Subject:     "test",
			WeekDay:     "Monday",
			TimeStart:   time.Now(),
			TimeEnd:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, {
			Id:          kitUtils.NewId(),
			TimesheetId: "123",
			Subject:     "test2",
			WeekDay:     "Monday2",
			TimeStart:   time.Now(),
			TimeEnd:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}}}, nil)
}

func TestTimesheetSuite(t *testing.T) {
	suite.Run(t, new(timesheetTestSuite))
}

func (s *timesheetTestSuite) Test_Create_WhenEmptyName_Fail() {
	_, err := s.timesheetService.Create(s.ctx, &domain.Timesheet{})
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetOwnerEmpty)
}

func (s *timesheetTestSuite) Test_Create_Ok() {

	tt, err := s.timesheetService.Create(s.ctx, &domain.Timesheet{
		Owner:    kitUtils.NewId(),
		DateFrom: time.Now(),
		DateTo:   time.Now(),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), tt)
	assert.NotEmpty(s.T(), tt.Id)
	assert.NotEmpty(s.T(), tt.Owner)
	assert.NotEmpty(s.T(), tt.DateFrom)
	assert.NotEmpty(s.T(), tt.DateTo)
	assert.NotEmpty(s.T(), tt.CreatedAt)
	assert.NotEmpty(s.T(), tt.UpdatedAt)
}

func (s *timesheetTestSuite) Test_Update_Fail() {
	_, err := s.timesheetService.Update(s.ctx, &domain.Timesheet{})
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetIdEmpty)
}

func (s *timesheetTestSuite) Test_Update_Ok() {
	_, err := s.timesheetService.Update(s.ctx, &domain.Timesheet{
		Id:       kitUtils.NewId(),
		Owner:    "123",
		DateFrom: time.Now(),
		DateTo:   time.Now(),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)

}

func (s *timesheetTestSuite) Test_Get_Fail() {
	_, _, err := s.timesheetService.Get(s.ctx, "")
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetIdEmpty)
}

func (s *timesheetTestSuite) Test_Get_Ok() {
	_, tt, err := s.timesheetService.Get(s.ctx, kitUtils.NewId())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tt)
}

func (s *timesheetTestSuite) Test_Search_Fail() {
	_, _, err := s.timesheetService.Search(s.ctx, &domain.SearchTimesheetRequest{})
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetOwnerEmpty)
}

func (s *timesheetTestSuite) Test_Search_Ok() {
	_, tt, err := s.timesheetService.Search(s.ctx, &domain.SearchTimesheetRequest{
		Owner: "123",
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tt)
}

func (s *timesheetTestSuite) Test_Delete_Ok() {
	err := s.timesheetService.Delete(s.ctx, kitUtils.NewId())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
}

func (s *timesheetTestSuite) Test_CreateEvent_WhenEmptyName_Fail() {
	_, err := s.timesheetService.CreateEvent(s.ctx, &domain.Event{})
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetIdEmpty)
}

func (s *timesheetTestSuite) Test_CreateEvent_Ok() {

	ev, err := s.timesheetService.CreateEvent(s.ctx, &domain.Event{
		Id:          kitUtils.NewId(),
		TimesheetId: kitUtils.NewId(),
		Subject:     "test subject",
		WeekDay:     "Monday",
		TimeStart:   time.Now(),
		TimeEnd:     time.Now(),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), ev)
	assert.NotEmpty(s.T(), ev.Id)
	assert.NotEmpty(s.T(), ev.TimesheetId)
	assert.NotEmpty(s.T(), ev.Subject)
	assert.NotEmpty(s.T(), ev.WeekDay)
	assert.NotEmpty(s.T(), ev.TimeStart)
	assert.NotEmpty(s.T(), ev.TimeEnd)
	assert.NotEmpty(s.T(), ev.CreatedAt)
	assert.NotEmpty(s.T(), ev.UpdatedAt)
}

func (s *timesheetTestSuite) Test_UpdateEvent_Fail() {
	_, err := s.timesheetService.UpdateEvent(s.ctx, &domain.Event{})
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetNotFound)
}

func (s *timesheetTestSuite) Test_UpdateEvent_Ok() {
	_, err := s.timesheetService.UpdateEvent(s.ctx, &domain.Event{
		Id:          kitUtils.NewId(),
		TimesheetId: kitUtils.NewId(),
		Subject:     "test subject 2",
		WeekDay:     "Monday",
		TimeStart:   time.Now(),
		TimeEnd:     time.Now(),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
}

func (s *timesheetTestSuite) Test_GetEvent_Fail() {
	_, _, err := s.timesheetService.GetEvent(s.ctx, "")
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetNotFound)
}

func (s *timesheetTestSuite) Test_GetEvent_Ok() {
	_, tt, err := s.timesheetService.GetEvent(s.ctx, kitUtils.NewId())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tt)
}

func (s *timesheetTestSuite) Test_DeleteEvent_Ok() {
	err := s.timesheetService.DeleteEvent(s.ctx, kitUtils.NewId())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
}

func (s *timesheetTestSuite) Test_SearchEvents_Fail() {
	_, _, err := s.timesheetService.SearchEvents(s.ctx, "")
	kitTest.AssertAppErr(s.T(), err, pb.ErrCodeTimesheetNotFound)
}

func (s *timesheetTestSuite) Test_SearchEvents_Ok() {
	_, tt, err := s.timesheetService.SearchEvents(s.ctx, kitUtils.NewId())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tt)
}

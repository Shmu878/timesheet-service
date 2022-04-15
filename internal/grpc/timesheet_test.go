//go:build integration
// +build integration

package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

type timesheetGrpcTestSuite struct {
	suite.Suite
	ctx         context.Context
	clTimesheet pb.TimesheetServiceClient
}

func (s *timesheetGrpcTestSuite) SetupSuite() {

	// setup context
	s.ctx = kitContext.NewRequestCtx().Test().ToContext(context.Background())

	// load config
	cfg, err := config.Load()
	if err != nil {
		s.T().Fatal(err)
	}

	// create GRPC client
	cl, err := kitGrpc.NewClient(&kitGrpc.ClientConfig{Host: cfg.Grpc.Host, Port: cfg.Grpc.Port})
	if err != nil {
		s.T().Fatal(err)
	}
	s.clTimesheet = pb.NewTimesheetServiceClient(cl.Conn)
}

func TestTimesheetSuite(t *testing.T) {
	suite.Run(t, new(timesheetGrpcTestSuite))
}

func (s *timesheetGrpcTestSuite) getCreateTimesheetRequest() *pb.CreateTimesheetRequest {

	return &pb.CreateTimesheetRequest{
		Owner:    kitUtils.NewId(),
		DateFrom: timestamppb.Now(),
		DateTo:   timestamppb.Now(),
	}
}

func (s *timesheetGrpcTestSuite) getCreateEventRequest() *pb.CreateEventRequest {

	return &pb.CreateEventRequest{
		TimesheetId: kitUtils.NewId(),
		Subject:     "language",
		WeekDay:     "Monday",
		TimeStart:   timestamppb.Now(),
		TimeEnd:     timestamppb.Now(),
	}
}

func (s *timesheetGrpcTestSuite) TestTimesheetCRUD() {

	// create a new consultant
	cl, err := s.clTimesheet.Create(s.ctx, s.getCreateTimesheetRequest())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), cl.Id)

	// get by id
	cl, err = s.clTimesheet.Get(s.ctx, &pb.TimesheetIdRequest{Id: cl.Id})
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), cl)
	assert.NotEmpty(s.T(), cl.Id)

	// set status to active
	cl, err = s.clTimesheet.Update(s.ctx, &pb.UpdateTimesheetRequest{
		Id:       cl.Id,
		Owner:    cl.Owner,
		DateFrom: timestamppb.Now(),
		DateTo:   timestamppb.Now(),
	})
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), cl.Owner)

	//search
	sl, err := s.clTimesheet.Search(s.ctx, &pb.SearchTimesheetRequest{
		Owner:          kitUtils.NewId(),
		DateFromSearch: timestamppb.Now(),
		DateToSearch:   timestamppb.Now(),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), sl)

	// delete sample
	_, err = s.clTimesheet.Delete(s.ctx, &pb.TimesheetIdRequest{Id: cl.Id})
	if err != nil {
		s.T().Fatal()
	}
}

func (s *timesheetGrpcTestSuite) TestEventCRUD() {

	// create a new consultant
	cl, err := s.clTimesheet.CreateEvent(s.ctx, s.getCreateEventRequest())
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), cl.Id)

	// get by Id
	cl, err = s.clTimesheet.GetEvent(s.ctx, &pb.EventIdRequest{Id: cl.Id})
	if err != nil {
		s.T().Fatal()
	}
	assert.NotEmpty(s.T(), cl)
	assert.NotEmpty(s.T(), cl.Id)
	assert.NotEmpty(s.T(), cl.WeekDay)

	// set status to active
	cl, err = s.clTimesheet.UpdateEvent(s.ctx, &pb.UpdateEventRequest{
		Id:          cl.Id,
		TimesheetId: kitUtils.NewId(),
		Subject:     "language",
		WeekDay:     "Monday",
		TimeStart:   timestamppb.Now(),
		TimeEnd:     timestamppb.Now(),
	})
	if err != nil {
		s.T().Fatal()
	}
	assert.Equal(s.T(), "language", cl.Subject)

	// search Events by timesheetId
	sl, err := s.clTimesheet.SearchEvents(s.ctx, &pb.TimesheetIdRequest{Id: cl.TimesheetId})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), sl)

	// delete Event
	_, err = s.clTimesheet.DeleteEvent(s.ctx, &pb.EventIdRequest{Id: cl.Id})
	if err != nil {
		s.T().Fatal(err)
	}
}

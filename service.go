// this file is generated by servgen util based on a template at 2021-06-26 10:37:24 +0300 MSK
package timesheet

import (
	"context"
	""
)

// serviceImpl implements a service bootstrapping
// all dependencies between layers must be specified here
type serviceImpl struct {
	service.Cluster
	cfg              *config.Config
	monitoring       kitMonitoring.MetricsServer
	timesheetService domain.TimesheetService
	grpc             *grpc.Server
	storageAdapterT  storage.Adapter
	storageAdapterE  storage.Adapter
}

// New creates a new instance of the service
func New() service.Service {

	s := &serviceImpl{
		Cluster:    service.NewCluster(logger.LF(), meta.Meta),
		monitoring: kitMonitoring.NewMetricsServer(logger.LF()),
	}

	s.storageAdapterT = storage.NewAdapter()
	s.storageAdapterE = storage.NewAdapter()

	s.timesheetService = impl.NewTimesheetService(s.storageAdapterT, s.storageAdapterE)

	s.grpc = grpc.New(s.timesheetService)

	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.Meta.ServiceCode()
}

// Init does all initializations
func (s *serviceImpl) Init(ctx context.Context) error {

	// load config
	var err error
	s.cfg, err = config.Load()
	if err != nil {
		return err
	}

	// set log config
	logger.Logger.Init(s.cfg.Log)

	// init cluster
	if err := s.Cluster.Init(s.cfg.Cluster, s.cfg.Nats.Host, s.cfg.Nats.Port, s.onClusterLeaderChanged(ctx)); err != nil {
		return err
	}

	// init storage
	if err := s.storageAdapterE.Init(s.cfg.Storages); err != nil {
		return err
	}
	if err := s.storageAdapterT.Init(s.cfg.Storages); err != nil {
		return err
	}
	// init grpc server
	if err := s.grpc.Init(s.cfg.Grpc); err != nil {
		return err
	}

	if s.cfg.Monitoring.Enabled {
		if err := s.monitoring.Init(s.cfg.Monitoring); err != nil {
			return err
		}
	}

	return nil

}

func (s *serviceImpl) onClusterLeaderChanged(ctx context.Context) service.OnLeaderChangedEvent {

	// if the current node is getting a leader, run daemons
	return func(l bool) {
		if l {
			// do something if the node is turned into a leader
			logger.L().C(ctx).Cmp("cluster").Mth("on-leader-change").Dbg("leader")
		}
	}

}

func (s *serviceImpl) Start(ctx context.Context) error {

	// start cluster
	if err := s.Cluster.Start(); err != nil {
		return err
	}

	// serve gRPC connection
	s.grpc.ListenAsync()

	if s.cfg.Monitoring.Enabled {
		s.monitoring.Listen()
	}

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.Cluster.Close()
	s.storageAdapterT.Close()
	s.storageAdapterE.Close()
	s.grpc.Close()
	if s.cfg.Monitoring.Enabled {
		s.monitoring.Close()
	}
}

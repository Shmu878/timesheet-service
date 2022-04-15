package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"time"
)

const (
	CacheKeyTimesheetId       = "timesheet.id:"
	CacheKeyEventId           = "event.id:"
	CacheKeyEventsTimesheetId = "event.timesheet.id:"
)

type timesheet struct {
	Id        string         `gorm:"column:id"`
	Owner     string         `gorm:"column:owner"`
	DateFrom  time.Time      `gorm:"column:date_from"`
	DateTo    time.Time      `gorm:"column:date_to"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type event struct {
	Id          string         `gorm:"column:id"`
	TimesheetId string         `gorm:"column:timesheet_id"`
	Subject     string         `gorm:"column:subject"`
	WeekDay     string         `gorm:"column:weekday"`
	TimeStart   time.Time      `gorm:"column:time_start"`
	TimeEnd     time.Time      `gorm:"column:time_end"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (a *adapterImpl) l() log.CLogger {
	return logger.L().Cmp("timesheet-storage")
}

func (a *adapterImpl) setEventTimesheetCacheAsync(ctx context.Context, key string, dto interface{}) {

	go func() {

		l := a.l().Mth("set-cache").C(ctx).Dbg()

		j, err := json.Marshal(dto)
		if err != nil {
			l.E(err).St().Err()
		}
		dtoStr := string(j)

		// set cache for id key
		if err := a.container.Cache.Instance.Set(key, dtoStr, time.Hour).Err(); err != nil {
			l.E(errors.ErrTimesheetStorageSetCache(err, ctx, key)).St().Err()
		}
	}()
}

func (a *adapterImpl) CreateTimetable(ctx context.Context, timesheet *domain.Timesheet) (*domain.Timesheet, error) {
	a.l().C(ctx).Mth("create")
	// save to DB
	if err := a.container.Db.Instance.Create(a.toTimesheetDto(timesheet)).Error; err != nil {
		return nil, errors.ErrTimesheetStorageCreate(err, ctx)
	}

	return timesheet, nil
}

func (a *adapterImpl) UpdateTimetable(ctx context.Context, timesheet *domain.Timesheet) (*domain.Timesheet, error) {
	a.l().Mth("update").C(ctx).Dbg()

	// update DB
	if err := a.container.Db.Instance.Save(a.toTimesheetDto(timesheet)).Error; err != nil {
		return nil, errors.ErrTimesheetStorageUpdate(err, ctx, timesheet.Id)
	}

	// clear cache
	keys := []string{CacheKeyTimesheetId + timesheet.Id}
	a.container.Cache.Instance.Del(keys...)

	return timesheet, nil
}

func (a *adapterImpl) GetTimetable(ctx context.Context, id string) (bool, *domain.Timesheet, error) {
	l := a.l().Mth("get").C(ctx).F(log.FF{"id": id}).Dbg()

	key := CacheKeyTimesheetId + id
	if j, err := a.container.Cache.Instance.Get(key).Result(); err == nil {
		// found in cache
		l.Dbg("found in cache")
		dto := &timesheet{}
		if err := json.Unmarshal([]byte(j), &dto); err != nil {
			return true, nil, err
		}
		return true, a.toTimesheetDomain(dto), nil
	} else {
		if err == redis.Nil {
			// not found in cache
			dto := &timesheet{Id: id}
			if res := a.container.Db.Instance.Limit(1).Find(&dto); res.Error == nil {
				l.DbgF("db: found %d", res.RowsAffected)
				if res.RowsAffected == 0 {
					return false, nil, nil
				} else {
					// set cache async
					a.setEventTimesheetCacheAsync(ctx, key, dto)
					return true, a.toTimesheetDomain(dto), nil
				}
			} else {
				return false, nil, errors.ErrTimesheetStorageGetDb(res.Error, ctx, id)
			}

		} else {
			return false, nil, errors.ErrTimesheetStorageGetCache(err, ctx, id)
		}
	}
}

func (a *adapterImpl) SearchTimetable(ctx context.Context, rq *domain.SearchTimesheetRequest) (bool, *domain.SearchTimesheetsResponse, error) {
	a.l().Mth("search").C(ctx).F(log.FF{"owner": rq.Owner}).Dbg()
	res := []*timesheet{}

	db := a.container.Db.Instance.Where(&timesheet{
		Owner: rq.Owner,
	})

	if rq.DateFromSearch != nil {
		db.Where("date_from < ?", rq.DateToSearch)
	}
	if rq.DateToSearch != nil {
		db.Where("date_to > ?", rq.DateFromSearch)
	}

	if err := db.Find(&res).Error; err != nil {
		return false, nil, errors.ErrTimesheetIndexSearch(err, ctx)
	}

	return true, &domain.SearchTimesheetsResponse{Timesheets: a.toTimesheetsDomain(res)}, nil
}

func (a *adapterImpl) DeleteTimetable(ctx context.Context, id string) error {
	a.l().C(ctx).Mth("create")
	// delete to DB
	if err := a.container.Db.Instance.Delete(&timesheet{Id: id}).Error; err != nil {
		return errors.ErrTimesheetStorageDelete(err, ctx, id)
	}

	keys := []string{CacheKeyTimesheetId + id}
	a.container.Cache.Instance.Del(keys...)

	return nil
}

func (a *adapterImpl) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	a.l().C(ctx).Mth("create")
	// save to DB
	if err := a.container.Db.Instance.Create(a.toEventDto(event)).Error; err != nil {
		return nil, errors.ErrTimesheetStorageCreate(err, ctx)
	}

	// clear cache
	keys := []string{CacheKeyEventsTimesheetId + event.TimesheetId}
	a.container.Cache.Instance.Del(keys...)

	return event, nil
}

func (a *adapterImpl) UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	a.l().Mth("update").C(ctx).Dbg()

	// update DB
	if err := a.container.Db.Instance.Save(a.toEventDto(event)).Error; err != nil {
		return nil, errors.ErrTimesheetStorageUpdate(err, ctx, event.Id)
	}

	// clear cache
	keys := []string{CacheKeyEventId + event.Id, CacheKeyEventsTimesheetId + event.TimesheetId}
	a.container.Cache.Instance.Del(keys...)

	return event, nil
}

func (a *adapterImpl) DeleteEvent(ctx context.Context, id string, timesheetId string) error {
	a.l().C(ctx).Mth("create")

	// delete to DB
	if err := a.container.Db.Instance.Delete(&event{Id: id}).Error; err != nil {
		return errors.ErrTimesheetStorageDelete(err, ctx, id)
	}
	keys := []string{CacheKeyEventId + id, CacheKeyEventsTimesheetId + timesheetId}
	a.container.Cache.Instance.Del(keys...)

	return nil
}

func (a *adapterImpl) GetEvent(ctx context.Context, id string) (bool, *domain.Event, error) {
	l := a.l().Mth("get").C(ctx).F(log.FF{"id": id}).Dbg()

	key := CacheKeyEventId + id
	if j, err := a.container.Cache.Instance.Get(key).Result(); err == nil {
		// found in cache
		l.Dbg("found in cache")
		dto := &event{}
		if err := json.Unmarshal([]byte(j), &dto); err != nil {
			return false, nil, err
		}
		return true, a.toEventDomain(dto), nil
	} else {
		if err == redis.Nil {
			// not found in cache
			dto := &event{Id: id}
			if res := a.container.Db.Instance.Limit(1).Find(&dto); res.Error == nil {
				l.DbgF("db: found %d", res.RowsAffected)
				if res.RowsAffected == 0 {
					return false, nil, nil
				} else {
					// set cache async
					a.setEventTimesheetCacheAsync(ctx, key, dto)
					return true, a.toEventDomain(dto), nil
				}
			} else {
				return false, nil, errors.ErrTimesheetStorageGetDb(res.Error, ctx, id)
			}

		} else {
			return false, nil, errors.ErrTimesheetStorageGetCache(err, ctx, id)
		}
	}
}

func (a *adapterImpl) SearchEvents(ctx context.Context, id string) (bool, *domain.SearchResponse, error) {
	l := a.l().Mth("get").C(ctx).F(log.FF{"timesheetId": id}).Dbg()

	key := CacheKeyEventsTimesheetId + id
	if j, err := a.container.Cache.Instance.Get(key).Result(); err == nil {
		var dtos []*event
		// found in cache
		l.Dbg("found in cache")
		if err := json.Unmarshal([]byte(j), &dtos); err != nil {
			return false, nil, err
		}
		return true, &domain.SearchResponse{Events: a.toEventsDomain(dtos)}, nil
	} else {
		if err == redis.Nil {
			// not found in cache
			var dtos []*event
			if res := a.container.Db.Instance.Find(&dtos, "timesheet_id = ?", id); res.Error == nil {
				l.DbgF("db: found %d", res.RowsAffected)
				if res.RowsAffected == 0 {
					return false, nil, nil
				} else {
					// set cache async
					a.setEventTimesheetCacheAsync(ctx, key, dtos)
					return true, &domain.SearchResponse{Events: a.toEventsDomain(dtos)}, nil
				}
			} else {
				return false, nil, errors.ErrTimesheetStorageGetDb(res.Error, ctx, id)
			}

		} else {
			return false, nil, errors.ErrTimesheetStorageGetCache(err, ctx, id)
		}
	}
}

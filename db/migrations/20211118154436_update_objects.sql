-- +goose Up
-- +goose StatementBegin
set schema 'timesheet';

ALTER TABLE timesheets
    DROP subject;
ALTER TABLE timesheets
    ADD date_from timestamp;
ALTER TABLE timesheets
    ADD date_to timestamp;
ALTER TABLE events
    DROP due_date;
ALTER TABLE events
    DROP calendar_id;
ALTER TABLE events
    DROP description;
ALTER TABLE events
    ADD weekday varchar;
ALTER TABLE events
    ADD timesheet_id varchar;
ALTER TABLE events
    ADD time_start timestamp;
ALTER TABLE events
    ADD time_end timestamp;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'timesheet';

ALTER TABLE timesheets
    ADD subject varchar;
ALTER TABLE  timesheets
    DROP date_from;
ALTER TABLE  timesheets
    DROP date_to;
ALTER TABLE events
    ADD due_date timestamp;
ALTER TABLE events
    ADD calendar_id varchar;
ALTER TABLE events
    ADD description varchar;
ALTER TABLE events
    DROP weekday;
ALTER TABLE events
    DROP timesheet_id;
ALTER TABLE events
    DROP time_start;
ALTER TABLE events
    DROP time_end;

-- +goose StatementEnd
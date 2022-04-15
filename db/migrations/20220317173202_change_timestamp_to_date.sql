-- +goose Up
-- +goose StatementBegin
alter table timesheets
alter column date_from type date using date_from::date;

alter table timesheets
alter column date_to type date using date_to::date;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'timesheet';
-- +goose StatementEnd

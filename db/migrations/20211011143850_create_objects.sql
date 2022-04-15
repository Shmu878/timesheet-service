-- +goose Up
-- +goose StatementBegin
set schema 'timesheet';

create table timesheets (
    id uuid primary key,
    owner varchar not null,
    subject varchar not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

create table events (
    id uuid primary key,
    subject varchar not null,
    due_date timestamp not null,
    calendar_id varchar not null,
    description varchar,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'timesheet';

drop table timesheets;
drop table events;

-- +goose StatementEnd
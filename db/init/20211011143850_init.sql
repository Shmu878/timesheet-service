-- +goose Up
CREATE ROLE timesheet LOGIN PASSWORD 'timesheet' NOINHERIT CREATEDB;
CREATE SCHEMA timesheet AUTHORIZATION timesheet;
GRANT USAGE ON SCHEMA timesheet TO PUBLIC;

-- +goose Down
DROP SCHEMA timesheet;
DROP ROLE timesheet;

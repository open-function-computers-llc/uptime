-- +goose Up
-- +goose StatementBegin
ALTER TABLE sites ADD COLUMN meta TEXT NOT NULL DEFAULT '' AFTER is_deleted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sites DROP COLUMN meta;
-- +goose StatementEnd

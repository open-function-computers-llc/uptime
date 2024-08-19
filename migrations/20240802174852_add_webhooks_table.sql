-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `webhooks` (
    id int(11) PRIMARY KEY AUTO_INCREMENT,
    website_id int(11) NOT NULL,
    hook_name VARCHAR(255) NOT NULL,
    hook_url VARCHAR(255) NOT NULL,
    hook_verb CHAR(4),
    hook_type tinyint(3) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `webhooks`;
-- +goose StatementEnd

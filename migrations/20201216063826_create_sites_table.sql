-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `sites` (
    `id` int(11) PRIMARY KEY AUTO_INCREMENT,
    `url` VARCHAR(255) NOT NULL DEFAULT "",
    `is_up` int(11) DEFAULT 1,
    `is_deleted` int(11) DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
    `last_checked` DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
    CONSTRAINT url_unique UNIQUE (url)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `sites`;
-- +goose StatementEnd

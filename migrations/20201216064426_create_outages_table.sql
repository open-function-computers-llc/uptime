-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `outages` (
    id int(11) PRIMARY KEY AUTO_INCREMENT,
    website_id int(11) NOT NULL,
    outage_start DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
    outage_end DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `outages`;
-- +goose StatementEnd

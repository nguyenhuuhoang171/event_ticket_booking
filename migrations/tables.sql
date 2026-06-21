-- Table user
CREATE TABLE `event_ticket_booking`.`user` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `email`           VARCHAR(255)     NOT NULL,
    `name`            VARCHAR(255)     NOT NULL DEFAULT '',
    `hashed_password` VARCHAR(255)     NOT NULL,
    `role`            VARCHAR(50)      NOT NULL DEFAULT 'user',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      BIGINT UNSIGNED  NULL,
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by`      BIGINT UNSIGNED  NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_email` (`email`)
);

-- Table event
CREATE TABLE `event_ticket_booking`.`event` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `name`          VARCHAR(255)     NOT NULL,
    `description`   TEXT             NULL,
    `date_time`     DATETIME         NOT NULL,
    `total_tickets` BIGINT UNSIGNED  NOT NULL DEFAULT 0,
    `ticket_price`  BIGINT UNSIGNED  NOT NULL DEFAULT 0,
    `status`        TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 = active, 2 = deleted',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`    BIGINT UNSIGNED  NOT NULL,
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by`    BIGINT UNSIGNED  NOT NULL,
    `deleted_at`    DATETIME         NULL DEFAULT NULL,
    `deleted_by`    BIGINT UNSIGNED  NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_event_date_time` (`date_time`)
);
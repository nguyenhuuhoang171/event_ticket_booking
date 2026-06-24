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

-- Insert user
INSERT INTO event_ticket_booking.`user`
(id, email, name, hashed_password, `role`, created_at, created_by, updated_at, updated_by)
VALUES(4, 'user@mail.com', '', '$2a$10$J.IW56JZLb75Zb98cppne.JAIXVK0jfoziJjo783lizK9JWv4vt/e', 'user', '2026-06-24 07:16:15', 0, '2026-06-24 07:16:15', 0);

-- Table event
CREATE TABLE `event_ticket_booking`.`event` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `name`          VARCHAR(255)     NOT NULL,
    `description`   TEXT             NULL,
    `date_time`     DATETIME         NOT NULL,
    `total_tickets` BIGINT UNSIGNED  NOT NULL DEFAULT 0,
    `sold_tickets`  BIGINT UNSIGNED  NOT NULL DEFAULT 0,
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

-- Table booking
CREATE TABLE `event_ticket_booking`.`booking` (
    `id`         BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `event_id`   BIGINT UNSIGNED  NOT NULL,
    `user_id`    BIGINT UNSIGNED  NOT NULL,
    `quantity`   BIGINT UNSIGNED  NOT NULL DEFAULT 0,
    `status`     TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 = pending, 2 = confirmed, 3 = cancelled',
    `created_at` DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` BIGINT UNSIGNED  NOT NULL,
    `updated_at` DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by` BIGINT UNSIGNED  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_booking_event_id` (`event_id`),
    KEY `idx_booking_user_id` (`user_id`)
);

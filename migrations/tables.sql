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
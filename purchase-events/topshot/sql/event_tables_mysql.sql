
CREATE TABLE `moment_events` (
  `type` CHAR(1) NOT NULL,
  `MomentId` INT UNSIGNED NOT NULL,
  `BlockHeight` BIGINT UNSIGNED NOT NULL,
  `PlayId` SMALLINT UNSIGNED NOT NULL,
  `SerialNumber` INT UNSIGNED NOT NULL,
  `SetId` SMALLINT UNSIGNED NOT NULL,
  `Price` DECIMAL(11,2) NOT NULL,
  `BlockId` CHAR(64) NOT NULL,
  `SellerAddr` CHAR(16) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`type`, `MomentId`, `BlockHeight`));

CREATE TABLE `moment_events_collectors` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `CollectorId` SMALLINT UNSIGNED NOT NULL,
  `State` CHAR(2) NOT NULL,
  `UpdatesInInterval` INT NOT NULL,
  `BlockHeight` BIGINT UNSIGNED NOT NULL,
  `CreatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`));

CREATE INDEX `idx_moment_events_BlockHeight`  ON `moment_events` (BlockHeight) COMMENT 'Need to quickly get last blocks' ALGORITHM DEFAULT LOCK NONE;

CREATE INDEX `idx_moment_events_collectors_CreatedAt` ON `moment_events_collectors` (CreatedAt) COMMENT 'Need access to last collector state' ALGORITHM DEFAULT LOCK NONE;

CREATE TABLE users (
  id            binary(16)      NOT NULL DEFAULT (UUID_TO_BIN(UUID(),true)),
  email_address varchar(64)     NOT NULL,
  given_name    varchar(64),
  family_name   varchar(64),
  status        enum('INVITED', 'ACCEPTED', 'VERIFIED', 'DELETED', 'DISABLED', 'BANNED') NOT NULL DEFAULT 'INVITED',
  role          enum('ADMIN', 'CSM', 'EMPLOYEE', 'USER') NOT NULL DEFAULT 'USER',
  password      blob,
  picture_url   varchar(256),
  version       integer   NOT NULL DEFAULT 0,
  created       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted       timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE User (
  id            binary(16)      NOT NULL DEFAULT (UUID_TO_BIN(UUID(),true)),
  emailAddress  varchar(64)     NOT NULL,
  givenName     varchar(64),
  familyName    varchar(64),
  status        enum('INVITED', 'ACCEPTED', 'VERIFIED', 'DELETED', 'DISABLED', 'BANNED') NOT NULL DEFAULT 'INVITED',
  role          enum('ADMIN', 'CSM', 'EMPLOYEE', 'USER') NOT NULL DEFAULT 'USER',
  password      blob,
  picture_url   varchar(256),
  created       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modified      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted       timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
);

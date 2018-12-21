CREATE TABLE users (
  userid BIGINT(20) UNSIGNED,
  name VARCHAR(255),
  email VARCHAR(255),
  PRIMARY KEY (userid)
) ENGINE=InnoDB;

CREATE TABLE files (
  fileid BIGINT(20) UNSIGNED,
  userid BIGINT(20) UNSIGNED,
  content VARCHAR(255),
  PRIMARY KEY (fileid)
) ENGINE=InnoDB;
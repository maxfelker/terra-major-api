DROP TABLE IF EXISTS Characters;

CREATE TABLE Characters(
  ID UUID PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  bio VARCHAR(255) NOT NULL,
  age INT NOT NULL,
  strength INT NOT NULL,
  intelligence INT NOT NULL,
  endurance INT NOT NULL,
  agility INT NOT NULL,
);
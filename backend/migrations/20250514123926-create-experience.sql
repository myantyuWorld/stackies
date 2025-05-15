-- +migrate Up
CREATE TABLE experiences (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE experiences;

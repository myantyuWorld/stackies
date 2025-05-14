-- +migrate Up
INSERT INTO
  experiences (title)
VALUES
  ('Go');

INSERT INTO
  experiences (title)
VALUES
  ('Python');

INSERT INTO
  experiences (title)
VALUES
  ('Ruby');

-- +migrate Down
DELETE FROM
  experiences
WHERE
  title = 'Go';

DELETE FROM
  experiences
WHERE
  title = 'Python';

DELETE FROM
  experiences
WHERE
  title = 'Ruby';

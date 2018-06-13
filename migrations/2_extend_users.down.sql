DROP TABLE IF EXISTS users;
CREATE TABLE users (
  user_id integer unique,
  email   varchar(100)
);
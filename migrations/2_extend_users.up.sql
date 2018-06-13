DROP TABLE IF EXISTS users;
CREATE TABLE users (
  "user_id" SERIAL,
  "email"   varchar(100),
  "role" varchar(10),
   PRIMARY KEY("user_id")
);
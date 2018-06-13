CREATE TABLE token (
  "token_id" SERIAL,
  "token"   varchar(200),
  "user_id" integer REFERENCES "users" (user_id),
  "expiry_date" integer,
   PRIMARY KEY("token_id")
);

CREATE TABLE password (
  "password_id" SERIAL,
  "password"   varchar(100),
  "user_id" integer REFERENCES "users" (user_id),
   PRIMARY KEY("password_id")
);
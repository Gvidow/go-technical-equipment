ALTER TABLE users ADD COLUMN password varchar(80) NOT NULL DEFAULT 'default';
ALTER TABLE users ALTER COLUMN password DROP DEFAULT;

CREATE SEQUENCE email;
ALTER TABLE users ADD COLUMN email varchar(50) NOT NULL DEFAULT nextval('email')::TEXT || '@example.com' CONSTRAINT users_email_uniq UNIQUE;
ALTER TABLE users ALTER COLUMN email DROP DEFAULT;
DROP SEQUENCE email;

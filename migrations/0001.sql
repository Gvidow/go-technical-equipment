CREATE TABLE IF NOT EXISTS equipment (
    id serial PRIMARY KEY,
    title varchar(80) NOT NULL,
    description text,
    picture varchar(50),
    status varchar(8) NOT NULL,
    CONSTRAINT equipment_status CHECK (status = 'active' OR status = 'delete')
);

CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username varchar(30) UNIQUE NOT NULL,
    role varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS request (
    id serial PRIMARY KEY,
    equipment_id int NOT NULL,
    status varchar(10) NOT NULL,
    moderator int NOT NULL,
    creation_date timestamp,
    formation_date timestamp,
    completion_date timestamp,
    CONSTRAINT fk_request FOREIGN KEY (moderator) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT request_status CHECK (status = 
    ANY(ARRAY['entered', 'operation', 'completed', 'canceled', 'deleted']))
);


CREATE TABLE IF NOT EXISTS orders (
    id int PRIMARY KEY,
    equipment_id int REFERENCES equipment (id) ON DELETE CASCADE NOT NULL,
    request_id int REFERENCES request (id) ON DELETE CASCADE NOT NULL
);


CREATE TABLE IF NOT EXISTS equipment (
    id serial PRIMARY KEY,
    title varchar(80) NOT NULL,
    description text,
    picture varchar(50),
    status varchar(8) NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz,
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
    creator int NOT NULL,
    created_at timestamptz,
    formated_at timestamptz,
    completed_at timestamptz,
    CONSTRAINT request_fk_moderator FOREIGN KEY (moderator) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT request_fk_creator FOREIGN KEY (creator) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT request_status CHECK (status = 
    ANY(ARRAY['entered', 'operation', 'completed', 'canceled', 'deleted']))
);


CREATE TABLE IF NOT EXISTS orders (
    equipment_id int REFERENCES equipment (id) ON DELETE CASCADE NOT NULL,
    request_id int REFERENCES request (id) ON DELETE CASCADE NOT NULL,
    CONSTRAINT orders_pk PRIMARY KEY (equipment_id, request_id)
);

CREATE TABLE IF NOT EXISTS users(
    uid varchar(36) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    pass text NOT NULL,
    age integer,
    RegisteredAt timestamp NOT NULL default NOW()
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS books(
    bid varchar(36) NOT NULL PRIMARY KEY,
    lable text NOT NULL,
    author text NOT NULL,
    descriptons text NOT NULL,
    WritedAt timestamp NOT NULL
);

CREATE TABLE users(
    id integer primary key autoincrement,
    email varchar not null unique,
    encryptedPassword varchar not null
);
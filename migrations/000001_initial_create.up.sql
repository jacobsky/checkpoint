CREATE TABLE comments (
    id INTEGER NOT NULL PRIMARY KEY,
    postdate DATETIME,
    pinned BOOLEAN,
    poster VARCHAR(80),
    message TEXT
);

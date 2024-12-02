CREATE TABLE programs (
    id INTEGER PRIMARY KEY,
    program text NOT NULL,
    name text NOT NULL,
    author text NOT NULL,
    disabled int NOT NULL DEFAULT 0
);
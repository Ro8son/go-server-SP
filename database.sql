CREATE TABLE Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT,
    is_admin INTEGER NOT NULL DEFAULT 0
);

INSERT INTO Users (login, password) VALUES
('Tako', '$2a$12$owvRo/QyIoq1n4rfXx2D/uLA8i5cSpFNrjHY6KWx5ijU/oXe2c.1G'), -- password: Tako1234
('aa', '$2a$12$YRpJ.CFCxfv6i/3RMzzdTOl3T/EeYEL5nHKqVDcXTHFoQs3qdE9xG');   -- password: aa

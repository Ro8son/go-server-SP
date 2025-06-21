CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT,
  profile TEXT,
  is_admin INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  owner_id INTEGER NOT NULL,
  file_name TEXT NOT NULL,
  title TEXT,
  description TEXT,
  coordinates TEXT,
  FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE fileGuestShares (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  file_id INTEGER NOT NULL,
  url TEXT NOT NULL,       -- Unique shareable link token
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME,              -- Optional expiration
  max_uses INTEGER,                 -- Limit access attempts
  FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE fileTags (
  file_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  PRIMARY KEY (file_id, tag_id),
  FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE TABLE album (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT
);

INSERT INTO users (login, password) VALUES
('Tako', '$2a$12$owvRo/QyIoq1n4rfXx2D/uLA8i5cSpFNrjHY6KWx5ijU/oXe2c.1G'), -- password: Tako1234
('aa', '$2a$12$YRpJ.CFCxfv6i/3RMzzdTOl3T/EeYEL5nHKqVDcXTHFoQs3qdE9xG');   -- password: aa

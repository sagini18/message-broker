CREATE TABLE IF NOT EXISTS message (
    id INTEGER PRIMARY KEY,
    channel_name TEXT NOT NULL,
    content BLOB
);

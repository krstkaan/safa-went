CREATE TABLE IF NOT EXISTS books (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    barcode VARCHAR(255) UNIQUE,
    author_id BIGINT NOT NULL REFERENCES authors(id),
    publisher_id BIGINT NOT NULL REFERENCES publishers(id),
    language VARCHAR(100) NOT NULL,
    page_count INT NOT NULL,
    is_donation BOOLEAN NOT NULL DEFAULT FALSE,
    shelf_code VARCHAR(100) NOT NULL,
    fixture_no INT NOT NULL UNIQUE,
    level VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_books_author_id ON books (author_id);
CREATE INDEX IF NOT EXISTS idx_books_publisher_id ON books (publisher_id);
CREATE INDEX IF NOT EXISTS idx_books_level ON books (level);
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON books (deleted_at);

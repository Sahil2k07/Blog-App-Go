-- +goose Up

CREATE TABLE Blog (
    id CHAR(36) PRIMARY KEY,
    profileId CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags JSON,
    published BOOLEAN DEFAULT TRUE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (profileId) REFERENCES Profile(id) ON DELETE CASCADE
);

CREATE INDEX idx_blog_profileId ON Blog(profileId);
CREATE INDEX idx_blog_published ON Blog(published);
CREATE INDEX idx_blog_updatedAt ON Blog(updatedAt);

CREATE TABLE Otp (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    otp VARCHAR(6) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_otp_email ON Otp(email);

-- +goose Down

DROP TABLE IF EXISTS Blog;

DROP TABLE IF EXISTS Otp;
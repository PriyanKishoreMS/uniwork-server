CREATE TABLE IF NOT EXISTS colleges (
    id BIGINT AUTO_INCREMENT PRIMARY KEY UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    domain VARCHAR(255) NULL UNIQUE,
    version INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY UNIQUE,
    college_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    profile_pic TEXT,
    dept VARCHAR(255),
    review DECIMAL(2,1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    FOREIGN KEY (college_id) REFERENCES colleges(id),
    INDEX (college_id)
);

CREATE TABLE IF NOT EXISTS services (
    id BIGINT AUTO_INCREMENT PRIMARY KEY UNIQUE,
    user_id BIGINT NOT NULL,
    college_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255) NOT NULL,
    price BIGINT,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expiry DATETIME,
    images TEXT,
    version INT NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (college_id) REFERENCES colleges(id),
    INDEX (user_id),
    INDEX (college_id)
);

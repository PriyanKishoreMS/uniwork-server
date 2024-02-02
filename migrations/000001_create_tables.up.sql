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
    mobile VARCHAR(15) NOT NULL UNIQUE,
    profile_pic VARCHAR(512) NOT NULL DEFAULT "https://upload.wikimedia.org/wikipedia/commons/thumb/b/bc/Unknown_person.jpg/434px-Unknown_person.jpg",
    dept VARCHAR(255) NOT NULL,
    services_completed INT NOT NULL DEFAULT 0, 
    earned BIGINT NOT NULL DEFAULT 0,
    review DECIMAL(2,1) NOT NULL DEFAULT 5,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    FOREIGN KEY (college_id) REFERENCES colleges(id),
    INDEX (college_id),
    INDEX (name)
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

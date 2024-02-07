CREATE TABLE IF NOT EXISTS tasks (
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
    CONSTRAINT task_user_fk FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT task_clg_fk FOREIGN KEY (college_id) REFERENCES colleges(id),
    INDEX (user_id),
    INDEX (college_id)
);




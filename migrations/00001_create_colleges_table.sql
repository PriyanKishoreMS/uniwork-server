-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS colleges (
    id BIGSERIAL PRIMARY KEY UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    domain VARCHAR(255) NOT NULL UNIQUE,
    version INT NOT NULL DEFAULT 1
);

INSERT INTO colleges (name, domain) VALUES ('Hindustan University', 'hitsuniv@student.ac.in'), ('SRM University', 'srmuniv@student.ac.in'),
('VIT University', 'vituniv@student.ac.in'),
('Sathyabama University', 'sathyabamauniv@student.ac.in'),
('KCG College of Technology', 'kcgtech@student.ac.in'),
('Srinivasan Medical College and Hospital', 'smch@student.ac.in');
-- +goose StatementEnd





-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS colleges;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) UNIQUE,
    display_name VARCHAR(50) NOT NULL,
    dio TEXT,
    location VARCHAR(255),
    website VARCHAR(255),
    birth_date DATE,
    profile_image VARCHAR(255),
    header_image VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    activation_token VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);
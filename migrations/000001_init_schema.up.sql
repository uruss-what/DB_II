CREATE TYPE user_role AS ENUM ('superuser', 'admin', 'editor', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, password_hash, role)
VALUES ('admin', '$2a$10$ZY6Ue5YGjH5YFOhJPXYz8OVB0P6pGRQB8bPrVwJXJ5ZYyZH1X5Zy.', 'superuser'); 
CREATE TABLE images (
                        id SERIAL PRIMARY KEY,
                        filename VARCHAR(255) NOT NULL,
                        format VARCHAR(50) NOT NULL,
                        size INT NOT NULL,
                        upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        filepath TEXT NOT NULL
);
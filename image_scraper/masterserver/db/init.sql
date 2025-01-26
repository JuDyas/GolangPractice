CREATE TABLE images (
                        id SERIAL PRIMARY KEY,              -- Уникальный идентификатор изображения
                        filename VARCHAR(255) NOT NULL,     -- Название файла
                        format VARCHAR(50) NOT NULL,        -- Формат изображения (например, jpg, png)
                        size INT NOT NULL,                  -- Размер файла в байтах
                        upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата загрузки
                        filepath TEXT NOT NULL              -- Путь к файлу на сервере
);

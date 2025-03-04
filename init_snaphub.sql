-- Создание базы данных snabhub, если она не существует
CREATE DATABASE IF NOT EXISTS snaphub;
CREATE USER IF NOT EXISTS 'kranid'@'%' IDENTIFIED BY 'your_password';

-- Назначаем права на SELECT, INSERT и UPDATE для базы данных snaphub
GRANT SELECT, INSERT, UPDATE ON snaphub.* TO 'kranid'@'%';

-- Применяем изменения
FLUSH PRIVILEGES;
-- Использование созданной базы данных
USE snaphub;

-- Создание таблицы с нужными столбцами
CREATE TABLE IF NOT EXISTS snap_info (
    id INT AUTO_INCREMENT PRIMARY KEY,        -- Уникальный автоинкрементный столбец
    name VARCHAR(255) NOT NULL,               -- Строковый столбец, индексируемый
    JBId VARCHAR(255),                         -- Строковый столбец
    packagename VARCHAR(255),                  -- Строковый столбец
    INDEX idx_name (name)                     -- Индекс по столбцу name
);

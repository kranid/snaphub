-- Создание базы данных snabhub, если она не существует
CREATE DATABASE IF NOT EXISTS snaphub;
CREATE USER IF NOT EXISTS 'kranid'@'%' IDENTIFIED BY 'your_password';

-- Назначаем права на SELECT, INSERT, UPDATE и INSERT для базы данных snaphub
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

-- Создание таблицы для хранения информации о снапшотах
CREATE TABLE IF NOT EXISTS snapshots (
    id INT AUTO_INCREMENT PRIMARY KEY,
    package_name VARCHAR(255) NOT NULL,
    activity_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Новая таблица для связывания снапшотов с данными из JSONBin через snap_info
CREATE TABLE IF NOT EXISTS snapshot_json_links (
    id INT AUTO_INCREMENT PRIMARY KEY,
    snapshot_id INT NOT NULL,
    snap_info_id INT NOT NULL,
    data_type VARCHAR(50) NOT NULL, -- 'original', 'expected', 'technical', 'human'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (snapshot_id) REFERENCES snapshots(id) ON DELETE CASCADE,
    FOREIGN KEY (snap_info_id) REFERENCES snap_info(id) ON DELETE CASCADE
);
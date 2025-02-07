CREATE TABLE IF NOT EXISTS seed (
    seed_id INT AUTO_INCREMENT PRIMARY KEY,
    description TEXT NOT NULL,
    img TEXT,
    variety_name VARCHAR(100) NOT NULL,
    vegetable VARCHAR(100) NOT NULL
);


CREATE TABLE IF NOT EXISTS orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    sender_user_id INT NOT NULL,
    reciver_user_id INT NOT NULL,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    state ENUM('In attesa', 'In preparazione', 'In spedizione', 'Arrivato') DEFAULT 'In attesa',
    FOREIGN KEY (sender_user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (reciver_user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

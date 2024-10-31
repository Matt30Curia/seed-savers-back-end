CREATE TABLE IF NOT EXISTS order_detail (
    detail_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    seed_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (seed_id) REFERENCES seed(seed_id) ON DELETE CASCADE
);
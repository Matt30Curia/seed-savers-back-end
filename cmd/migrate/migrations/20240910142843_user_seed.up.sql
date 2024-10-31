CREATE TABLE IF NOT EXISTS users_seed (
    user_id INT,
    seed_id INT,
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, seed_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (seed_id) REFERENCES seed(seed_id) ON DELETE CASCADE
);
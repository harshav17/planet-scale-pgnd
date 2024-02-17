CREATE TABLE IF NOT EXISTS item_splits_nu (
    item_split_id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT NOT NULL,
    user_id VARCHAR(255),
    initials VARCHAR(5),
    amount DECIMAL(19,4) NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (item_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);
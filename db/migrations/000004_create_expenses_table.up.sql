CREATE TABLE IF NOT EXISTS expenses (
    expense_id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT,
    paid_by VARCHAR(255),
    amount DECIMAL(10, 2),
    description VARCHAR(100),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(group_id),
    FOREIGN KEY (paid_by) REFERENCES users(user_id)
    -- Add additional details if needed like receipt image path.
);

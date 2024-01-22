CREATE TABLE IF NOT EXISTS settlements (
    settlement_id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT NOT NULL, -- References the group_id from the expense_groups table
    paid_by VARCHAR(255) NOT NULL,
    paid_to VARCHAR(255) NOT NULL,
    amount DECIMAL(19,4),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES expense_groups(group_id),
    FOREIGN KEY (paid_by) REFERENCES users(user_id),
    FOREIGN KEY (paid_to) REFERENCES users(user_id)
    -- Additional details or notes can be added if necessary.
);

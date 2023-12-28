CREATE TABLE IF NOT EXISTS settlements (
    settlement_id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT, -- References the group_id from the groups table
    paid_by VARCHAR(255),
    paid_to VARCHAR(255),
    amount DECIMAL(10, 2),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(group_id),
    FOREIGN KEY (paid_by) REFERENCES users(auth0_id),
    FOREIGN KEY (paid_to) REFERENCES users(auth0_id)
    -- Additional details or notes can be added if necessary.
);

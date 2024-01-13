CREATE TABLE IF NOT EXISTS expenses (
    expense_id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT,
    split_type_id INT NOT NULL,
    paid_by VARCHAR(255),
    amount DECIMAL(19,4),
    description VARCHAR(100),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),  -- References the auth0_id from the users table
    updated_by VARCHAR(255),  -- References the auth0_id from the users table
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES expense_groups(group_id),
    FOREIGN KEY (paid_by) REFERENCES users(user_id),
    FOREIGN KEY (split_type_id) REFERENCES split_types(split_type_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id),
    FOREIGN KEY (updated_by) REFERENCES users(user_id),
    CONSTRAINT expense_amount CHECK (amount > 0)
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id INT,
    user_id VARCHAR(255),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES expense_groups(group_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
    -- You might want to add more information about the membership, like status or role.
);

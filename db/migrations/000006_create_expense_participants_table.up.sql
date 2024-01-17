CREATE TABLE IF NOT EXISTS expense_participants (
    expense_id INT,
    user_id VARCHAR(255),
    amount_owed DECIMAL(19,4),  -- The absolute amount the user owes for this expense
    share_percentage DECIMAL(5, 2),  -- The percentage of the total expense this user is responsible for
    note VARCHAR(255),  -- Optional field for any notes related to the split (e.g., reasons for uneven split)
    PRIMARY KEY (expense_id, user_id),
    FOREIGN KEY (expense_id) REFERENCES expenses(expense_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

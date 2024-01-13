CREATE TABLE IF NOT EXISTS items (
    item_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    expense_id INT NOT NULL,
    FOREIGN KEY (expense_id) REFERENCES expenses (expense_id)
);
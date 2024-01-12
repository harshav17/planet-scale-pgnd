BEGIN;

CREATE TABLE IF NOT EXISTS split_types (
    split_type_id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO split_types (type_name, description) VALUES ('Equal', 'Divides the bill equally among all participants');
INSERT INTO split_types (type_name, description) VALUES ('Unequal', 'Allows each participant to specify how much they are paying');
INSERT INTO split_types (type_name, description) VALUES ('ItemBased', 'Allows each participant to specify which items they are paying for');
INSERT INTO split_types (type_name, description) VALUES ('ShareBased', 'Splits the bill based on the share of each participant');
INSERT INTO split_types (type_name, description) VALUES ('PercentageBased', 'Splits the bill based on the percentage of each participant');

COMMIT;
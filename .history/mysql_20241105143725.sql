CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT,
    name string,
    created_at datetime,
    updated_at datetime,
)
CREATE TABLE IF NOT EXISTS users (
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    customer_id INT(11),
    name varchar(64),
    created_at datetime,
    updated_at datetime,
)
CREATE TABLE IF NOT EXISTS users (
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    customer_id INT(11) NOT NULL,
    username VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
);
CREATE TABLE IF NOT EXISTS customer_admins(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    cstomer_id INT(11) NOT NULL,
    username  VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
)


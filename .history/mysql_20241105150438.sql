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

CREATE TABLE IF NOT EXISTS customer_admin_chat_settings(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    admin_id INT(11) NOT NULL,
    background VARCHAR(255) NOT NULL,
    is_auto_accept TINYINT(4) NOT NULL DEFAULT 0,
    welcome_content VARCHAR(512) DEFAULT "",
    offline_content VARCHAR(512) DEFAULT "",
    name VARCHAR(32) DEFAULT "",
    last_online DATETIME DEFAULT NULL,
    avatar VARCHAR(11) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
)
CREATE TABLE IF NOT EXISTS customer_admin_auto_messages(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,
    content VARCHAR(512) NOT NULL,
    customer_id INT(11) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
)


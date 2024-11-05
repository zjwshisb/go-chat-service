CREATE TABLE IF NOT EXISTS users (
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    customer_id INT(11) NOT NULL,
    username VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX (username)
);
CREATE TABLE IF NOT EXISTS customer_admins(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    cstomer_id INT(11) NOT NULL,
    username  VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    index (username)
);

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
    index (admin_id)
);
CREATE TABLE IF NOT EXISTS customer_chat_auto_messages(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,
    content VARCHAR(512) NOT NULL,
    customer_id INT(11) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
);
CREATE TABLE IF NOT EXISTS customer_chat_auto_rule_scenes(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    rule_id INT(11) NOT NULL,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    index (rule_id)
);
CREATE TABLE IF NOT EXISTS customer_chat_auto_rules(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    customer_id INT(11) INDEX, 
    name VARCHAR(64) NOT NULL,
    match VARCHAR(64) NOT NULL,
    match_type VARCHAR(64) NOT NULL,
    reply_type VARCHAR(64) NOT NULL,
    message_id  INT(11) NOT NULL,
    is_system TINYINT(11) NOT NULL,
    sort INT(11) NOT NULL,
    is_open TINYINT(11) NOT NULL,
    count BIGINT(29) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
);
CREATE TABLE IF NOT EXISTS customer_chat_messages(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    user_id INT(11) INDEX, 
    admin_id INT(11) NOT NULL default 0,
    customer_id  INT(11) NOT NULL,
    type VARCHAR(16) NOT NULL,
    content VARCHAR(512) NOT NULL,
    received_at BIGINT(20) DEFAULT 0,
    send_at BIGINT(20) DEFAULT 0,
    source TINYINT(4) NOT NULL,
    session_id INT(11) DEFAULT 0,
    req_id VARCHAR(64) NOT NULL,
    read_at BIGINT(20) DEFAULT 0,
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX (user_id),
    INDEX (admin_id)
);
CREATE TABLE IF NOT EXISTS customer_chat_sessions(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    user_id INT(11) INDEX, 
    queried_at BIGINT(20) NOT NULL default 0,
    accepted_at BIGINT(20) NOT NULL default 0,
    canceled_at BIGINT(20) NOT NULL default 0,
    broked_at BIGINT(20) NOT NULL default 0,
    customer_id INT(11) NOT NULL,
    admin_id INT(11) NOT NULL,
    type TINYINT(4) DEFAULT 0,
    rate SMALLINT(6) DEFAULT NULL,
    INDEX (user_id),
    INDEX (admin_id)
);
CREATE TABLE IF NOT EXISTS customer_chat_settings(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255), 
    title  VARCHAR(255), 
    customer_id INT(11) NOT NULL,
    value VARCHAR(512),
    options VARCHAR(512),
    type VARCHAR(32),
    description VARCHAR(64),
    created_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
);
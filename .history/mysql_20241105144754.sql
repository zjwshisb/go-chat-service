CREATE TABLE IF NOT EXISTS users (
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    customer_id INT(11),
    username varchar(64),
    created_at datetime,
    updated_at datetime,
);
CREATE TABLE IF NOT EXISTS customer_admins(
    id INT(11) PRIMARY KEY AUTO_INCREMENT,
    cstomer_id INT(11),
    username  varchar(64),
    password varchar(255),
    created_at datetime,
    update_at datetime,
    deleted_at datetime default null,
)


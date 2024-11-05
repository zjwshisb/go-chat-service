-- dev.customer_admins definition

CREATE TABLE `customer_admins` (
                                   `id` int NOT NULL AUTO_INCREMENT,
                                   `customer_id` int NOT NULL,
                                   `username` varchar(64) NOT NULL,
                                   `password` varchar(255) NOT NULL,
                                   `created_at` datetime NOT NULL,
                                   `update_at` datetime NOT NULL,
                                   `deleted_at` datetime DEFAULT NULL,
                                   PRIMARY KEY (`id`) USING BTREE,
                                   KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- dev.users definition

CREATE TABLE `users` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `customer_id` int NOT NULL,
                         `username` varchar(64) NOT NULL,
                         `created_at` datetime NOT NULL,
                         `updated_at` datetime NOT NULL,
                         `deleted_at` datetime DEFAULT NULL,
                         PRIMARY KEY (`id`) USING BTREE,
                         KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
#
# CREATE TABLE
#     IF NOT EXISTS customer_admin_chat_settings (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         admin_id INT (11) NOT NULL,
#         background VARCHAR(255) NOT NULL,
#         is_auto_accept TINYINT (4) NOT NULL DEFAULT 0,
#         welcome_content VARCHAR(512) DEFAULT '',
#         offline_content VARCHAR(512) DEFAULT '',
#         name VARCHAR(32) DEFAULT '',
#         last_online DATETIME DEFAULT NULL,
#         avatar VARCHAR(11) NOT NULL,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (admin_id)
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_auto_messages (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         `name` VARCHAR(64) NOT NULL,
#         `type` VARCHAR(32) NOT NULL,
#         content VARCHAR(512) NOT NULL,
#         customer_id INT (11) NOT NULL,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_auto_rule_scenes (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         `name` VARCHAR(64) NOT NULL,
#         rule_id INT (11) NOT NULL,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (rule_id)
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_auto_rules (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         customer_id INT (11) ,
#         `name` VARCHAR(64) NOT NULL,
#         `match` VARCHAR(64) NOT NULL,
#         match_type VARCHAR(64) NOT NULL,
#         reply_type VARCHAR(64) NOT NULL,
#         message_id INT (11) NOT NULL,
#         is_system TINYINT (11) NOT NULL,
#         sort INT (11) NOT NULL,
#         is_open TINYINT (11) NOT NULL,
#         count BIGINT (29) NOT NULL DEFAULT 0,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         index (customer_id) using BTREE
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_messages (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         user_id INT (11),
#         admin_id INT (11) NOT NULL default 0,
#         customer_id INT (11) NOT NULL,
#         `type` VARCHAR(16) NOT NULL,
#         content VARCHAR(512) NOT NULL,
#         received_at DATETIME NOT NULL DEFAULT NULL,
#         send_at DATETIME NOT NULL DEFAULT NULL,
#         source TINYINT (4) NOT NULL,
#         session_id INT (11) DEFAULT 0,
#         req_id VARCHAR(64) NOT NULL,
#         read_at DATETIME NOT NULL DEFAULT NULL,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (user_id) USING btree ,
#         INDEX (admin_id) USING btree
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_sessions (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         user_id INT (11),
#         queried_at DATETIME NOT NULL,
#         accepted_at DATETIME NOT NULL DEFAULT NULL,
#         canceled_at DATETIME NOT NULL DEFAULT NULL,
#         broken_at DATETIME NOT NULL DEFAULT NULL,
#         customer_id INT (11) NOT NULL,
#         admin_id INT (11) NOT NULL,
#         `type` TINYINT (4) DEFAULT 0,
#         rate SMALLINT (6) DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (user_id) USING BTREE ,
#         INDEX (admin_id) USING  BTREE
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_settings (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         name VARCHAR(255),
#         title VARCHAR(255),
#         customer_id INT (11) NOT NULL,
#         `value` VARCHAR(512),
#         options VARCHAR(512),
#         `type` VARCHAR(32),
#         `description` VARCHAR(64),
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (customer_id)
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer_chat_transfers (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         user_id INT (11) NOT NULL,
#         from_session_id INT (11) NOT NULL,
#         to_session_id INT (11) NOT NULL DEFAULT 0,
#         from_admin_id INT (11) NOT NULL,
#         to_admin_id INT (11) NOT NULL,
#         customer_id INT (11) NOT NULL,
#         remark VARCHAR(512) DEFAULT '',
#         accepted_at DATETIME DEFAULT NULL,
#         canceled_at DATETIME DEFAULT NULL,
#         created_at DATETIME NOT NULL,
#         update_at DATETIME NOT NULL,
#         PRIMARY KEY (`id`) USING BTREE,
#         INDEX (customer_id)
#     );
#
# CREATE TABLE
#     IF NOT EXISTS customer (
#         id INT (11) NOT NULL AUTO_INCREMENT,
#         name VARCHAR(64) NOT NULL,
#         created_at DATETIME NOT NULL,
#         updated_at DATETIME NOT NULL,
#         deleted_at DATETIME DEFAULT NULL,
#         PRIMARY KEY (`id`) USING BTREE
#     );
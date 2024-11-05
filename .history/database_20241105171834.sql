-- dev.customer definition
CREATE TABLE
    `customers` (
        `id` int NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_admin_chat_settings definition
CREATE TABLE
    `customer_admin_chat_settings` (
        `id` int NOT NULL AUTO_INCREMENT,
        `admin_id` int NOT NULL,
        `background` varchar(255) NOT NULL,
        `is_auto_accept` tinyint NOT NULL DEFAULT '0',
        `welcome_content` varchar(512) DEFAULT '',
        `offline_content` varchar(512) DEFAULT '',
        `name` varchar(32) DEFAULT '',
        `last_online` datetime DEFAULT NULL,
        `avatar` varchar(11) NOT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `admin_id` (`admin_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_admins definition
CREATE TABLE
    `customer_admins` (
        `id` int NOT NULL AUTO_INCREMENT,
        `customer_id` int NOT NULL,
        `username` varchar(64) NOT NULL,
        `password` varchar(255) NOT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `username` (`username`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_auto_messages definition
CREATE TABLE
    `customer_chat_auto_messages` (
        `id` int NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `type` varchar(32) NOT NULL,
        `content` varchar(512) NOT NULL,
        `customer_id` int NOT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_auto_rules definition
CREATE TABLE
    `customer_chat_auto_rules` (
        `id` int NOT NULL AUTO_INCREMENT,
        `customer_id` int DEFAULT NULL,
        `name` varchar(64) NOT NULL,
        `match` varchar(64) NOT NULL,
        `match_type` varchar(64) NOT NULL,
        `reply_type` varchar(64) NOT NULL,
        `message_id` int NOT NULL,
        `is_system` tinyint NOT NULL,
        `sort` int NOT NULL,
        `is_open` tinyint NOT NULL,
        `count` bigint NOT NULL DEFAULT '0',
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_messages definition
CREATE TABLE
    `customer_chat_messages` (
        `id` int NOT NULL AUTO_INCREMENT,
        `user_id` int DEFAULT NULL,
        `admin_id` int NOT NULL DEFAULT '0',
        `customer_id` int NOT NULL,
        `type` varchar(16) NOT NULL,
        `content` varchar(512) NOT NULL,
        `received_at` datetime DEFAULT NULL,
        `send_at` datetime DEFAULT NULL,
        `source` tinyint NOT NULL,
        `session_id` int DEFAULT '0',
        `req_id` varchar(64) NOT NULL,
        `read_at` datetime DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `user_id` (`user_id`) USING BTREE,
        KEY `admin_id` (`admin_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_sessions definition
CREATE TABLE
    `customer_chat_sessions` (
        `id` int NOT NULL AUTO_INCREMENT,
        `user_id` int DEFAULT NULL,
        `queried_at` datetime NOT NULL,
        `accepted_at` datetime DEFAULT NULL,
        `canceled_at` datetime DEFAULT NULL,
        `broken_at` datetime DEFAULT NULL,
        `customer_id` int NOT NULL,
        `admin_id` int NOT NULL,
        `type` tinyint DEFAULT '0',
        `rate` smallint DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `user_id` (`user_id`) USING BTREE,
        KEY `admin_id` (`admin_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_settings definition
CREATE TABLE
    `customer_chat_settings` (
        `id` int NOT NULL AUTO_INCREMENT,
        `name` varchar(255) DEFAULT NULL,
        `title` varchar(255) DEFAULT NULL,
        `customer_id` int NOT NULL,
        `value` varchar(512) DEFAULT NULL,
        `options` varchar(512) DEFAULT NULL,
        `type` varchar(32) DEFAULT NULL,
        `description` varchar(64) DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.customer_chat_transfers definition
CREATE TABLE
    `customer_chat_transfers` (
        `id` int NOT NULL AUTO_INCREMENT,
        `user_id` int NOT NULL,
        `from_session_id` int NOT NULL DEFAULT '0',
        `to_session_id` int NOT NULL DEFAULT '0',
        `from_admin_id` int NOT NULL DEFAULT '0',
        `to_admin_id` int NOT NULL DEFAULT '0',
        `customer_id` int NOT NULL,
        `remark` varchar(512) DEFAULT '',
        `accepted_at` datetime DEFAULT NULL,
        `canceled_at` datetime DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `update_at` datetime NOT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- dev.users definition
CREATE TABLE
    `users` (
        `id` int NOT NULL AUTO_INCREMENT,
        `customer_id` int NOT NULL,
        `username` varchar(64) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `username` (`username`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_auto_rule_scenes` (
        `id` int NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `rule_id` int NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
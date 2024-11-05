CREATE TABLE
    `customers` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_admin_chat_settings` (
        `id`  int  UNSIGNED NOT NULL AUTO_INCREMENT,
        `admin_id`  int UNSIGNED NOT NULL,
        `background` varchar(255) NOT NULL,
        `is_auto_accept` tinyint UNSIGNED NOT NULL DEFAULT '0',
        `welcome_content` varchar(512) DEFAULT '',
        `offline_content` varchar(512) DEFAULT '',
        `name` varchar(32) DEFAULT '',
        `last_online` datetime DEFAULT NULL,
        `avatar` varchar(11) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `admin_id` (`admin_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_admins` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `customer_id` int UNSIGNED NOT NULL,
        `username` varchar(64) NOT NULL,
        `password` varchar(255) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `username` (`username`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_auto_messages` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `type` varchar(32) NOT NULL,
        `content` varchar(512) NOT NULL,
        `customer_id` int UNSIGNED NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_auto_rules` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `customer_id` int UNSIGNED DEFAULT NULL,
        `name` varchar(64) NOT NULL,
        `match` varchar(64) NOT NULL,
        `match_type` varchar(64) NOT NULL,
        `reply_type` varchar(64) NOT NULL,
        `message_id` int UNSIGNED NOT NULL,
        `is_system` tinyint UNSIGNED NOT NULL,
        `sort` int UNSIGNED NOT NULL,
        `is_open` tinyint NOT NULL,
        `count` bigint NOT NULL DEFAULT '0',
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_messages` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `user_id` int UNSIGNED DEFAULT NULL,
        `admin_id` int UNSIGNED NOT NULL DEFAULT '0',
        `customer_id` int UNSIGNED NOT NULL,
        `type` varchar(16) NOT NULL,
        `content` varchar(512) NOT NULL,
        `received_at` datetime DEFAULT NULL,
        `send_at` datetime DEFAULT NULL,
        `source` tinyint UNSIGNED NOT NULL,
        `session_id` int UNSIGNED DEFAULT '0',
        `req_id` varchar(64) NOT NULL,
        `read_at` datetime DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `user_id` (`user_id`) USING BTREE,
        KEY `admin_id` (`admin_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_sessions` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `user_id` int UNSIGNED DEFAULT NULL,
        `queried_at` datetime NOT NULL,
        `accepted_at` datetime DEFAULT NULL,
        `canceled_at` datetime DEFAULT NULL,
        `broken_at` datetime DEFAULT NULL,
        `customer_id` int UNSIGNED NOT NULL,
        `admin_id` int UNSIGNED NOT NULL,
        `type` tinyint UNSIGNED DEFAULT '0',
        `rate` smallint UNSIGNED DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `user_id` (`user_id`) USING BTREE,
        KEY `admin_id` (`admin_id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_settings` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` varchar(255) DEFAULT NULL,
        `title` varchar(255) DEFAULT NULL,
        `customer_id` int UNSIGNED NOT NULL,
        `value` varchar(512) DEFAULT NULL,
        `options` varchar(512) DEFAULT NULL,
        `type` varchar(32) DEFAULT NULL,
        `description` varchar(64) DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_transfers` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `user_id` int UNSIGNED NOT NULL,
        `from_session_id` int UNSIGNED NOT NULL DEFAULT '0',
        `to_session_id` int UNSIGNED NOT NULL DEFAULT '0',
        `from_admin_id` int UNSIGNED NOT NULL DEFAULT '0',
        `to_admin_id` int UNSIGNED NOT NULL DEFAULT '0',
        `customer_id` int UNSIGNED NOT NULL,
        `remark` varchar(512) DEFAULT '',
        `accepted_at` datetime DEFAULT NULL,
        `canceled_at` datetime DEFAULT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `customer_id` (`customer_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `users` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `customer_id` int UNSIGNED NOT NULL,
        `username` varchar(64) NOT NULL,
        `created_at` datetime NOT NULL,
        `updated_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE,
        KEY `username` (`username`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE
    `customer_chat_auto_rule_scenes` (
        `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` varchar(64) NOT NULL,
        `rule_id` int UNSIGNED NOT NULL,
        `updated_at` datetime NOT NULL,
        `created_at` datetime NOT NULL,
        `deleted_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
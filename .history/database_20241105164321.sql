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
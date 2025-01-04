/*
 Navicat Premium Data Transfer

 Source Server         : aliyun
 Source Server Type    : MySQL
 Source Server Version : 50742
 Source Host           : 120.77.242.145:3306
 Source Schema         : chat

 Target Server Type    : MySQL
 Target Server Version : 50742
 File Encoding         : 65001

 Date: 04/01/2025 11:03:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for customer_admin_chat_settings
-- ----------------------------
DROP TABLE IF EXISTS `customer_admin_chat_settings`;
CREATE TABLE `customer_admin_chat_settings` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `admin_id` int(10) unsigned NOT NULL,
  `background` int(10) unsigned NOT NULL,
  `is_auto_accept` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `welcome_content` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `offline_content` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `name` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `last_online` datetime DEFAULT NULL,
  `avatar` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `customer_admin_chat_settings_admin_id_IDX` (`admin_id`) USING BTREE,
  KEY `admin_id` (`admin_id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_admins
-- ----------------------------
DROP TABLE IF EXISTS `customer_admins`;
CREATE TABLE `customer_admins` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` int(10) unsigned NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_auto_messages
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_auto_messages`;
CREATE TABLE `customer_chat_auto_messages` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_auto_rule_scenes
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_auto_rule_scenes`;
CREATE TABLE `customer_chat_auto_rule_scenes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `rule_id` int(10) unsigned NOT NULL,
  `updated_at` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_auto_rules
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_auto_rules`;
CREATE TABLE `customer_chat_auto_rules` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` int(10) unsigned DEFAULT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `match` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `match_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `reply_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `message_id` int(10) unsigned NOT NULL,
  `is_system` tinyint(3) unsigned NOT NULL,
  `sort` int(10) unsigned NOT NULL,
  `is_open` tinyint(4) NOT NULL,
  `scenes` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '""',
  `count` bigint(20) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `customer_id` (`customer_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_files
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_files`;
CREATE TABLE `customer_chat_files` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` int(10) unsigned NOT NULL COMMENT '客户Id',
  `disk` varchar(16) COLLATE utf8mb4_bin NOT NULL DEFAULT 'local' COMMENT '存储引擎',
  `path` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '路径',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `from_model` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '来源模型',
  `from_id` int(10) unsigned NOT NULL COMMENT '来源id',
  `type` varchar(16) COLLATE utf8mb4_bin NOT NULL DEFAULT 'image' COMMENT '文件类型',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `is_resource` tinyint(4) NOT NULL DEFAULT '0',
  `parent_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=106 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- Table structure for customer_chat_messages
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_messages`;
CREATE TABLE `customer_chat_messages` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned DEFAULT NULL,
  `admin_id` int(10) unsigned NOT NULL DEFAULT '0',
  `customer_id` int(10) unsigned NOT NULL,
  `type` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL,
  `received_at` datetime DEFAULT NULL,
  `send_at` datetime DEFAULT NULL,
  `source` tinyint(3) unsigned NOT NULL,
  `session_id` int(10) unsigned DEFAULT '0',
  `req_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `read_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_id` (`user_id`) USING BTREE,
  KEY `admin_id` (`admin_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2175 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_sessions
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_sessions`;
CREATE TABLE `customer_chat_sessions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned DEFAULT NULL,
  `queried_at` datetime NOT NULL,
  `accepted_at` datetime DEFAULT NULL,
  `canceled_at` datetime DEFAULT NULL,
  `broken_at` datetime DEFAULT NULL,
  `customer_id` int(10) unsigned NOT NULL,
  `admin_id` int(10) unsigned NOT NULL,
  `type` tinyint(3) unsigned DEFAULT '0',
  `rate` smallint(5) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_id` (`user_id`) USING BTREE,
  KEY `admin_id` (`admin_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_settings
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_settings`;
CREATE TABLE `customer_chat_settings` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `customer_id` int(10) unsigned NOT NULL,
  `value` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `options` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `customer_id` (`customer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customer_chat_transfers
-- ----------------------------
DROP TABLE IF EXISTS `customer_chat_transfers`;
CREATE TABLE `customer_chat_transfers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `from_session_id` int(10) unsigned NOT NULL DEFAULT '0',
  `to_session_id` int(10) unsigned NOT NULL DEFAULT '0',
  `from_admin_id` int(10) unsigned NOT NULL DEFAULT '0',
  `to_admin_id` int(10) unsigned NOT NULL DEFAULT '0',
  `customer_id` int(10) unsigned NOT NULL,
  `remark` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `accepted_at` datetime DEFAULT NULL,
  `canceled_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `customer_id` (`customer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for customers
-- ----------------------------
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` int(10) unsigned NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

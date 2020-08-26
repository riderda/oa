/*
 Navicat Premium Data Transfer

 Source Server         : test
 Source Server Type    : MySQL
 Source Server Version : 50724
 Source Host           : localhost:3306
 Source Schema         : demo

 Target Server Type    : MySQL
 Target Server Version : 50724
 File Encoding         : 65001

 Date: 26/08/2020 22:18:04
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for userinfo
-- ----------------------------
DROP TABLE IF EXISTS `userinfo`;
CREATE TABLE `userinfo`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user id` int(20) NOT NULL,
  `full name` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `join time` date NOT NULL,
  `leave time` date NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user id`(`user id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of userinfo
-- ----------------------------
INSERT INTO `userinfo` VALUES (1, 1, '陈伟达', '2020-08-26', NULL);

SET FOREIGN_KEY_CHECKS = 1;

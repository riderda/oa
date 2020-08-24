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

 Date: 24/08/2020 23:46:32
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for blog
-- ----------------------------
DROP TABLE IF EXISTS `blog`;
CREATE TABLE `blog`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `keyword` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `title` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '因为有保存为草稿的功能，所以可以为空',
  `author` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `record` int(255) NOT NULL DEFAULT 0,
  `public status` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `public time` datetime(0) NULL DEFAULT NULL COMMENT '保存草稿的时候没有发布时间，确认发布后修改',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of blog
-- ----------------------------
INSERT INTO `blog` VALUES (1, 'go', '陈伟达', 'ss', 'ss', 0, 'ss', '2020-08-24 22:44:15');
INSERT INTO `blog` VALUES (2, 'javasript', '细粒丁', 'dd', 'dd', 0, 'dd', '2020-08-24 22:44:53');
INSERT INTO `blog` VALUES (3, 'java', '洪伟龙', 'ssd', 'dss', 0, 'ddd', '2020-08-24 22:45:31');

SET FOREIGN_KEY_CHECKS = 1;

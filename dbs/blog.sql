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

 Date: 27/08/2020 15:30:32
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
  `summary` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '摘要，内容的前100个字',
  `author` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `record` int(255) NOT NULL DEFAULT 0,
  `public status` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `public time` datetime(0) NULL DEFAULT NULL COMMENT '保存草稿的时候没有发布时间，确认发布后修改',
  `is show` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用于绝对是否展示',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of blog
-- ----------------------------
INSERT INTO `blog` VALUES (1, 'dsddasdaa', 'sdasdasddasd', 'ddasdsddasdasdasdasdasdasdadasdas', 'ddasdsddasdasdasdasdasdasdadasdas', '', 0, 'release', '2020-08-26 22:14:24', '');
INSERT INTO `blog` VALUES (2, 'dsddasdaa', 'sdasdasddasd', 'ddasdsddasdasdasdasdasdasdadasdas', 'ddasdsddasdasdasdasdasdasdadasdas', '', 0, '', '2020-08-26 22:14:24', '');
INSERT INTO `blog` VALUES (3, 'dsddasdaa', 'sdasdasddasd', 'ddasdsddasdasdasdasdasdasdadasdas', 'ddasdsddasdasdasdasdasdasdadasdas', '', 0, '', '2020-08-26 22:14:24', '');
INSERT INTO `blog` VALUES (4, 'dsddasdaa', 'sdasdasddasd', 'llllllllllllllllllllllllllllllkjkjkljklhkhjj', 'llllllllllllllllllllllllllllllkjkjkljklhkhjj', '', 0, 'release', '2020-08-26 22:16:18', NULL);

SET FOREIGN_KEY_CHECKS = 1;

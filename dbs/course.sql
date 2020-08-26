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

 Date: 26/08/2020 22:17:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for course
-- ----------------------------
DROP TABLE IF EXISTS `course`;
CREATE TABLE `course`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 26 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of course
-- ----------------------------
INSERT INTO `course` VALUES (1, 'javascript', 'https://wangdoc.com/es6/', 'ES6-网道', '全面介绍 ECMAScript 6 新引入的语法特性，对涉及的语法知识给予详细介绍，并给出大量简洁易懂的示例代码。');
INSERT INTO `course` VALUES (2, 'javascript', 'https://wangdoc.com/javascript/', 'javascript教程-网道', '本教程全面介绍 JavaScript 核心语法，覆盖了 ES5 和 DOM 规范的所有内容。内容上从最简单的讲起，循序渐进、由浅入深，力求清晰易懂。所有章节都带有大量的代码实例，便于理解和模仿，可以用到实际项目中，即学即用。');
INSERT INTO `course` VALUES (3, 'go', 'https://www.kancloud.cn/kancloud/web-application-with-golang/44105', 'Go Web 编程', '学习Go过程以及以前从事Web开发过程中的一些经验总结');
INSERT INTO `course` VALUES (4, 'go', 'https://www.w3cschool.cn/go/', 'Go 教程', 'Go语言是谷歌2009年发布的第二款开源编程语言,它专门针对多处理器系统应用程序的编程进行了优化，它是一种系统语言其非常有用和强大,其程序可以媲美C或C++代码的速度，而且更加安全、支持并行进程。Go支持面向对象，而且具有真正的闭包(closures)和反射 (reflection)等功能。Go可以在不损失应用程序性能的情况下降低代码的复杂性');
INSERT INTO `course` VALUES (5, 'javascript', 'https://enable-javascript.com/', 'How to enable JavaScript in your browser', 'Nowadays almost all web pages contain JavaScript, a scripting programming language that runs on visitor\'s web browser. It makes web pages functional for specific purposes and if disabled for some reason, the content or the functionality of the web page can be limited or unavailable. Here you can find instructions on how to enable (activate) JavaScript in five most commonly used browsers.');
INSERT INTO `course` VALUES (6, 'javascript', 'https://www.baidu.com', 'javascript', 'javascript');
INSERT INTO `course` VALUES (7, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (8, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (9, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (10, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (11, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (12, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (13, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (14, 'javascript', 'javascript', 'javascript', 'javascript');
INSERT INTO `course` VALUES (15, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (16, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (17, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (18, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (19, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (20, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (21, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (22, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (23, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (24, 'go', 'go', 'go', 'go');
INSERT INTO `course` VALUES (25, 'go', 'go', 'go', 'go');

SET FOREIGN_KEY_CHECKS = 1;

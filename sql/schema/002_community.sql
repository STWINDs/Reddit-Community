-- +goose up
CREATE TABLE IF NOT EXISTS `community` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `community_id` bigint(20) UNSIGNED NOT NULL,
  `community_name` varchar(255) NOT NULL COLLATE utf8mb4_general_ci COMMENT '社区名称',
  `introduction` VARCHAR(256) COLLATE utf8mb4_general_ci COMMENT '社区描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `community_id` (`community_id`) USING BTREE,
  UNIQUE KEY `community_name` (`community_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



-- +goose down
DROP TABLE IF EXISTS `community`;
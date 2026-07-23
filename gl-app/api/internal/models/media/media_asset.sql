CREATE TABLE `media_asset` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除标记',
  `user_id` bigint unsigned NOT NULL COMMENT '上传用户ID',
  `resource_type` varchar(16) NOT NULL COMMENT '资源类型: image/video',
  `bucket` varchar(128) NOT NULL COMMENT '对象存储桶',
  `object_key` varchar(512) NOT NULL COMMENT '对象存储路径',
  `origin_name` varchar(255) NOT NULL COMMENT '原始文件名',
  `content_type` varchar(128) NOT NULL DEFAULT '' COMMENT '文件MIME类型',
  `ext` varchar(32) NOT NULL DEFAULT '' COMMENT '文件扩展名',
  `file_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '对象存储中的文件大小',
  `status` varchar(32) NOT NULL DEFAULT 'uploading' COMMENT '上传状态: uploading/uploaded/failed',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_object_key` (`object_key`),
  KEY `idx_user_status` (`user_id`, `status`),
  KEY `idx_user_resource_type` (`user_id`, `resource_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='媒体素材表';

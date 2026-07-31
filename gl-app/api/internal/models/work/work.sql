CREATE TABLE `work` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '作品主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  `user_id` bigint unsigned NOT NULL COMMENT '作品作者用户ID',
  `type` varchar(16) NOT NULL COMMENT '作品类型：image图片作品、video视频作品',
  `title` varchar(100) NOT NULL COMMENT '作品主标题',
  `content` varchar(4000) NOT NULL DEFAULT '' COMMENT '作品正文内容',
  `visibility` varchar(32) NOT NULL COMMENT '可见范围：public/private/mutual/selected/excluded',
  `visibility_user_ids` json DEFAULT NULL COMMENT '指定可见或排除可见的用户ID列表',
  `collection_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属合集ID，0表示未加入合集',
  `original` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否声明原创：0否、1是',
  `cover_asset_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '封面媒体素材ID',
  `process_status` varchar(32) NOT NULL DEFAULT 'pending' COMMENT '媒体处理状态：pending/processing/succeeded/failed',
  `review_status` varchar(32) NOT NULL DEFAULT 'waiting_process' COMMENT '审核状态：waiting_process/pending_review/approved/rejected',
  `publish_status` varchar(32) NOT NULL DEFAULT 'pending' COMMENT '发布状态：pending/scheduled/published/deleted',
  `scheduled_at` datetime(3) DEFAULT NULL COMMENT '计划发布时间，空表示审核通过后立即发布',
  `published_at` datetime(3) DEFAULT NULL COMMENT '实际发布时间',
  `reviewed_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '审核员用户ID',
  `reviewed_at` datetime(3) DEFAULT NULL COMMENT '审核完成时间',
  `review_reason` varchar(1000) NOT NULL DEFAULT '' COMMENT '审核意见或拒绝原因',
  `idempotency_key` varchar(64) NOT NULL COMMENT '用户维度的发布幂等键',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_idempotency` (`user_id`,`idempotency_key`),
  KEY `idx_user_created` (`user_id`,`created_at`),
  KEY `idx_review_status` (`review_status`,`created_at`),
  KEY `idx_schedule` (`publish_status`,`scheduled_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='作品';

CREATE TABLE `collection` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '合集主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  `user_id` bigint unsigned NOT NULL COMMENT '合集创建者用户ID',
  `name` varchar(100) NOT NULL COMMENT '合集名称',
  PRIMARY KEY (`id`),
  KEY `idx_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='作品合集';

CREATE TABLE `work_asset` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '作品素材关系主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `work_id` bigint unsigned NOT NULL COMMENT '作品ID',
  `media_asset_id` bigint unsigned NOT NULL COMMENT '媒体素材ID',
  `role` varchar(16) NOT NULL COMMENT '素材角色：image图片、video视频、cover封面',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '素材展示顺序，从0开始',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_work_media` (`work_id`,`media_asset_id`),
  UNIQUE KEY `idx_media_asset` (`media_asset_id`),
  KEY `idx_work_role_sort` (`work_id`,`role`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='作品资源';

CREATE TABLE `topic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '话题主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `name` varchar(64) NOT NULL COMMENT '话题展示名称',
  `normalized_name` varchar(64) NOT NULL COMMENT '用于去重检索的标准化话题名称',
  `view_num` bigint unsigned NOT NULL default 0 COMMENT '话题浏览量',
  `comment_num` bigint unsigned NOT NULL default 0 COMMENT '话题讨论度',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_normalized_name` (`normalized_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='话题';

CREATE TABLE `work_topic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '作品话题关系主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `work_id` bigint unsigned NOT NULL COMMENT '作品ID',
  `topic_id` bigint unsigned NOT NULL COMMENT '话题ID',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '话题展示顺序，从0开始',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_work_topic` (`work_id`,`topic_id`),
  KEY `idx_topic_work` (`topic_id`,`work_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='作品话题';

CREATE TABLE `media_process_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '媒体处理任务主键ID',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `work_id` bigint unsigned NOT NULL COMMENT '待处理作品ID',
  `task_type` varchar(32) NOT NULL COMMENT '任务类型：process_work处理作品媒体',
  `status` varchar(32) NOT NULL DEFAULT 'pending' COMMENT '任务状态：pending/processing/succeeded/failed',
  `stage` varchar(32) NOT NULL DEFAULT 'queued' COMMENT '处理阶段：queued/preparing/processing_image/processing_video/waiting_review/failed',
  `progress` int unsigned NOT NULL DEFAULT '0' COMMENT '处理进度百分比，范围0到100',
  `retry_count` int unsigned NOT NULL DEFAULT '0' COMMENT '任务重试次数',
  `error_message` varchar(2000) NOT NULL DEFAULT '' COMMENT '最近一次处理失败信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_work_task_type` (`work_id`,`task_type`),
  KEY `idx_status_updated` (`status`,`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='媒体处理任务';

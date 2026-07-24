-- 为已经执行过001迁移的数据库补充作品管理页统计字段。
ALTER TABLE `work`
  ADD COLUMN `view_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作品浏览次数' AFTER `review_reason`,
  ADD COLUMN `like_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作品点赞次数' AFTER `view_count`,
  ADD COLUMN `favorite_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作品收藏次数' AFTER `like_count`,
  ADD COLUMN `share_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作品转发次数' AFTER `favorite_count`;

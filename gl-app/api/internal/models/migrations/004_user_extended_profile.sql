-- 扩展用户个人资料字段；已有数据库升级时执行一次。
ALTER TABLE `user`
  ADD COLUMN `bio` varchar(256) NOT NULL DEFAULT '' COMMENT '用户简介' AFTER `sex`,
  ADD COLUMN `both_day` date DEFAULT NULL COMMENT '用户生日' AFTER `bio`;

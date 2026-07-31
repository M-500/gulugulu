-- 为用户个人中心增加头像对象Key；昵称字段已存在，无需重复添加。
ALTER TABLE `user`
  ADD COLUMN `avatar` varchar(512) NOT NULL DEFAULT '' COMMENT '用户头像在正式对象存储桶中的对象Key' AFTER `nickname`;

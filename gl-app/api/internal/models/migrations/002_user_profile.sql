-- 用户资料扩展：老库如果没有 avatar 字段，执行该语句。
ALTER TABLE `user`
  ADD COLUMN `avatar` varchar(512) NOT NULL DEFAULT '' COMMENT '用户头像在公共对象存储桶中的对象Key' AFTER `nickname`;

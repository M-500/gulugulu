# Repo 层约定

后端持久化统一使用 GORM，业务 Logic 和后台 Worker 只能依赖 Repo 接口，不允许直接执行 SQL，也不应直接使用 `ServiceContext.GormDB`。

## 目录职责

- `user_repo`：用户 DAO、Redis 缓存和缓存旁路策略。
- `work_repo`：作品创建、列表、详情、审核以及作品相关事务。
- `media_repo`：素材记录、处理任务以及媒体 Worker 状态流转事务。

## 依赖注入

MySQL 和 Redis 的客户端构造分别封装在 `pkg/gormx`、`pkg/redisx`。`ServiceContext` 只负责调用统一构造函数并把客户端注入 DAO、Cache 和 Repo，业务层只持有 Repo 接口。

## 缓存策略

用户详情采用 Cache Aside：

1. 优先读取 Redis。
2. 缓存未命中或 Redis 暂时不可用时回源 MySQL。
3. MySQL 查询成功后回写 Redis，缓存失败只记录日志，不阻断主流程。
4. 更新数据库成功后删除缓存，避免数据库写入失败时提前删除有效缓存。

作品和媒体任务以强一致状态流转为主，目前不缓存；如果后续增加作品缓存，必须在对应 Repo 内完成失效处理，不能把缓存操作泄漏到 Logic。

## 事务边界

跨表写入由 Repo 自己开启 GORM 事务。例如创建作品会在同一事务内校验并锁定素材、写入作品关系和话题、绑定素材并创建处理任务；媒体状态切换也会同步更新任务、作品和素材状态。

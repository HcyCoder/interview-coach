# ADR-002: 后端 API 网关使用 Go + Hertz + GORM

## Status
Accepted

## Date
2026-05-07

## Context
后端需要同时承担认证、会话管理、历史记录、任务编排和对外 API 网关职责。它需要高并发、低延迟、稳定的中间件生态，并且要与 MySQL 事务模型和 Redis/RabbitMQ 协作。

## Decision
使用 Go 作为主后端语言，采用 Hertz 处理 HTTP 入口，GORM 处理 MySQL 持久化。业务逻辑集中在 `service` 层，HTTP 层只负责校验、鉴权、DTO 转换和协议适配。

## Alternatives Considered

### Gin / Echo / Fiber
- Pros: 生态成熟，上手快。
- Cons: 相比 Hertz，在高并发场景和工程化约束上没有明显优势。
- Rejected: 该项目强调 API 网关稳定性和吞吐，Hertz 更适合作为默认选型。

### 直接在前端做 BFF
- Pros: 少一层服务。
- Cons: 业务状态机、异步任务和审计日志不适合放在前端侧。
- Rejected: 破坏单一事实来源和服务边界。

### 纯 SQL 手写数据访问
- Pros: 控制力强。
- Cons: 重复样板代码多，实体和迁移维护成本高。
- Rejected: GORM 足以覆盖本项目的 CRUD + 事务需求。

## Consequences
- API 入口、状态机和权限逻辑统一收敛到 Go 服务，前端只关心界面。
- GORM 迁移和查询抽象提高生产效率，但复杂查询仍需谨慎控制。
- Hertz 带来更高的性能余量，适合面试会话这种高频交互场景。

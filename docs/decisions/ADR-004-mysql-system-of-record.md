# ADR-004: MySQL 作为系统事实来源，并按服务划分 schema

## Status
Accepted

## Date
2026-05-07

## Context
产品需要强一致的账号、会话、历史记录、报告和审计数据。数据模型是高度关系化的，且面试记录、报告和权限访问都需要事务保证。与此同时，AI 服务和业务后端的关注点不同，不应共享同一组表的写权限。

## Decision
使用单一 MySQL 集群作为关系数据库，但按服务划分 schema：`app` 由 Go 后端拥有，`rag` 由 AI 服务拥有。敏感 PII 字段采用应用层加密后入库，查询所需的非敏感派生字段单独保存。

## Alternatives Considered

### PostgreSQL
- Pros: JSON 和扩展能力强。
- Cons: 本方案已被技术栈约束为 MySQL。
- Rejected: 不符合强制技术栈要求。

### MongoDB
- Pros: 灵活。
- Cons: 会话、报告、权限、令牌和审计是典型关系模型。
- Rejected: 事务和约束表达不如关系库直接。

### 各服务各自独立数据库实例
- Pros: 边界更强。
- Cons: 运维复杂度更高，当前规模没有必要拆成多个集群。
- Rejected: 采用单集群 + schema 级隔离即可满足边界和效率。

## Consequences
- 关系约束、唯一键、事务和审计都能稳定落地。
- 服务边界通过 schema 隔离，而不是靠约定。
- 需要明确迁移归属，避免跨服务随意改表。

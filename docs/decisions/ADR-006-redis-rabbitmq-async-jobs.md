# ADR-006: 缓存与异步任务采用 Redis + RabbitMQ

## Status
Accepted

## Date
2026-05-07

## Context
简历解析和复盘报告生成是长耗时任务，必须异步化；同时登录态、幂等键、限流和短期任务状态也需要高速存取。单靠一个系统无法同时覆盖“缓存”和“可靠队列”两类需求。

## Decision
Redis 负责缓存、幂等键、限流、短期状态和会话辅助数据；RabbitMQ 负责 durable job queue，承载 `resume_parse`、`report_generate` 等异步任务。Go 后端负责生产消息和更新任务状态，AI 服务负责消费并回写结果。

## Alternatives Considered

### Kafka
- Pros: 吞吐高，适合事件流。
- Cons: 维护成本高于当前需求，报告生成并不需要完整流平台。
- Rejected: 该项目更需要明确的任务投递语义，而不是流处理平台。

### 只用 Redis Queue / Streams
- Pros: 简单。
- Cons: 对任务持久性和复杂重试策略不如专用消息队列稳妥。
- Rejected: 复盘报告属于必须落地的后台任务，RabbitMQ 更合适。

### 只用 RabbitMQ
- Pros: 足够覆盖队列。
- Cons: 缺少低延迟缓存能力，登录态和幂等控制会变差。
- Rejected: Redis 仍然是必要的补充组件。

## Consequences
- 异步任务可以重试、可观测、可回放。
- Redis 降低数据库压力，并支持短期状态查询。
- 需要明确消息幂等和消费确认策略，防止重复生成报告。

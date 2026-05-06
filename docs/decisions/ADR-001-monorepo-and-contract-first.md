# ADR-001: Monorepo 与 Contract-First 边界

## Status
Accepted

## Date
2026-05-07

## Context
Interview Coach 由三个运行时组成：Next.js 前端、Go 业务后端、Python AI 服务。它们的语言、依赖和部署节奏不同，但共享同一套用户流程和数据语义。如果不先固定仓库与契约边界，接口会快速漂移，前后端和 AI 服务会反复互相阻塞。

## Decision
采用 monorepo，根目录固定为 `frontend/`、`backend/`、`ai-service/`、`infra/`、`contracts/` 和 `docs/`。所有跨服务接口以 `contracts/openapi/*.yaml` 作为单一事实源，生成前端 client、后端 DTO 校验和 AI 服务请求模型。

## Alternatives Considered

### Polyrepo
- Pros: 各服务独立演进，权限边界清晰。
- Cons: 契约同步成本高，开发联调频繁跨仓库，对小团队不友好。
- Rejected: 当前阶段优先降低协作成本，不优先追求组织级隔离。

### 单体应用
- Pros: 开发初期最简单。
- Cons: Go、Python、Node 运行时被强耦合，AI 依赖会污染业务后端。
- Rejected: 不符合服务职责分离，也不利于后续水平扩展。

## Consequences
- 接口变更必须先改 `contracts/`，再生成代码，减少漂移。
- 本地开发可以通过 Docker Compose 一键联调三套运行时。
- 架构更清晰，但需要维护生成脚本和契约校验流程。

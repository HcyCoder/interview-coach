# ADR-003: AI 编排与 RAG 服务独立为 Python FastAPI

## Status
Accepted

## Date
2026-05-07

## Context
简历解析、知识检索、提示词编排、结构化输出校验和报告生成都强依赖 Python 生态中的文本处理和 AI 工具链。把这些能力塞进 Go 后端会增加依赖复杂度，也会让模型代码与业务代码互相污染。

## Decision
将 AI 能力独立为 `ai-service`，使用 Python + FastAPI 实现。该服务负责简历解析、JD 结构化、RAG 检索、追问生成、报告生成和知识库向量化；Go 后端只通过内部 REST 或队列调用它。

## Alternatives Considered

### 直接把 AI 逻辑写进 Go
- Pros: 少一次网络调用。
- Cons: Python 生态的解析、向量化、推理工具链优势无法利用。
- Rejected: 维护成本高，迭代速度慢。

### 用前端 Server Actions 直接调模型
- Pros: 路径短。
- Cons: 会把高延迟、易失败的 AI 逻辑暴露给用户请求链。
- Rejected: 不符合服务隔离，也不利于审计和扩展。

### serverless / 外部编排平台
- Pros: 快速搭建。
- Cons: 复杂状态流和 RAG 依赖会变得难以调试和回放。
- Rejected: 需要可控、可追踪、可本地复现的编排服务。

## Consequences
- AI 相关依赖与业务后端解耦，部署和扩容可以独立进行。
- 需要额外维护内部鉴权、请求追踪和队列消费者。
- 模型输出校验、prompt 版本和检索证据都可以在 Python 层集中处理。

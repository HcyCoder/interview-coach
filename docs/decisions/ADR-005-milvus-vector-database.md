# ADR-005: 向量数据库采用 Milvus

## Status
Accepted

## Date
2026-05-07

## Context
RAG 需要对面经、八股文、项目片段和候选人上下文做语义召回，并且要支持元数据过滤、版本控制和后续扩容。单纯靠关系型数据库做向量检索会在性能和可维护性上迅速碰到上限。

## Decision
使用 Milvus 作为唯一向量数据库，`ai-service` 负责 embedding、入库、召回和重排。Milvus 中只保存向量与检索必要的 scalar metadata，chunk 原文和业务语义仍由 MySQL 的 `rag` schema 保存。

## Alternatives Considered

### pgvector
- Pros: 和关系库绑定紧密。
- Cons: 受限于关系库扩展方式，不适合作为独立 RAG 服务的主检索层。
- Rejected: 项目明确需要独立向量数据库。

### Elasticsearch 向量检索
- Pros: 搜索能力强。
- Cons: 主要面向全文检索，向量检索与过滤的建模不如专用向量库直接。
- Rejected: 目标是稳定的语义检索，不是搜索引擎化的混合检索平台。

### Pinecone / 云厂商托管向量库
- Pros: 运维简单。
- Cons: 对本地开发和自托管不友好，成本和可控性较差。
- Rejected: 该项目需要 Docker Compose 和 Kubernetes YAML 的自托管路径。

## Consequences
- 召回路径清晰：chunk -> embedding -> Milvus -> rerank -> citation。
- 检索性能可独立优化，不会拖累业务数据库。
- 需要维护索引重建、embedding 版本和知识库同步流程。

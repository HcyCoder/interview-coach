# ADR-007: 文件与中间产物使用 S3 兼容对象存储

## Status
Accepted

## Date
2026-05-07

## Context
系统必须存储简历 PDF、解析中间产物、生成的报告附件和未来可能的导出文件。把二进制大对象塞进 MySQL 会让数据库膨胀，并且不利于生命周期管理。

## Decision
所有二进制文件统一放到 S3 兼容对象存储中。本地开发使用 MinIO，生产环境接入云厂商 S3/OSS 类存储。数据库只保存对象 key、hash、大小、content-type 和生命周期状态。

## Alternatives Considered

### 直接存 MySQL BLOB
- Pros: 最直接。
- Cons: 不利于备份、缓存和大文件处理。
- Rejected: PDF 和报告附件都属于典型对象存储场景。

### 本地磁盘
- Pros: 实现快。
- Cons: 不能水平扩容，也不适合生产多副本部署。
- Rejected: 与 Kubernetes 和多实例部署冲突。

### 省略对象存储
- Pros: 少一个组件。
- Cons: 无法安全、持久地保存上传文件。
- Rejected: 业务流程本身就依赖文件持久化。

## Consequences
- 上传、下载、扫描和生命周期管理都可以独立于应用容器。
- `resumes` 表只保存元数据，不承担大文件压力。
- 本地开发和生产环境的文件访问路径保持一致。

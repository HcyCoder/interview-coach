# ADR-008: 认证采用邮箱密码 + JWT Access Token + Refresh Token

## Status
Accepted

## Date
2026-05-07

## Context
MVP 需要最小可用的账号体系，但同时要兼顾前后端分离、跨页面会话保持和安全注销。纯前端态会话不利于多设备访问；纯长效 token 又不利于撤销和轮换。

## Decision
采用邮箱 + 密码登录。短期 access token 通过响应体传递，refresh token 通过 `httpOnly`、`Secure`、`SameSite=Lax` Cookie 传输，并在 `refresh_tokens` 表中保存 hash 和轮换状态。访问令牌过期后由刷新接口轮换。

## Alternatives Considered

### 仅用服务端 session cookie
- Pros: 简单。
- Cons: 前后端分离下跨服务扩展和移动端复用不够灵活。
- Rejected: 需要显式 token 模式更适合该架构。

### 第三方 OAuth 登录
- Pros: 用户体验好。
- Cons: 依赖外部身份源，且不符合当前 MVP 的最小边界。
- Rejected: 不是本阶段必须项。

### 永久 JWT
- Pros: 实现快。
- Cons: 无法优雅撤销，安全边界差。
- Rejected: 不满足最基本的会话治理要求。

## Consequences
- 前端可以安全地在页面间维持登录状态。
- 后端可以按设备撤销 refresh token，支持强制登出。
- 需要维护 token 轮换、防重放和 cookie 安全属性。

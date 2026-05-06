# Spec: Interview Coach

## Objective
Interview Coach 是一个面向候选人的 AI 面试模拟与复盘平台，帮助用户把简历 PDF 和目标岗位 JD 转成结构化画像、定制化面试计划、多轮技术面试和可追溯的复盘报告。

用户故事：
- 候选人上传简历 PDF 后，系统自动解析出姓名、技术栈、项目经验、教育/工作经历等结构化 JSON。
- 候选人输入目标岗位 JD 后，系统生成定制化面试计划，包括重点考察方向、题目顺序和建议深度。
- 候选人开始面试后，AI 面试官基于简历 + JD 发起多轮追问，追问内容从 RAG 知识库检索面经/八股文后动态生成。
- 面试结束后，系统生成多维度评分报告，至少覆盖技术深度、沟通清晰度、项目描述质量和改进建议。
- 候选人可查看历史面试记录、评分详情和复盘建议。

MVP 范围：
- In scope：基础登录/账号体系、PDF 上传、画像提取、JD 解析、面试计划生成、多轮 AI 面试、评分报告、历史记录保存与查询。
- Out of scope：语音面试、coding sandbox、企业级多租户、公开分享、完整知识库运营后台、扫描版 PDF 的 OCR 兜底。

非功能性需求：
- 性能：10MB 以内 PDF 的解析 + 画像生成 p95 < 60s；首轮 AI 回复 p95 < 3s；报告生成 p95 < 30s。
- 安全：上传文件类型/大小校验、恶意内容扫描、用户级权限隔离、PII 加密存储、模型输出 schema 校验、审计日志。
- 可扩展性：解析/嵌入/报告生成走异步队列；Web/API 保持无状态；RAG 检索与评分服务可水平扩展。

用户交互流程：
1. 候选人登录。
2. 上传简历 PDF，系统展示解析结果并允许确认后继续。
3. 输入 JD，系统生成面试计划并让用户预览。
4. 启动面试，AI 根据简历 + JD 发题并从知识库检索追问依据。
5. 结束面试后，系统生成评分报告和复盘建议。
6. 用户进入历史页，查看过往面试详情并继续复盘。

后续迭代计划：
- 迭代 2：加入 OCR 识别、简历画像手动编辑、报告导出。
- 迭代 3：加入语音面试、实时转写和表达能力评分。
- 迭代 4：加入知识库后台、题库运营、团队复盘和数据分析。
- 迭代 5：加入多语言面试、ATS/招聘流程集成和企业级权限。

## Tech Stack
- Next.js 15、React 19、TypeScript 5、Node.js 22
- PostgreSQL 16 + Prisma 6
- Redis 7 + BullMQ
- pgvector
- Auth.js
- S3 兼容对象存储
- Zod
- OpenAI API
- Vitest、Testing Library、Playwright

## Commands
Build: `pnpm build`
Test: `pnpm test -- --coverage`
Lint: `pnpm lint`
Dev: `pnpm dev`

## Project Structure
- `apps/web/` → 候选人前台、登录、面试、历史记录、API 路由
- `apps/worker/` → 简历解析、向量化、RAG 检索、评分报告生成
- `packages/db/` → Prisma schema、迁移、数据库访问封装
- `packages/ai/` → Prompt、面试编排、评分规则、结构化输出
- `packages/rag/` → Chunk、embedding、召回、重排、引用
- `packages/shared/` → 类型、常量、校验器
- `tests/` → 单元/集成测试
- `e2e/` → 浏览器端到端测试
- `docs/` → PRD、架构说明、运行手册

## Code Style
- TypeScript 优先，`camelCase` 用于变量/函数，`PascalCase` 用于组件/类型，`use*` 用于 Hooks，`*Schema` 用于校验器，`*Service` 用于编排服务。
- 2 空格缩进，单引号，尾随逗号，禁止 `any` 和未校验的模型输出。
- 所有外部输入先过 Zod，再进入业务逻辑。

```ts
import { z } from 'zod';

const InterviewPlanSchema = z.object({
  focusAreas: z.array(z.string()).min(1),
  questionCount: z.number().int().min(3).max(10),
});

export function createInterviewPlan(input: GeneratePlanInput): InterviewPlan {
  return InterviewPlanSchema.parse({
    focusAreas: deriveFocusAreas(input.resume, input.jd),
    questionCount: 6,
  });
}
```

## Testing Strategy
- 单元测试：提示词拼装、评分规则、画像归一化、召回排序。
- 集成测试：PDF 上传到画像提取、JD 到面试计划、报告持久化、历史查询。
- E2E 测试：从上传到出报告再到历史查看的全链路。
- LLM 相关测试全部 mock；模型输出和报告结构做 schema/回归测试。
- 覆盖率目标：核心域模块 >= 80%，整体 >= 70%。

## Boundaries
- Always: 先校验输入和模型输出；保存 prompt/model 版本；记录审计日志；变更前先跑测试。
- Ask first: 新增依赖、改数据库 schema、改 CI/CD、换模型供应商、改存储/保留策略、调整评分维度或报告结构。
- Never: 提交密钥、跳过校验、直接信任模型自由文本、删除失败测试来让 CI 过、泄露其他用户的简历或历史记录。

## Success Criteria
- 候选人可在一个连续会话内完成上传、计划生成、面试、报告和历史保存。
- 候选人能上传简历 PDF 并拿到包含必需字段的结构化画像记录。
- 候选人能输入 JD 并生成面试计划后启动多轮 AI 面试。
- 每个追问都能追溯到知识库检索结果或来源 ID。
- 面试结束后能生成并保存评分报告，报告包含至少 4 个维度和复盘建议。
- 历史记录页能列出并打开过往面试详情，且跨设备登录后仍可访问。
- 系统在目标负载下满足上述性能、安全和隔离要求。

## Open Questions
- MVP 是否需要 OCR 识别扫描版 PDF，还是只支持文本型 PDF？
- 简历画像在开始面试前是否允许候选人手动编辑确认？
- 账号体系采用邮箱验证码、密码还是第三方登录？
- RAG 知识库的初始来源、更新频率和责任人是谁？
- 是否保存完整逐轮对话，以及对应的导出/删除策略是什么？

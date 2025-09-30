# Copilot 使用指南（为 AI 编码代理准备）

说明：此仓库当前没有可识别的源代码或说明文件（README、package.json、pyproject.toml 等未找到）。下面的指令面向在仓库被填充后，帮助 AI 代理快速发现架构、构建/测试/调试流程以及项目特有约定。请把本文件作为活文档：当项目添加实际代码时，更新此处的“关键文件/命令”以反映真实内容。

## 目标（简洁）
- 快速识别项目类型（Node/Python/Go/Java/多语言）及其入口点。
- 找到构建、测试、运行和 CI 的“真实命令”（若存在），并优先使用仓库内声明的脚本。
- 抽取项目约定（目录布局、配置位置、环境变量、集成点）并在自动化更改时遵守。

## 当下仓库状态（重要）
- 当前仓库为空或未包含可识别的工程文件。所有建议均基于可发现文件出现后的自动推断。
- 如果你是人类维护者，请把实际项目文件（至少 README 或 package.json / pyproject.toml）推送到仓库，之后本指南会被补充具体示例。

## 发现步骤（对 AI 代理的具体动作顺序 — 可自动执行）
1. 在仓库根及子目录查找这些文件以判定语言/框架：`package.json`, `pyproject.toml`, `requirements.txt`, `Pipfile`, `setup.py`, `go.mod`, `pom.xml`, `build.gradle`, `Dockerfile`, `README.md`。
2. 若发现 `package.json`：读取 `scripts` 节，优先使用 `build`/`test`/`start` 的脚本名称作为执行入口（不要假设命令，始终使用脚本条目）。
3. 若发现 `pyproject.toml` 或 `setup.py`：查找 `tool.poetry.scripts`/`entrypoints` 或 `if __name__ == '__main__'` 等入口点。
4. 查找 CI 配置：`.github/workflows/`、`azure-pipelines.yml`、`circleci/config.yml`，以找出常用构建/测试命令与环境矩阵。
5. 查找运行时配置与密钥线索：`.env`, `config/`, `infra/`, `terraform/`，以及 `secrets` 被引用的位置。
6. 检索测试目录（`tests/`, `spec/`, `__tests__`）并执行快速静态检查以判断测试框架和典型断言模式。

## 架构线索（当文件存在时应如何阅读）
- 单体 Web 应用：优先查看 `src/`（前端）或 `server/`、`app/`（后端），寻找 `index.js`/`app.py`/`main.go` 作为入口。
- 后端服务/微服务：查找 `cmd/`、`internal/`、`pkg/`（Go 风格），或 `services/`、`api/`（JS/Python）。
- 前端：查找 `public/`, `src/`，并检查 `webpack.config.js`, `vite.config.js`, `next.config.js`。
- 基础设施/部署：`Dockerfile`, `k8s/`, `infra/`, `terraform/` 指示部署目标及凭证/环境依赖。

## 项目特有约定（要记录并优先遵守）
- 优先使用仓库中显式声明的脚本（例如 `package.json#scripts`、Makefile 目标、pyproject 的脚本）。
- 环境变量读取通常发生在 `.env` 或 `config/*.yaml|json`。修改配置时请保持相同的命名与层级。
- 若仓库包含 `mono-repo` 风格（有 `packages/`、`workspaces` 或 `lerna`），只在单一包上下文内做变更或同时调整受影响的包并更新根级锁文件。

## 集成点与外部依赖（AI 需检查的位置）
- 第三方 API 客户端通常在 `lib/`、`clients/`、`integrations/`。查找 `API_KEY`, `TOKEN`, `BASE_URL` 的引用位置。
- 持久化：搜索 `orm`、`models`、`migrations`，并识别 DB 连接字符串来源（环境变量或 secrets 管理）。
- 消息/异步：搜索 `queue`, `rabbitmq`, `kafka`, `celery` 等关键词，确定消息边界与重试策略实现点。

## 代码模式示例（在对应文件出现时给出的具体提示）
- package.json 示例：优先调用 `npm run test` 中定义的脚本，而不是猜测 `npm test` 的行为；在 PR 中修改依赖时同步更新 `package-lock.json`/`yarn.lock`。
- Python 项目示例：若使用 `pytest`，请在 `pytest.ini` 或 `pyproject.toml` 中查阅额外选项（例如 markers 或 testpaths），运行测试时使用仓库声明的配置。

## 合并已有 AI 指南（若仓库已有旧文件）
- 如果仓库未来包含 `.github/copilot-instructions.md`、`AGENT.md` 或 `README.md` 中的 AI 指导，合并时保留具体命令、脚本名与业务上下文（不要覆盖维护者写的部署秘钥位置说明）。
- 优先保持人类作者写明的“生产环境差异”和“是否可在本地模拟”的段落不被移除。

## 交互样例（供人类维护者/审查员参考）
- 可对 AI 询问的短促示例：
  - “仓库中哪个文件声明了运行服务的命令？”
  - “CI 使用的测试命令是什么？在哪里定义？”
  - “有哪些外部服务（DB、队列、第三方 API）在代码中被引用？它们的配置在什么位置？”

## 变更与审核准则（给自动 PR 的提示）
- PR 应说明变更影响的运行命令、环境变量与对 CI 的影响。
- 若修改依赖或构建配置，提交时包含一个能重现的本地验证步骤（例如：在 CI 中已通过的脚本名与测试覆盖的简短说明）。

## 需要维护者补充的信息（请补全）
- 项目语言与框架（Node/Python/Go/Java/其它）。
- 核心运行命令（build/test/start）或指向定义这些命令的文件。
- 是否有特殊分支策略或 commit message 规范。

---

如果你愿意，我可以：
- 在你推送一些关键文件（至少 README 或 package.json / pyproject.toml）后，把本文件补充为包含 repo 内具体命令与示例的版本；或
- 现在把这个基础版本 commit 到仓库，并随后根据你提供的代码/文件继续迭代修改。

请告诉我你希望我做哪一步，或直接推送上述任一文件，我会立刻把本指南升级为针对性更强的版本。
# Trans - 命令行翻译与分支命名工具

一个简洁的 CLI 工具，支持中英文互译、交互式翻译，以及基于描述生成规范化 Git 分支名（始终调用 LLM，失败报错）。

## 特性
- 中文→英文翻译：`-e`
- 英文→中文翻译：`-c`
- 交互式翻译（演示版）：`-i`
- 分支名生成（调用 LLM）：`-g`

## 安装
- 全局安装：
  ```bash
  go install github.com/leijiangnan/trans@latest
  ```
- 本地编译：
  ```bash
  go build -o trans
  ```

## 快速使用

### 翻译
- 中文→英文：
  ```bash
  trans -e 你好世界
  # 输出: Hello world
  ```
- 英文→中文：
  ```bash
  trans -c hello world
  # 输出: 你好世界
  ```

### 交互式翻译（演示）
```bash
trans -i
```
- 支持中英文自动检测
- 输入 `help` 查看帮助；输入 `exit`/`quit` 退出

示例：
```bash
$ trans -i
=== 交互式翻译工具 ===
输入 'exit' 或 'quit' 退出
输入 'help' 查看帮助

请输入要翻译的文本: 你好
[DEMO] 正在翻译: "你好"...
翻译结果: "Hello"
```

### 生成 Git 分支名
```bash
trans -g "新增登录页面"
# 示例输出：feature/add-login-page

trans -g "修复支付错误 issue-123"
# 示例输出：fix/issue-123

trans -g "release v1.2.0"
# 示例输出：release/v1.2.0

trans -g "更新依赖与文档"
# 示例输出：chore/update-dependencies-and-documentation
```

## LLM 配置
`-g` 命令始终调用 LLM 生成分支名；若未配置或调用失败，将报错。

使用环境变量配置：
```bash
export OPENAI_BASE_URL=https://api.moonshot.cn/v1
export OPENAI_API_KEY=YOUR_KEY
export OPENAI_MODEL=kimi-k2-turbo-preview
```

说明：
- 不要给值加反引号或多余空格（例如不要写成：`export OPENAI_BASE_URL= `https://...``）。
- 提示词包含严格的命名规范（前缀、字符限制、禁止连续/首尾连字符、长度不超过 48）。
- LLM 输出会进行严格校验；不符合规范或过长将报错，不再降级回退。

## 约定式分支规范（摘要）
- 前缀：`feature/`（或 `feat/`）、`bugfix/`（或 `fix/`）、`hotfix/`、`release/`、`chore/`
- 基本规则：
  - 使用小写字母、数字与连字符；避免特殊字符、下划线或空格
  - 禁止连续、开头或结尾的连字符或点；`release` 分支允许版本号中的点（如 `release/v1.2.0`）
  - 清晰简洁；如有工单编号（`issue-123`、`ABC-123`），建议包含

## 注意事项
- 翻译功能使用 Google Translate（`client=gtx`），无需 API Key，但需网络连接
- 交互模式为演示版本，后续可替换为实际 AI 翻译 API
- `-g` 在中文描述下会自动翻译生成英文短语并规范化为分支名

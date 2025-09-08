# Trans - 中英文翻译命令行工具

一个简单的命令行翻译工具，支持中英文互译和交互式翻译。

### 环境变量(交互模式需要API_KEY)
```bash
export ANTHROPIC_BASE_URL=https://api.moonshot.cn/anthropic
export ANTHROPIC_API_KEY=sk-DZdvzKifi3BbalHjXFInqkG
````

## 安装

```bash
go install github.com/leijiangnan/trans
```

## 使用方法

### 1. 直接翻译模式

#### 中文翻译成英文
```bash
trans -e 你好世界
```
输出: Hello World

#### 英文翻译成中文
```bash
trans -c hello world
```
输出: 你好世界

### 2. 交互式翻译模式

```bash
trans -i
```

进入交互模式后，你可以：
- 直接输入文本进行翻译
- 支持中英文自动检测
- 输入 `help` 查看帮助
- 输入 `exit` 或 `quit` 退出

#### 交互模式示例
```bash
$ trans -i
=== 交互式翻译工具 ===
输入 'exit' 或 'quit' 退出
输入 'help' 查看帮助

请输入要翻译的文本: 你好
[DEMO] 翻译结果: "Hello"

请输入要翻译的文本: hello world
[DEMO] 翻译结果: "你好世界"

请输入要翻译的文本: exit
感谢使用，再见！
```

## 使用示例

```bash
# 中文翻译成英文
trans -e "你好，这是一个测试"
# 输出: Hello, this is a test

# 英文翻译成中文
trans -c "Good morning, how are you?"
# 输出: 早上好，你好吗？

# 进入交互模式
trans -i
```

## 注意事项

- 需要网络连接才能进行翻译
- 使用Google翻译API，无需API密钥
- 支持长文本翻译
- 交互模式目前为演示版本，后续将集成AI翻译API
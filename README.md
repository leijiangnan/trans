# Trans - 中英文翻译命令行工具

一个简单的命令行翻译工具，支持中英文互译。

## 安装

```bash
go install github.com/leijiangnan/trans
```

## 使用方法

### 中文翻译成英文
```bash
trans -e 你好世界
```
输出: Hello World

### 英文翻译成中文
```bash
trans -c Hello World
```
输出: 你好世界

## 使用示例

```bash
# 中文翻译成英文
trans -e "你好，这是一个测试"
# 输出: Hello, this is a test

# 英文翻译成中文
trans -c "Good morning, how are you?"
# 输出: 早上好，你好吗？
```

## 注意事项

- 需要网络连接才能进行翻译
- 使用Google翻译API，无需API密钥
- 支持长文本翻译
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// GenerateBranchNameWithLLM 调用LLM根据描述生成分支名，失败则返回错误
// 需要环境变量：
// - OPENAI_API_KEY（必需）
// - OPENAI_BASE_URL（可选，默认 https://api.openai.com/v1）
// - OPENAI_MODEL（可选，默认 gpt-4o-mini）
func GenerateBranchNameWithLLM(description string) (string, error) {
	apiKey := cleanEnvValue(os.Getenv("OPENAI_API_KEY"))
	if strings.TrimSpace(apiKey) == "" {
		return "", fmt.Errorf("missing OPENAI_API_KEY")
	}
	base := cleanEnvValue(os.Getenv("OPENAI_BASE_URL"))
	if strings.TrimSpace(base) == "" {
		base = "https://api.openai.com/v1"
	}
	model := cleanEnvValue(os.Getenv("OPENAI_MODEL"))
	if strings.TrimSpace(model) == "" {
		model = "gpt-4o-mini"
	}

	systemPrompt := `你是一个 Git 分支命名助理，请严格遵守以下规范并仅输出分支名字符串：
根据描述自动选择合适的分支前缀（从下面选择一个）：
- feature/ (或 feat/)：用于新功能（例如 feature/add-login-page, feat/add-login-page）
- bugfix/ (或 fix/)：用于错误修复（例如 bugfix/fix-header-bug, fix/header-bug）
- hotfix/：用于紧急修复（例如 hotfix/security-patch）
- release/：用于准备发布的分支（例如 release/v1.2.0）
- chore/：用于非代码任务，如依赖项、文档更新（例如 chore/update-dependencies）
基本规则：
- 仅凭分支名称就能清楚地了解代码更改的目的。
- 使用小写字母、数字、连字符和点：分支名称应全部使用小写字母 (a-z)、数字 (0-9) 及连字符 (-) 进行分隔。避免使用特殊字符、下划线或空格。对于 release 分支，可以在描述中使用点 (.) 来表示版本号（例如 release/v1.2.0）。
- 禁止连续、开头或结尾的连字符或点：确保连字符和点不能连续出现（例如 feature/new--login, release/v1.-2.0），也不能出现在描述的开头或结尾（例如 feature/-new-login, release/v1.2.0.）。
- 保持清晰简洁：分支名称应简明扼要，清楚表达工作的内容和目的。
- 包含工单编号：如果适用，应包含项目管理工具中的工单编号，以便于追踪。例如，对于工单 issue-123，分支名称可以是 feature/issue-123-new-login。
输出要求：
- 仅输出分支名（不包含解释或附加文本）。
- 分支名长度不要太长，整体不超过 48 个字符。`

	userPrompt := "Description: " + strings.TrimSpace(description)

	req := openAIChatRequest{
		Model: model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.2,
		MaxTokens:   32,
	}

	buf, _ := json.Marshal(req)
	url := strings.TrimRight(base, "/") + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 读取部分响应体供诊断
		var msg string
		var out map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&out)
		if out != nil {
			if m, ok := out["error"].(map[string]interface{}); ok {
				if s, ok := m["message"].(string); ok {
					msg = s
				}
			}
		}
		if msg == "" {
			msg = resp.Status
		}
		// 常见错误指引
		switch resp.StatusCode {
		case 401:
			return "", fmt.Errorf("llm http 401 unauthorized: %s. 请检查 OPENAI_API_KEY 是否有效，并确认 OPENAI_BASE_URL 与该 Key 对应的厂商一致（例如 Moonshot/Kimi 常见为 https://api.moonshot.cn/v1）。", msg)
		case 404:
			return "", fmt.Errorf("llm http 404 not found: %s. 请检查 OPENAI_BASE_URL/chat/completions 路径是否正确。", msg)
		case 429:
			return "", fmt.Errorf("llm http 429 rate limited: %s. 请稍后重试或降低请求频率。", msg)
		default:
			return "", fmt.Errorf("llm http %d: %s", resp.StatusCode, msg)
		}
	}

	var out openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("llm empty choices")
	}
	name := strings.TrimSpace(out.Choices[0].Message.Content)
	if name == "" {
		return "", fmt.Errorf("llm empty content")
	}

	// 后处理与严格校验（不做规则降级）
	name = postProcessBranchName(name)
	if name == "" {
		return "", fmt.Errorf("llm output invalid by policy")
	}
	return name, nil
}

func cleanEnvValue(s string) string {
	return strings.Trim(s, " \t\r\n\"'`")
}

func postProcessBranchName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// 校验：必须且只能包含一个前缀斜杠
	allowed := []string{"main", "master", "develop", "feature/", "feat/", "bugfix/", "fix/", "hotfix/", "release/", "chore/"}
	hasPrefix := false
	for _, p := range allowed {
		if strings.HasPrefix(s, p) {
			hasPrefix = true
			break
		}
	}
	if !hasPrefix {
		return ""
	}
	if strings.Count(s, "/") != 1 {
		return ""
	}
	if strings.HasSuffix(s, "/") {
		return ""
	}

	// 校验字符集
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c == '/' || c == '-' || c == '.' || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return ""
		}
	}
	// 禁止连续、开头或结尾的连字符或点
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, ".") || strings.HasSuffix(s, "-") || strings.HasSuffix(s, ".") {
		return ""
	}
	if strings.Contains(s, "--") || strings.Contains(s, "..") {
		return ""
	}

	// release 分支：版本号中允许点，且需合法
	if strings.HasPrefix(s, "release/") && parseReleaseVersion(s) == "" {
		return ""
	}
	// 非 release 分支不应包含点
	if !strings.HasPrefix(s, "release/") && strings.Contains(s, ".") {
		return ""
	}

	// 长度限制：整体不超过 48 字符
	if len(s) > 48 {
		return ""
	}
	return s
}

// 提取 release 版本号，支持 vX.Y.Z 或 X.Y.Z
func parseReleaseVersion(s string) string {
	re := regexp.MustCompile(`\brelease\/v?(\d+\.\d+\.\d+)\b`)
	m := re.FindStringSubmatch(s)
	if len(m) == 2 && m[1] != "" {
		return "v" + m[1]
	}
	return ""
}

// OpenAI API 结构
type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

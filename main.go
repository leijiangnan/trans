package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

func translate(text, targetLang string) (string, error) {
	// 使用Google翻译API
	apiURL := "https://translate.googleapis.com/translate_a/single"
	params := url.Values{}
	params.Add("client", "gtx")
	params.Add("sl", "auto")
	params.Add("tl", targetLang)
	params.Add("dt", "t")
	params.Add("q", text)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("翻译请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析Google翻译API的响应
	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(result) > 0 {
		if translations, ok := result[0].([]interface{}); ok && len(translations) > 0 {
			var translated strings.Builder
			for _, t := range translations {
				if parts, ok := t.([]interface{}); ok && len(parts) > 0 {
					if text, ok := parts[0].(string); ok {
						translated.WriteString(text)
					}
				}
			}
			return translated.String(), nil
		}
	}

	return "", fmt.Errorf("无法解析翻译结果")
}

func main() {
	var (
		englishFlag     = flag.Bool("e", false, "将中文翻译成英文")
		chineseFlag     = flag.Bool("c", false, "将英文翻译成中文")
		interactiveFlag = flag.Bool("i", false, "进入交互式翻译模式")
	)
	flag.Parse()

	// 检查是否使用了交互模式
	if *interactiveFlag {
		InteractiveTranslate()
		return
	}

	if !*englishFlag && !*chineseFlag {
		fmt.Println("使用方法:")
		fmt.Println("  trans -e <中文文本>    # 将中文翻译成英文")
		fmt.Println("  trans -c <英文文本>    # 将英文翻译成中文")
		fmt.Println("  trans -i               # 进入交互式翻译模式")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("错误: 请提供要翻译的文本")
		os.Exit(1)
	}

	text := strings.Join(args, " ")

	var targetLang string
	if *englishFlag {
		targetLang = "en"
	} else if *chineseFlag {
		targetLang = "zh-CN"
	}

	translated, err := translate(text, targetLang)
	if err != nil {
		fmt.Printf("翻译失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(translated)
}

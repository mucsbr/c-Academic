package s2s

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func extractTextCode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(extractText(c))
	}
	return buf.String()
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "code" {
			buf.WriteString("`" + extractText(c) + "`")
			continue
		}
		buf.WriteString(extractText(c))
	}
	return buf.String()
}

// renderHTML renders the provided HTML content and extracts relevant sections.
func renderHTML(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic(err)
	}
	content := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		//fmt.Println(n.Type, n.Data)
		if n.Type == html.ElementNode && n.Data == "pre" {
			fmt.Println("Code Block:")
			fmt.Println(extractTextCode(n))
			tmpContent := extractText(n)
			if tmpContent != "" && tmpContent != "\n" {
				content += "```" + extractTextCode(n) + "```"
			}
		} else if n.Type == html.ElementNode && n.Data == "p" {
			fmt.Println("Paragraph:")
			fmt.Println(extractText(n))
			content += extractText(n) + "\n"
		} else if n.Type == html.ElementNode && n.Data == "li" {
			fmt.Println("text:")
			fmt.Println(extractText(n))
			content += extractText(n) + "\n"
		}
		//if n.Type == html.TextNode {
		//	//fmt.Println(n.Data)
		//	content += n.Data
		//}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return content
}
func DealRes(input string) string {
	return strings.TrimRight(renderHTML(input), "\n")
}

func ProcessCodeSegments(msg string, subStr string) string {
	codeIndices := findAllIndices(msg, subStr)

	var result strings.Builder
	prevIndex := 0

	for _, codeIndex := range codeIndices {
		nearestTag, tagName := findNearestTag(msg[:codeIndex])
		fmt.Println(nearestTag)

		if tagName == "pre" {
			// 替换 <code> 为 ```
			result.WriteString(msg[prevIndex:codeIndex])
			result.WriteString("```\n")
			prevIndex = codeIndex + len(subStr)
		} else if tagName == "p" {
			result.WriteString(msg[prevIndex:codeIndex])
			result.WriteString("`")
			prevIndex = codeIndex + len(subStr)
		}
	}

	// 添加最后一部分内容
	result.WriteString(msg[prevIndex:])
	return result.String()
}

func ProcessCodeSegmentsEx(msg string, subStr string) string {
	codeIndices := findAllIndices(msg, subStr)

	var result strings.Builder
	prevIndex := 0

	for _, codeIndex := range codeIndices {
		nearestTag, tagName := findNearestTag(msg[:codeIndex])
		fmt.Println(nearestTag)

		if tagName == "pre" {
			// 替换 <code> 为 ```
			//result.WriteString(msg[prevIndex:codeIndex])
			//result.WriteString("```")
			//prevIndex = codeIndex + len(subStr)
		} else if tagName == "p" {
			result.WriteString(msg[prevIndex:codeIndex])
			result.WriteString("`")
			prevIndex = codeIndex + len(subStr)
		}
	}

	// 添加最后一部分内容
	result.WriteString(msg[prevIndex:])
	return result.String()
}

// 找到所有 <code> 标签的位置
func findAllIndices(msg, subStr string) []int {
	var indices []int
	idx := strings.Index(msg, subStr)
	for idx != -1 {
		indices = append(indices, idx)
		start := idx + len(subStr)
		idx = strings.Index(msg[start:], subStr)
		if idx != -1 {
			idx += start // adjust the index relative to the original string
		}
	}
	return indices
}

// 向前查找最近的 <p> 和 <pre> 标签
func findNearestTag(subMsg string) (string, string) {
	pIndex := strings.LastIndex(subMsg, "<p>")
	preIndex := strings.LastIndex(subMsg, "<pre>")

	if pIndex > preIndex {
		return subMsg[pIndex:], "p"
	} else if preIndex > pIndex {
		return subMsg[preIndex:], "pre"
	}
	return "", ""
}

func DealLine(html string) string {
	//// 匹配 <li> 标签的正则表达式
	//re := regexp.MustCompile(`<li>(.*?)</li>`)
	//// 匹配可能存在的编号的正则表达式
	//numberRe := regexp.MustCompile(`^\s*\d+\.\s*`)
	//
	//// 查找所有 <li> 标签内容
	//matches := re.FindAllStringSubmatch(htmlContent, -1)
	//
	//// 替换后的内容
	//for i, match := range matches {
	//	// 去除 <li> 标签内的编号
	//	cleanedContent := numberRe.ReplaceAllString(match[1], "")
	//	replacement := fmt.Sprintf("%d.%s", i+1, cleanedContent)
	//	htmlContent = strings.Replace(htmlContent, match[0], replacement, 1)
	//}
	//
	//// 移除 <li> 标签后的 <ol> 标签
	//htmlContent = strings.Replace(htmlContent, "<ol>", "<ol>\n", 1)
	//htmlContent = strings.Replace(htmlContent, "</ol>", "\n</ol>", 1)
	// 创建正则表达式匹配<li>标签内容中的现有编号
	//re := regexp.MustCompile(`<li>\s*\d+\.\s*`)
	//html = re.ReplaceAllString(html, "<li>")
	//
	//// 查找所有<li>标签内容
	//re = regexp.MustCompile(`<li>(.*?)</li>`)
	//matches := re.FindAllStringSubmatch(html, -1)
	//
	//// 重新构建<li>标签并编号
	//for i, match := range matches {
	//	item := match[1]
	//	// 如果内容以```开头，则跳过处理
	//	if strings.HasPrefix(strings.TrimSpace(item), "`") {
	//		continue
	//	}
	//	newItem := fmt.Sprintf("%d. %s", i+1, strings.TrimSpace(item))
	//	html = strings.Replace(html, match[0], newItem, 1)
	//}
	//if html == ". " {
	//	return ""
	//}
	//
	//// 输出结果
	//return html
	// 处理<ol>标签中的<li>标签
	html = reformatList(html, "ol", true)

	// 处理<ul>标签中的<li>标签
	html = reformatList(html, "ul", false)

	return html
}

func reformatList(html, listType string, numbered bool) string {
	// 匹配对应listType标签内的内容
	re := regexp.MustCompile(fmt.Sprintf(`<%s>(.*?)</%s>`, listType, listType))
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		listContent := match[1]

		// 创建正则表达式匹配<li>标签内容中的现有编号
		reItem := regexp.MustCompile(`<li>\s*\d*\.*\s*`)
		listContent = reItem.ReplaceAllString(listContent, "<li>")

		// 查找所有<li>标签内容
		reItem = regexp.MustCompile(`<li>(.*?)</li>`)
		items := reItem.FindAllStringSubmatch(listContent, -1)

		// 重新构建<li>标签并编号
		var newListContent strings.Builder
		count := 1
		for _, item := range items {
			content := item[1]
			if strings.HasPrefix(strings.TrimSpace(content), "`") {
				newListContent.WriteString(item[0])
				continue
			}

			var newItem string
			if numbered {
				newItem = fmt.Sprintf("%d.%s", count, strings.TrimSpace(content))
			} else {
				newItem = fmt.Sprintf("- %s", strings.TrimSpace(content))
			}
			newListContent.WriteString(newItem)
			count++
		}

		html = strings.Replace(html, match[1], newListContent.String(), 1)
	}

	return html
}

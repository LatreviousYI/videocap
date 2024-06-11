package main

import (
    "fmt"
    "regexp"
)

func main() {
    // 示例字符串
    str := "{{today}}/{{deviceid}}-{{datetime}}{{jjjj}}"

    // 定义正则表达式
    re := regexp.MustCompile(`\{\{(.*?)\}\}`)

    // 查找所有匹配项
    matches := re.FindAllStringSubmatch(str, -1)

    // 打印所有匹配项
	fmt.Println(matches)
    for _, match := range matches {
        if len(match) > 1 {
            fmt.Println(match[1])
        }
    }
}
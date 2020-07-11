package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func modeFunc(mode string) mode {
	switch mode {
	case "test":
		return testMode
	case "line":
		return lineMode
	case "group":
		return groupMode
	default:
		return nil
	}
}

type mode func(string, string, io.StringWriter)

func main() {
	filePath := flag.String("f", "", "file path")
	mode := flag.String("m", "test", "execute mode, test/line/group")
	out := flag.String("o", "", "output path, default console")
	flag.Parse()
	if "" == *filePath {
		return
	}

	fmt.Println("file path: ", *filePath, " mode: ", *mode)

	dat, err := ioutil.ReadFile(*filePath)
	fmt.Println(*filePath)
	if err != nil {
		fmt.Println("file path error", err)
		return
	}

	var writer io.StringWriter

	if *out != "" {
		writer, err = os.OpenFile(*out, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("file path error", err)
			return
		}
	} else {
		writer = os.Stdout
	}

	text := string(dat)

	mof := modeFunc(*mode)
	if nil == mof {
		fmt.Println("mode not found")
		return
	}
	fmt.Println("input regex:")
	regexStr := ""
	for true {
		regexStr = strings.TrimSpace(readConsoleLine())
		if "" == regexStr {
			fmt.Println("input regex:")
			continue
		}
		fmt.Printf("regexStr=%s\n", regexStr)
		mof(text, regexStr, writer)
		fmt.Println("input regex:")
	}
}

func readConsoleLine() string {
	reader := bufio.NewReader(os.Stdin)
	data, _, e := reader.ReadLine()
	if e != nil {
		return ""
	}
	regexStr := string(data)
	return regexStr
}

// 只输出前5个匹配结果, 展示匹配内容和分组内容
func testMode(text string, regex string, writer io.StringWriter) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("regex error", err)
		return
	}
	result := reg.FindAllStringSubmatch(text, -1)
	writer.WriteString(fmt.Sprintf("total:%d\n", len(result)))
	for i, s := range result {
		if i > 4 {
			return
		}
		writer.WriteString(fmt.Sprintf("-----------%d-------------\n", i))
		writer.WriteString(fmt.Sprintf("%s\n", s[0]))
		for j, g := range s {
			if 0 == j {
				continue
			}
			writer.WriteString(fmt.Sprintf("group %d: %s\n", j, g))
		}
	}
}

// 只输所有结果, 从分组0开始到结尾. 以\t分割
func lineMode(text string, regex string, writer io.StringWriter) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("regex error", err)
		return
	}
	result := reg.FindAllStringSubmatch(text, -1)
	fmt.Printf("total:%d\n", len(result))
	for _, s := range result {
		for _, g := range s {
			writer.WriteString(fmt.Sprintf("\t%s", g))
		}
		writer.WriteString("\n")
	}
}

// 只输出分组, 从分组1开始到结尾. 以\t分割
func groupMode(text string, regex string, writer io.StringWriter) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("regex error", err)
		return
	}
	result := reg.FindAllStringSubmatch(text, -1)
	fmt.Printf("total:%d\n", len(result))
	for _, s := range result {
		for j, g := range s {
			if 0 == j {
				continue
			}
			writer.WriteString(fmt.Sprintf("\t%s", g))
		}
		writer.WriteString("\n")
	}
}

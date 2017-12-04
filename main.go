package main

import (
	"fmt"
	"regexp"
	"io/ioutil"
	"flag"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type testModel struct {
	InputStr     string
	FindGroupNum []int
	PrintNum     int
	RegexStr     string
}

func testModelParseParams(str string) (*testModel) {
	ss := strings.SplitN(str, "r:", 2)

	var testModel testModel
	testModel.InputStr = str
	testModel.PrintNum = -1
	testModel.FindGroupNum = make([]int, 0, 5)
	if len(ss) == 1 {
		testModel.RegexStr = ss[0]
	} else {
		testModel.RegexStr = ss[1]
		ps := strings.SplitN(ss[0], " ", 2)
		if len(ps) == 2 {
			testModel.PrintNum, _ = strconv.Atoi(ps[0])
			groupNums := strings.Split(ps[1], ",")
			for _, gns := range groupNums {
				gns, e := strconv.Atoi(strings.TrimSpace(gns))
				if e != nil {
					continue
				}
				testModel.FindGroupNum = append(testModel.FindGroupNum, gns)
			}
		} else {
			testModel.PrintNum, _ = strconv.Atoi(ss[0])
		}
	}

	return &testModel
}

func modelFunc(model string) func(*string, string) {
	switch model {
	case "test":
		return func(text *string, input string) {
			findF(text, testModelParseParams(input))
		}
	}
	return nil
}

func main() {
	filePath := flag.String("f", "", "file path")
	model := flag.String("m", "test", "excute model")
	//model := flag.StringVar
	flag.Parse()
	if "" == *filePath {
		return
	}

	fmt.Println("file path: ", *filePath, " model: ", *model)

	dat, err := ioutil.ReadFile(*filePath)
	fmt.Println(*filePath)
	if err != nil {
		fmt.Println("file path error", err)
		return
	}

	text := string(dat)

	mof := modelFunc(*model)
	if nil == mof {
		fmt.Println("model not found")
		return
	}
	fmt.Println("input regex:")
	regexStr := ""
	for true {
		regexStr = strings.TrimSpace(readConsoleLine())
		if "" == regexStr {
			fmt.Println("input regex")
			continue
		}
		fmt.Println("regexStr=", regexStr)
		mof(&text, regexStr)
		fmt.Println("input regex")
	}
}

func readConsoleLine() (string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, e := reader.ReadLine()
	if e != nil {
		return ""
	}
	regexStr := string(data)
	return regexStr
}

func findF(text *string, findPattern *testModel) {
	reg, _ := regexp.Compile(findPattern.RegexStr)
	result := reg.FindAllStringSubmatch(*text, findPattern.PrintNum)
	for i, s := range result {
		fmt.Println("-----------", i, "-------------")
		fmt.Println(s[0])
		if 0 == len(findPattern.FindGroupNum) {
			for j, g := range s {
				if 0 == j {
					continue
				}
				fmt.Println("group ", j, ": ", g)
			}
		} else {
			for _, gn := range findPattern.FindGroupNum {
				if len(s) >= gn {
					fmt.Println("group ", gn, ": ", s[gn])
				}
			}
		}
	}
}

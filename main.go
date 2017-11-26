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

type FindPattern struct {
	InputStr string
	FindGroupNum int
	PrintNum int
	RegexStr string
}

func parsePattern(str string) (*FindPattern) {
	ss := strings.SplitN(str," ",3)
	var pattern FindPattern
	pattern.InputStr = str
	pattern.PrintNum = -1
	pattern.FindGroupNum = -1
	if len(ss) == 3 {
		pattern.PrintNum,_ = strconv.Atoi(ss[0])
		pattern.FindGroupNum,_ = strconv.Atoi(ss[1])
		pattern.RegexStr = ss[2]
	} else if len(ss) == 2 {
		pattern.PrintNum,_ = strconv.Atoi(ss[0])
		pattern.RegexStr = ss[1]
	} else {
		pattern.RegexStr = ss[0]
	}
	return &pattern
}


func main()  {
	filePath := flag.String("f","file path","file path")
	flag.Parse()
	fmt.Println(*filePath)
	dat,err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println("file path error",err)
		return
	}

	text := string(dat)
	fmt.Println("input regex:")
	regexStr := ""

	for true {
		regexStr = strings.TrimSpace(readConsoleLine())
		if "" == regexStr {
			fmt.Println("input regex")
			continue
		}
		fmt.Println("regexStr=", regexStr)
		findF(&text,parsePattern(regexStr))
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

func findF(text *string,findPattern *FindPattern) {
	reg,_ := regexp.Compile(findPattern.RegexStr)
	result := reg.FindAllStringSubmatch(*text,findPattern.PrintNum)
	for i,s := range result {
		fmt.Println("-----------",i,"-------------")

		if findPattern.FindGroupNum <= 0 {
			fmt.Println(s[0])
			for j, g := range s {
				if 0 == j {
					continue
				}
				fmt.Println("group ", j, ": ", g)
			}
		} else {
			if len(s) >= findPattern.FindGroupNum {
				fmt.Println("group ", findPattern.FindGroupNum, ": ", s[findPattern.FindGroupNum])
			}
		}
	}
}
package main

import (
	"fmt"
	"regexp"
	//"io"
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
	if len(ss) == 3 {
		pattern.FindGroupNum,_ = strconv.Atoi(ss[0])
		pattern.PrintNum,_ = strconv.Atoi(ss[1])
		pattern.RegexStr = ss[2]
	} else if len(ss) == 2 {
		pattern.FindGroupNum,_ = strconv.Atoi(ss[0])
		pattern.RegexStr = ss[1]
	} else {
		pattern.RegexStr = ss[0]
	}
	return &pattern
}


func main()  {
	filePath := flag.String("f","file path","file path")
	flag.Parse()
	//filePath := "f://cache/test.txt"
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
		regexStr = readConsoleLine()
		//_, e := fmt.Scanf("%s",&regexStr)
		if strings.EqualFold("",regexStr) {
			fmt.Println("input regex")
		}
		fmt.Println("regexStr=", regexStr)
		find(text,regexStr)
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


func find(text string, pattern string) {
	reg,_ := regexp.Compile(pattern)
	result := reg.FindAllStringSubmatch(text,-1)
	for i,s := range result {
		fmt.Println("-----------",i,"-------------")
		if len(s) == 1 {
			fmt.Println(s[0])
		} else {
			fmt.Println(s[0])
			for j, g := range s {
				if 0 == j {
					continue
				}
				fmt.Println("group ",j,": ",g)
			}
		}

	}
}
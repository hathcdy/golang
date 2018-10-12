package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

var startPage int
var endPage int
var lineNumber int
var pageType int
var printDest string
var inputFileName string

func main() {
	flag.IntVarP(&startPage, "startNumber", "s", 1, "Input the start page")
	flag.IntVarP(&endPage, "endNumber", "e", 1, "Input the end page")
	flag.IntVarP(&lineNumber, "lineNumber", "l", 72, "Input the number of lines in a page")
	flag.StringVarP(&printDest, "printDest", "d", "", "The destination to print")

	flagF := flag.BoolP("flagF", "f", false, "If the pages are delimited by '\\f'")

	flag.Parse()

	if *flagF {
		pageType = 2
	} else {
		pageType = 1
	}

	if flag.NArg() == 1 {
		inputFileName = flag.Arg(0)
	} else {
		inputFileName = ""
	}

	/*
		fmt.Printf("args = %s, num = %d\n", flag.Args(), flag.NArg())
			for i := 0; i != flag.NArg(); i++ {
				fmt.Printf("ard[%d] = %s\n", i, flag.Arg(i))
			}
	*/

	/*
		fmt.Printf("startPage = %d\n", startPage)
		fmt.Printf("endPage = %d\n", endPage)
		fmt.Printf("lineNumber = %d\n", lineNumber)
		fmt.Printf("pageType = %d\n", pageType)
		fmt.Printf("printDest = %s\n", printDest)
		fmt.Printf("inputFileName = %s\n", inputFileName)
	*/

	Validate(startPage, endPage, lineNumber, pageType, printDest, inputFileName, flag.NArg())
	execSelpg(startPage, endPage, lineNumber, pageType, printDest, inputFileName)
}

func Validate(startPage int, endPage int, lineNumber int, pageType int, printDest string, inputFileName string, notFlageNum int) {
	seValid := startPage >= 1 && startPage <= endPage
	nargValid := notFlageNum == 0 || notFlageNum == 1
	typeValid := !(pageType == 2 && lineNumber != -1)
	if !seValid || !nargValid || !typeValid {
		println("Invalid input, use '-h' for help")
		os.Exit(1)
	}
}

func execSelpg(startPage int, endPage int, lineNumber int, pageType int, printDest string, inputFileName string) {
	curPage := 1
	curLine := 0

	fin := os.Stdin
	fout := os.Stdout
	var inpipe io.WriteCloser
	var err error

	if inputFileName != "" {
		fin, err = os.Open(inputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "selpg: could not open input file \"%s\"\n", inputFileName)
			fmt.Println(err)
			os.Exit(1)
		}
		defer fin.Close()
	}

	if printDest != "" {
		cmd := exec.Command("cat", "-n")
		inpipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer inpipe.Close()
		cmd.Stdout = fout
		cmd.Start()
	}

	if pageType == 1 {
		sc := bufio.NewScanner(fin)
		for sc.Scan() {
			if curPage >= startPage && curPage <= endPage {
				fout.Write([]byte(sc.Text() + "\n"))
				if printDest != "" {
					inpipe.Write([]byte(sc.Text() + "\n"))
				}
			}
			curLine++
			if curLine%lineNumber == 0 {
				curPage++
				curLine = 0
			}
		}
	} else {
		rd := bufio.NewReader(fin)
		for {
			pageStr, err := rd.ReadString('\f')
			if err != nil || err == io.EOF {
				if err == io.EOF {
					if curPage >= startPage && curPage <= endPage {
						fmt.Fprintf(fout, "%s", pageStr)
					}
				}
				break
			}
			pageStr = strings.Replace(pageStr, "\f", "", -1)
			curPage++
			if curPage >= startPage && curPage <= endPage {
				fmt.Fprintf(fout, "%s", pageStr)
			}
		}
	}

	if curPage < endPage {
		fmt.Fprintf(os.Stderr, "./selpg: end page (%d) is greater than total page (%d)", endPage, curPage)
	}
}

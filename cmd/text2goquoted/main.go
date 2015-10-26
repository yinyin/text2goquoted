package main

import "os"
import "flag"
import "fmt"
import "log"

import "github.com/yinyin/text2goquoted/quoter"

func main() {
	var inputFilePath string
	var outputFilePath string
	var pkgName string
	var constNamePrefix string
	var keepPrefixSpace bool
	var keepSuffixSpace bool
	var keepNewLine bool
	
	flag.StringVar(&inputFilePath, "i", "", "input file path")
	flag.StringVar(&outputFilePath, "o", "", "output file path")
	flag.StringVar(&pkgName, "package", "", "package name")
	flag.StringVar(&pkgName, "p", "", "package name (short hand)")
	flag.StringVar(&constNamePrefix, "n", "-- ", "const name prefix")
	flag.BoolVar(&keepPrefixSpace, "keep_prefix_space", false, "keep prefix spaces")
	flag.BoolVar(&keepSuffixSpace, "keep_suffix_space", false, "keep suffix spaces")
	flag.BoolVar(&keepNewLine, "keep_new_line", false, "keep new line character")

	flag.Parse()
	
	if "" == inputFilePath {
		fmt.Fprint(os.Stderr, "ERR: input path not given.\n")
		return
	}
	if "" == outputFilePath {
		fmt.Fprint(os.Stderr, "ERR: output path not given.\n")
		return
	}
	if "" == pkgName {
		fmt.Fprint(os.Stderr, "ERR: package name not given.\n")
		return
	}
	
	fpIn, err := os.Open(inputFilePath)
	if nil != err {
		log.Fatal("cannot open input file for read:", err)
		return
	}
	defer fpIn.Close()
	
	fpOut, err := os.Create(outputFilePath)
	if nil != err {
		log.Fatal("cannot open output file for write:", err)
		return
	}
	defer fpOut.Close()
	
	err = quoter.QuoteText(fpOut, fpIn, pkgName, constNamePrefix, keepPrefixSpace, keepSuffixSpace, keepNewLine)
	if nil != err {
		log.Fatal("cannot translate text to quote string file:", err)
		return
	}
	return
}


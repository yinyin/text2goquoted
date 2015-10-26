package quoter

import "io"
import "bufio"
import "strconv"
import "strings"
import "unicode"

func getStringLead(count int) (l string) {
	if 0 == count {
		return "\""
	} else {
		return "\t+ \""
	}
}

func getStringFinish() (l string) {
	return "\"\n"
}

func convertFragmentedLine(w bufio.Writer, r bufio.Reader, keepPrefixSpace bool, keepSuffixSpace bool, line string, count int) (err error) {
	if !keepPrefixSpace {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
	}
	w.WriteString(getStringLead(count))
	w.WriteString(strconv.Quote(line))
	line, isPrefix, err := r.ReadLine()
	for nil == err {
		if isPrefix {
			w.WriteString(strconv.Quote(line))
		} else {
			if !keepSuffixSpace {
				line = strings.TrimRightFunc(line, unicode.IsSpace)
			}
			w.WriteString(strconv.Quote(line))
			w.WriteString(getStringFinish())
			break
		}
		line, isPrefix, err = r.ReadLine()
	}
	return nil
}

func convertSingleLine(w *bufio.Writer, keepPrefixSpace bool, keepSuffixSpace bool, line string, count int) {
	if !keepPrefixSpace {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
	}
	if !keepSuffixSpace {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	if 0 == len(line) {
		return
	}
	w.WriteString(getStringLead(count))
	w.WriteString(strconv.Quote(line))
	w.WriteString(getStringFinish())
}

func QuoteText(rd io.Reader, wr io.Writer, pkgName string, constNamePrefix string, keepPrefixSpace bool, keepSuffixSpace bool) (err error) {
	r := bufio.NewScanner(rd)
	w := bufio.NewWriter(wr)
	w.WriteString("package ")
	w.WriteString(pkgName)
	w.WriteString("\n\n")

	var currentConstName string
	var lineCount int = 0
	splitedContent := false

	for r.Scan() {
		line := r.Text()
		if strings.HasPrefix(line, constNamePrefix) {
			currentConstName = strings.TrimSpace(line[len(constNamePrefix):len(line)])
			w.WriteString("const ")
			w.WriteString(currentConstName)
			w.WriteString(" string = ")
			lineCount = 0
		} else if 0 != len(strings.TrimSpace(line)) {
			convertSingleLine(w, keepPrefixSpace, keepSuffixSpace, line, lineCount)
			lineCount = lineCount + 1
		}
	}
	return nil
}

package quoter

import "io"
import "bufio"
import "strconv"
import "strings"
import "unicode"

func getStringLead(count int) (l string) {
	if 0 == count {
		return ""
	} else {
		return "\t+ "
	}
}

func getStringFinish() (l string) {
	return "\n"
}

func convertSingleLine(w *bufio.Writer, keepPrefixSpace bool, keepSuffixSpace bool, keepNewLine bool, line string, count int) {
	if !keepPrefixSpace {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
	}
	if !keepSuffixSpace {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	if keepNewLine {
		line = strings.TrimRight(line, "\n") + "\n"
	}
	if 0 == len(line) {
		return
	}
	w.WriteString(getStringLead(count))
	w.WriteString(strconv.Quote(line))
	w.WriteString(getStringFinish())
}

func QuoteText(wr io.Writer, rd io.Reader, pkgName string, constNamePrefix string, keepPrefixSpace bool, keepSuffixSpace bool, keepNewLine bool) (err error) {
	r := bufio.NewScanner(rd)
	w := bufio.NewWriter(wr)
	w.WriteString("package ")
	w.WriteString(pkgName)
	w.WriteString("\n\n")

	var lineCount int = 0

	for r.Scan() {
		line := r.Text()
		if strings.HasPrefix(line, constNamePrefix) {
			currentConstName := strings.TrimSpace(line[len(constNamePrefix):len(line)])
			w.WriteString("const ")
			w.WriteString(currentConstName)
			w.WriteString(" string = ")
			lineCount = 0
		} else if 0 != len(strings.TrimSpace(line)) {
			convertSingleLine(w, keepPrefixSpace, keepSuffixSpace, keepNewLine, line, lineCount)
			lineCount = lineCount + 1
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return nil
}

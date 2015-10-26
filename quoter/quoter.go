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

func convertSingleLine(keepPrefixSpace bool, keepSuffixSpace bool, keepNewLine bool, line string, count int) (l string) {
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
		return ""
	}
	return strconv.Quote(line)
}

func outputLine(w *bufio.Writer, line string, suffix string, count int) (err error) {
	if count > 0 {
		if _, err = w.WriteRune('\t'); nil != err {
			return err
		}
	}
	if _, err = w.WriteString(line); nil != err {
		return err
	}
	if _, err = w.WriteString(suffix); nil != err {
		return err
	}
	if _, err = w.WriteRune('\n'); nil != err {
		return err
	}
	return nil
}

func QuoteText(wr io.Writer, rd io.Reader, pkgName string, constNamePrefix string, keepPrefixSpace bool, keepSuffixSpace bool, keepNewLine bool) (err error) {
	r := bufio.NewScanner(rd)
	w := bufio.NewWriter(wr)
	w.WriteString("package ")
	w.WriteString(pkgName)
	w.WriteString("\n\n")

	var bufferedLine string = ""
	var lineCount int = 0

	for r.Scan() {
		line := r.Text()
		if strings.HasPrefix(line, constNamePrefix) {
			if err := outputLine(w, bufferedLine, "", lineCount); nil != err {
				return err
			}
			currentConstName := strings.TrimSpace(line[len(constNamePrefix):])
			w.WriteString("const ")
			w.WriteString(currentConstName)
			w.WriteString(" string = ")
			bufferedLine = ""
			lineCount = 0
		} else if 0 != len(strings.TrimSpace(line)) {
			if len(bufferedLine) > 0 {
				if err := outputLine(w, bufferedLine, " +", lineCount); nil != err {
					return err
				}
				lineCount = lineCount + 1
			}
			bufferedLine = convertSingleLine(keepPrefixSpace, keepSuffixSpace, keepNewLine, line, lineCount)
		}
	}
	if len(bufferedLine) > 0 {
		if err := outputLine(w, bufferedLine, "", lineCount); nil != err {
			return err
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return w.Flush()
}

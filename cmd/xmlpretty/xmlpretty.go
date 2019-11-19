package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func indent(n int) string {
	return strings.Repeat(" ", n)
}

func attrs(as []xml.Attr) string {
	var b strings.Builder
	sep := ""
	for _, a := range as {
		b.WriteString(sep)
		b.WriteString(name(a.Name))
		b.WriteString("=\"")
		b.WriteString(a.Value)
		b.WriteString("\"")
		sep = " "
	}
	return b.String()
}

const (
	none = iota
	start
	end
	char
	comment
	proc
	directive
)

func name(n xml.Name) string {
	var b strings.Builder
	if n.Space != "" {
		b.WriteString(n.Space)
		b.WriteRune(':')
	}
	b.WriteString(n.Local)
	return b.String()
}

func main() {
	f := os.Stdin
	var err error
	if len(os.Args) > 1 {
		f, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("open failed: %s\n", err)
			os.Exit(-1)
		}
		defer f.Close()
	}
	d := xml.NewDecoder(f)
	prevType := none
	prevStartOrEnd := none
	close := ""
	ind := 0
	for {
		t, err := d.Token()
		if err != nil {
			fmt.Printf("err: %s\n", err)
			break
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if prevStartOrEnd == start {
				ind += 2
			}
			fmt.Printf("%s\n", close)
			fmt.Printf("%s<%s", indent(ind), name(se.Name))
			a := attrs(se.Attr)
			if a != "" {
				fmt.Printf(" %s", a)
			}
			close = ">"
			prevType = start
			prevStartOrEnd = start

		case xml.EndElement:
			if prevStartOrEnd == end {
				ind -= 2
			}
			switch prevType {
			case start:
				fmt.Print("/>")
			case char:
				fmt.Printf("</%s>", name(se.Name))
			case end:
				fmt.Printf("\n%s</%s>",
					indent(ind), name(se.Name))
			default:
				fmt.Print("\n<%s>",
					indent(ind), name(se.Name))
			}
			close = ""
			prevType = end
			prevStartOrEnd = end

		case xml.CharData:
			fmt.Printf("%s%s", close, se)
			close = ""
			prevType = char

		case xml.Comment:
			fmt.Printf("%s<!--%s-->", close, se)
			close = ""
			prevType = comment

		case xml.ProcInst:
			fmt.Printf("%s<?%s %s?>", close, se.Target, se.Inst)
			close = ""
			prevType = proc

		case xml.Directive:
			fmt.Printf("%s<!%s>", close, se)
			close = ""
			prevType = directive

		default:
			fmt.Printf("unknown element type %T\n", t)
			prevType = none
			break
		}
	}
}

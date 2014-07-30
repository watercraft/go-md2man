package mangen

import (
	"bytes"
	"fmt"

	"github.com/russross/blackfriday"
)

type Man struct{}

func ManRenderer(flags int) blackfriday.Renderer {
	return &Man{}
}

func (m *Man) GetFlags() int {
	return 0
}

func (m *Man) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	out.WriteString("\n.PP\n.RS\n\n.nf\n")
	escapeSpecialChars(out, text)
	out.WriteString("\n.fi\n")
}

func (m *Man) BlockQuote(out *bytes.Buffer, text []byte) {
	out.WriteString("\n.PP\n.RS\n")
	out.Write(text)
	out.WriteString("\n.RE\n")
}

func (m *Man) BlockHtml(out *bytes.Buffer, text []byte) {
	// a pretty lame thing to do...
	fmt.Errorf("man: BlockHtml not supported")
	out.Write(text)
}

func (m *Man) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()

	out.WriteString("\n.")
	switch level {
	case 1:
		out.WriteString("TH ")
	case 2:
		out.WriteString("SH ")
	default:
		out.WriteString("SS ")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

func (m *Man) HRule(out *bytes.Buffer) {
	out.WriteString("\n.ti 0\n\\l'\\n(.lu'\n")
}

func (m *Man) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	out.WriteString(".IP ")
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		out.WriteString("\\(bu 2")
	} else {
		out.WriteString("\\n+[step" + string(flags) + "]")
	}
	out.WriteString("\n")
	if !text() {
		out.Truncate(marker)
		return
	}

}

func (m *Man) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("\n\\item ")
	out.Write(text)
}

func (m *Man) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	out.WriteString("\n.TP\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

// TODO: This might now work
func (m *Man) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.WriteString(".TS\nallbox;\n")

	out.Write(header)
	out.Write(body)
	out.WriteString("\n.TE\n")
}

func (m *Man) TableRow(out *bytes.Buffer, text []byte) {
	if out.Len() > 0 {
		out.WriteString("\n")
	}
	out.Write(text)
	out.WriteString("\n")
}

func (m *Man) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" ")
	}
	out.Write(text)
	out.WriteString(" ")
}

// TODO: This is probably broken
func (m *Man) TableCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString("\t")
	}
	out.Write(text)
	out.WriteString("\t")
}

func (m *Man) Footnotes(out *bytes.Buffer, text func() bool) {

}

func (m *Man) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {

}

func (m *Man) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("\n\\[la]")
	out.Write(link)
	out.WriteString("\\[ra]")
}

func (m *Man) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("\\fB\\fC")
	escapeSpecialChars(out, text)
	out.WriteString("\\fR")
}

func (m *Man) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\fB")
	out.Write(text)
	out.WriteString("\\fP")
}

func (m *Man) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\fI")
	out.Write(text)
	out.WriteString("\\fP")
}

func (m *Man) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	fmt.Errorf("man: Image not supported")
}

func (m *Man) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n.br\n")
}

func (m *Man) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	m.AutoLink(out, link, 0)
}

func (m *Man) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	fmt.Errorf("man: Raw HTML not supported")
}

func (m *Man) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\s+2")
	out.Write(text)
	out.WriteString("\\s-2")
}

func (m *Man) StrikeThrough(out *bytes.Buffer, text []byte) {
	fmt.Errorf("man: strikethrough not supported")
}

func (m *Man) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {

}

func (m *Man) Entity(out *bytes.Buffer, entity []byte) {
	// TODO: convert this into a unicode character or something
	out.Write(entity)
}

func (m *Man) NormalText(out *bytes.Buffer, text []byte) {
	escapeSpecialChars(out, text)
}

// header and footer
func (m *Man) DocumentHeader(out *bytes.Buffer) {

}

func (m *Man) DocumentFooter(out *bytes.Buffer) {
}

func needsBackslash(c byte) bool {
	for _, r := range []byte("_{}%$&\\~") {
		if c == r {
			return true
		}
	}
	return false
}

func escapeSpecialChars(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		for i < len(text) && !needsBackslash(text[i]) {
			i++
		}
		if i > org {
			out.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		out.WriteByte('\\')
		out.WriteByte(text[i])
	}
}

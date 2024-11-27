package xml2go

import (
	"encoding/xml"
	"io"
	"log/slog"
)

func (b *builder) parse(f io.Reader) error {
	dec := xml.NewDecoder(f)
	for {
		tok, err := dec.Token()
		if tok == nil {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch elm := tok.(type) {
		case xml.StartElement:
			b.begin(elm)
		case xml.EndElement:
			b.end(elm)
		case xml.ProcInst:
		case xml.CharData:
			b.data(string(elm))
		default:
			slog.Warn("unhandled element", "elm", elm)
		}
	}
	return nil
}

type builder struct {
	StructsMap  map[string]Gostruct
	structStack *stack[Gostruct]
	fieldStack  *stack[Gofield]
	lastDATA    string
}

func NewBuilder() *builder {
	return &builder{
		StructsMap:  make(map[string]Gostruct),
		structStack: new(stack[Gostruct]),
		fieldStack:  new(stack[Gofield]),
	}
}

func (b *builder) data(text string) {
	b.lastDATA = text
}

func (b *builder) begin(elem xml.StartElement) {
	fieldAttrs := fieldAttributes(elem)
	if len(fieldAttrs) > 0 || b.structStack.empty() {
		s := newStruct(title(elem.Name.Local))
		for _, each := range fieldAttrs {
			s = s.withField(Gofield{Name: title(each.Name.Local), isAttr: true, XMLtag: each.Name.Local, Typ: "string"})
		}
		b.structStack.push(s)
		return
	}
	// no attr
	f := Gofield{Name: title(elem.Name.Local), XMLtag: elem.Name.Local, Typ: "string"}
	b.fieldStack.push(f)
}
func (b *builder) end(elem xml.EndElement) {
	// end of struct or end of field
	if !b.structStack.empty() {
		s := b.structStack.top()
		// closes top struct?
		if s.Name == title(elem.Name.Local) {
			b.StructsMap[s.Name] = s
			b.structStack.pop()
			return
		}
	}
	// closes top field?
	if !b.fieldStack.empty() {
		f := b.fieldStack.pop()
		if !b.structStack.empty() {
			s := b.structStack.pop()
			b.structStack.push(s.withField(f))
			return
		}
	}
	slog.Warn("unmatched", "end", elem.Name.Local)
}

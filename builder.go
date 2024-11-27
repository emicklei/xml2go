package xml2go

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"strings"
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

type Gostruct struct {
	Name   string
	Fields map[string]Gofield
}

func newStruct(name string) Gostruct {
	return Gostruct{
		Name:   name,
		Fields: make(map[string]Gofield),
	}
}

func (s Gostruct) addField(f Gofield) {
	// for now override
	s.Fields[f.Name] = f
}

func (s Gostruct) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "type %s struct {\n", s.Name)
	fmt.Fprintf(buf, "\tXMLName xml.Name `xml:\"%s\"`\n", s.Name)
	for _, each := range s.Fields {
		if each.isAttr {
			fmt.Fprintf(buf, "\t%s %s `xml:\"%s,attr\"`\n", each.Name, each.Typ, each.XMLtag)
		} else {
			fmt.Fprintf(buf, "\t%s %s `xml:\"%s\"`\n", each.Name, each.Typ, each.XMLtag)
		}
	}
	fmt.Fprintf(buf, "}\n")
	return buf.String()
}

type Gofield struct {
	Name   string
	Typ    string
	XMLtag string
	isAttr bool
}

type stack[T any] struct {
	elements []T
}

func (s *stack[T]) push(elem T) {
	s.elements = append(s.elements, elem)
}

func (s *stack[T]) pop() T {
	elem := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return elem
}

func (s *stack[T]) top() T {
	return s.elements[len(s.elements)-1]
}

func (s *stack[T]) empty() bool {
	return len(s.elements) == 0
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
			s.addField(Gofield{Name: title(each.Name.Local), isAttr: true, XMLtag: each.Name.Local, Typ: detectGoType(each.Value)})
		}
		b.structStack.push(s)
	}
	// no attr
	f := Gofield{Name: title(elem.Name.Local), XMLtag: elem.Name.Local}
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
			s := b.structStack.top()
			f.Typ = detectGoType(b.lastDATA)
			s.addField(f)
			return
		}
	}
	slog.Warn("unmatched end", "elem", elem.Name.Local)
}

func detectGoType(valueString string) string {
	// best guess
	switch valueString {
	case "true", "false":
		return "boolean"
	default:
		return "string"
	}
}

func title(name string) string {
	// not unicode safe
	return strings.Title(name)
}

func fieldAttributes(elem xml.StartElement) []xml.Attr {
	var attrs []xml.Attr
	for _, each := range elem.Attr {
		if each.Name.Space == "xmlns" {
			continue
		}
		attrs = append(attrs, each)
	}
	return attrs
}

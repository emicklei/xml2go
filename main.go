package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	f, _ := os.Open("./data/doc1.xml")
	dec := xml.NewDecoder(f)
	b := newBuilder()
	for {
		tok, err := dec.Token()
		if tok == nil {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("fail", "err", err.Error())
			break
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
			slog.Info("other", "elm", elm)
		}
	}
	for _, each := range b.structsMap {
		fmt.Println(each.String())
	}
}

type gostruct struct {
	name   string
	fields map[string]gofield
}

func newStruct(name string) gostruct {
	return gostruct{
		name:   name,
		fields: make(map[string]gofield),
	}
}

func (s gostruct) addField(f gofield) {
	// for now override
	s.fields[f.name] = f
}

func (s gostruct) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "type %s struct {\n", s.name)
	fmt.Fprintf(buf, "\tXMLName xml.Name `xml:\"%s\"`\n", s.name)
	for _, each := range s.fields {
		if each.isAttr {
			fmt.Fprintf(buf, "\t%s %s `xml:\"%s,attr\"`\n", each.name, each.typ, each.xmltag)
		} else {
			fmt.Fprintf(buf, "\t%s %s `xml:\"%s\"`\n", each.name, each.typ, each.xmltag)
		}
	}
	fmt.Fprintf(buf, "}\n")
	return buf.String()
}

type gofield struct {
	name   string
	typ    string
	xmltag string
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
	structsMap  map[string]gostruct
	structStack *stack[gostruct]
	fieldStack  *stack[gofield]
	lastDATA    string
}

func newBuilder() *builder {
	return &builder{
		structsMap:  make(map[string]gostruct),
		structStack: new(stack[gostruct]),
		fieldStack:  new(stack[gofield]),
	}
}

var titler = cases.Title(language.English)

func (b *builder) data(text string) {
	b.lastDATA = text
}

func (b *builder) begin(elem xml.StartElement) {
	slog.Info("begin", "elem", elem.Name.Local)
	if len(elem.Attr) > 0 {
		s := newStruct(titler.String(elem.Name.Local))
		for _, each := range elem.Attr {
			s.addField(gofield{name: titler.String(each.Name.Local), isAttr: true, xmltag: each.Name.Local, typ: detectGoType(each.Value)})
		}
		b.structStack.push(s)
	}
	// no attr
	f := gofield{name: titler.String(elem.Name.Local), xmltag: elem.Name.Local}
	b.fieldStack.push(f)
}
func (b *builder) end(elem xml.EndElement) {
	slog.Info("end", "elem", elem.Name.Local)
	// end of struct or end of field
	if !b.structStack.empty() {
		s := b.structStack.top()
		// closes top struct?
		if s.name == titler.String(elem.Name.Local) {
			b.structsMap[s.name] = s
			b.structStack.pop()
			return
		}
	}
	// closes top field?
	if !b.fieldStack.empty() {
		f := b.fieldStack.pop()
		if !b.structStack.empty() {
			s := b.structStack.top()
			f.typ = detectGoType(b.lastDATA)
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

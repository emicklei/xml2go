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
			slog.Info("inst", "elm", string(elm.Inst))
		case xml.CharData:
			slog.Info("data", "elm", string(elm))
		default:
			slog.Info("other", "elm", elm)
		}
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
		fmt.Fprintf(buf, "\t%s %s `xml:\"%s,attr\"`\n", each.name, each.typ, each.xmltag)
	}
	fmt.Fprintf(buf, "}\n")
	return buf.String()
}

type gofield struct {
	name   string
	typ    string
	xmltag string
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

type builder struct {
	structsMap  map[string]gostruct
	structStack *stack[gostruct]
	fieldStack  *stack[gofield]
}

func newBuilder() *builder {
	return &builder{
		structsMap:  make(map[string]gostruct),
		structStack: new(stack[gostruct]),
		fieldStack:  new(stack[gofield]),
	}
}

var titler = cases.Title(language.English)

func (b *builder) begin(elem xml.StartElement) {
	if len(elem.Attr) > 0 {
		s := newStruct(titler.String(elem.Name.Local))
		for _, each := range elem.Attr {
			s.addField(gofield{name: titler.String(each.Name.Local), typ: detectGoType(each.Value)})
		}
		b.structStack.push(s)
	}
	// no attr
	f := gofield{name: titler.String(elem.Name.Local)}
	b.structStack.top().addField(f)
}
func (b *builder) end(elem xml.EndElement) {

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

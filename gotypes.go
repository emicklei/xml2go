package xml2go

import (
	"fmt"
	"strings"
)

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

func (s Gostruct) withField(f Gofield) Gostruct {
	// for now override
	e, ok := s.Fields[f.Name]
	if !ok {
		s.Fields[f.Name] = f
		return s
	}
	// exists so make it repeated
	s.Fields[e.Name] = e.withSliceType()
	return s
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

func (f Gofield) withSliceType() Gofield {
	// already a slice?
	if strings.HasPrefix(f.Typ, "[]") {
		return f
	}
	f.Typ = "[]" + f.Typ
	return f
}

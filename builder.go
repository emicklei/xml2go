package xml2go

import (
	"encoding/xml"
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

type builder struct {
	StructsMap map[string]Gostruct
	nodeStack  *stack[node]
	lastDATA   string
}

func NewBuilder() *builder {
	return &builder{
		StructsMap: make(map[string]Gostruct),
		nodeStack:  new(stack[node]),
	}
}

func (b *builder) data(text string) {
	b.lastDATA = text
}

func (b *builder) begin(elem xml.StartElement) {
	n := node{
		name:   elem.Name.Local,
		attrs:  fieldAttributes(elem),
		xmltag: elem.Name.Local,
	}
	b.nodeStack.push(n)
}
func (b *builder) end(elem xml.EndElement) {
	n := b.nodeStack.pop()
	// closing top node?
	if n.name == elem.Name.Local {
		if b.nodeStack.empty() {
			s := b.makeStruct(n)
			b.StructsMap[s.Name] = s
		} else {
			// add to parent
			p := b.nodeStack.pop()
			p.nodes = append(p.nodes, n)
			b.nodeStack.push(p)
		}
	} else {
		// new child
		n.nodes = append(n.nodes, n)
		b.nodeStack.push(n)
	}
	// slog.Warn("unmatched", "end", elem.Name.Local)
}

func (b *builder) makeStruct(n node) Gostruct {
	s := newStruct(strings.Title(n.name))
	for _, each := range n.attrs {
		f := newField(strings.Title(each.Name.Local), each.Name.Local, true)
		s = s.withField(f)
	}
	for _, each := range n.nodes {
		f := newField(strings.Title(each.name), each.xmltag, false)
		f.Typ = strings.Title(each.name)
		s = s.withField(f)
	}
	return s
}

type node struct {
	name   string
	attrs  []xml.Attr
	xmltag string
	nodes  []node
}

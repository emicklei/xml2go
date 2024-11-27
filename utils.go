package xml2go

import (
	"encoding/xml"
	"strings"
)

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

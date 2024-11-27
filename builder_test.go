package xml2go

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	src := `<?xml version="1.0"?><typeA/>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(b.StructsMap))
	assert.Equal(t, "TypeA", b.StructsMap["TypeA"].Name)
}

func TestEmptyNS(t *testing.T) {
	src := `<?xml version="1.0"?><typeA xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(b.StructsMap))
	s := b.StructsMap["TypeA"]
	assert.Equal(t, "TypeA", s.Name)
	assert.Equal(t, 0, len(s.Fields))
}

func TestOneAttr(t *testing.T) {
	src := `<?xml version="1.0"?><typeA a="b"/>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "A", s.Fields["A"].Name)
	assert.Equal(t, "a", s.Fields["A"].XMLtag)
}

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
func TestOneField(t *testing.T) {
	src := `<?xml version="1.0"?><typeA><field>42</field></typeA>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "Field", s.Fields["Field"].Name)
	assert.Equal(t, "field", s.Fields["Field"].XMLtag)
}

func TestOneFieldNS(t *testing.T) {
	src := `<?xml version="1.0"?><typeA><field xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">42</field></typeA>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "Field", s.Fields["Field"].Name)
	assert.Equal(t, "field", s.Fields["Field"].XMLtag)
}

func TestTwoFields(t *testing.T) {
	src := `<?xml version="1.0"?><typeA><field1>42</field1><field2>true</field2></typeA>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 2, len(s.Fields))
	assert.Equal(t, "Field1", s.Fields["Field1"].Name)
}
func TestTwoRepeatedFields(t *testing.T) {
	src := `<?xml version="1.0"?><typeA><field1>42</field1><field1>12</field1></typeA>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "[]string", s.Fields["Field1"].Typ)
}
func TestOneNestedType(t *testing.T) {
	src := `<?xml version="1.0"?>
	<typeA>
		<typeB>
			<field1>42</field1>
		</typeB>
	</typeA>`
	b := NewBuilder()
	if err := b.parse(strings.NewReader(src)); err != nil {
		t.Error(err)
	}
	s := b.StructsMap["TypeA"]
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "TypeB", s.Fields["TypeB"].Typ)
}

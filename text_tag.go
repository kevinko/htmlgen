// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"html"
	"io"
)

// A pseudo tag that holds text.
type TextTag struct {
	parent *baseTag

	text string
}

func (t *TextTag) A(href string, text string) *TextTag {
	t.parent.A(href).T(text)
	return t.parent.T()
}

func (t *TextTag) Abbr(text string, title ...string) *TextTag {
	t.parent.Abbr(title...).T(text)
	return t.parent.T()
}

// Assigns t to tag and returns t.
func (t *TextTag) Assign(tag **TextTag) *TextTag {
	*tag = t
	return t
}

func (t *TextTag) B(text string) *TextTag {
	t.parent.B().T(text)
	return t.parent.T()
}

// Like B() but assigns the nested text tag to tag.
func (t *TextTag) BAssign(text string, tag **TextTag) *TextTag {
	t.parent.B().T(text).Assign(tag)
	return t.parent.T()
}

func (t *TextTag) Cite(text string) *TextTag {
	t.parent.Cite().T(text)
	return t.parent.T()
}

func (t *TextTag) Code(text string) *TextTag {
	t.parent.Code().T(text)
	return t.parent.T()
}

// Returns a copy of the TextTag.
func (t *TextTag) Copy() *TextTag {
	return &TextTag{text: t.text}
}

func (t *TextTag) Dfn(text string) *TextTag {
	t.parent.Dfn().T(text)
	return t.parent.T()
}

func (t *TextTag) Em(text string) *TextTag {
	t.parent.Em().T(text)
	return t.parent.T()
}

func (t *TextTag) I(text string) *TextTag {
	t.parent.I().T(text)
	return t.parent.T()
}

func (t *TextTag) Img(src, alt string, options ...*ImgOptions) *TextTag {
	t.parent.Img(src, alt, options...)
	return t.parent.T()
}

func (t *TextTag) isHidden() bool {
	return false
}

func (t *TextTag) Kbd(text string) *TextTag {
	t.parent.Kbd().T(text)
	return t.parent.T()
}

// Returns the parent tag for t.
func (t *TextTag) Parent() Tag {
	return t.parent
}

func (t *TextTag) Pre(text string) *TextTag {
	t.parent.Pre().T(text)
	return t.parent.T()
}

func (t *TextTag) Samp(text string) *TextTag {
	t.parent.Samp().T(text)
	return t.parent.T()
}

// Assigns text to the TextTag, replacing existing text.
func (t *TextTag) SetText(text string) *TextTag {
	t.text = text
	return t
}

func (t *TextTag) Small(text string) *TextTag {
	t.parent.Small().T(text)
	return t.parent.T()
}

func (t *TextTag) Strong(text string) *TextTag {
	t.parent.Strong().T(text)
	return t.parent.T()
}

func (t *TextTag) T(text ...string) *TextTag {
	if len(text) == 0 {
		return t
	}
	t.text += html.EscapeString(text[0])
	return t
}

// Returns the TextTag's text.
func (t *TextTag) Text() string {
	return t.text
}

func (t *TextTag) TUnsafe(text ...string) *TextTag {
	if len(text) == 0 {
		return t
	}
	t.text += text[0]
	return t
}

func (t *TextTag) TV(text ...string) *TextTagVar {
	return t.parent.TV(text...)
}

func (t *TextTag) TVUnsafe(text ...string) *TextTagVar {
	return t.parent.TVUnsafe(text...)
}

func (t *TextTag) U(text string) *TextTag {
	t.parent.U().T(text)
	return t.parent.T()
}

func (t *TextTag) Up(count ...int) Tag {
	realCount := 1
	if len(count) > 0 {
		realCount = count[0]
	}

	var currTag Tag
	currTag = t.parent
	for realCount > 1 {
		currTag = currTag.Parent()

		realCount--
	}
	return currTag
}

func (t *TextTag) Var(text string) *TextTag {
	t.parent.Var().T(text)
	return t.parent.T()
}

func (t *TextTag) write(writer io.Writer, env ...Environment) (int, error) {
	return io.WriteString(writer, t.text)
}

func (t *TextTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	if count, err := writeIndent(writer, indent); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := t.write(writer); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

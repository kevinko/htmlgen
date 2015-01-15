// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"html"
	"io"
	"regexp"
	"strings"
)

const (
	undefinedVarValue = StringValue("undefined")
)

var (
	reVar = regexp.MustCompile(`^(\$|[a-zA-Z_]+[a-zA-Z0-9_]*)`)
)

type Environment map[string]Value

type StringValue string

func (s StringValue) String() string {
	return string(s)
}

type Value interface {
	String() string
}

// Specifies a TextTag containing variables that can be expanded at render time.
//
// variables are of the form:
//
//  ^\$[a-zA-Z0-9_]+$
//
// $$ becomes $.
// $ will be deleted.
type TextTagVar struct {
	parent *baseTag

	text string

	vars []Var

	// Unsafe text strings will not be HTML-escaped.
	isUnsafe bool
}

func NewTextTagVar(text string) *TextTagVar {
	return &TextTagVar{
		text: text,
		vars: parseVars(text),
	}
}

func (t *TextTagVar) A(href string, text string) *TextTagVar {
	t.parent.A(href).TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Abbr(text string, title ...string) *TextTagVar {
	t.parent.Abbr(title...).TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) B(text string) *TextTagVar {
	t.parent.B().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Cite(text string) *TextTagVar {
	t.parent.Cite().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Code(text string) *TextTagVar {
	t.parent.Code().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Dfn(text string) *TextTagVar {
	t.parent.Dfn().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Em(text string) *TextTagVar {
	t.parent.Em().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) I(text string) *TextTagVar {
	t.parent.I().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Img(src, alt string, options ...*ImgOptions) *TextTagVar {
	t.parent.Img(src, alt, options...)
	return t.parent.TV()
}

func (t *TextTagVar) isHidden() bool {
	return false
}

func (t *TextTagVar) Kbd(text string) *TextTagVar {
	t.parent.Kbd().TV(text)
	return t.parent.TV()
}

// Returns the parent tag for t.
func (t *TextTagVar) Parent() Tag {
	return t.parent
}

func (t *TextTagVar) Pre(text string) *TextTagVar {
	t.parent.Pre().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Samp(text string) *TextTagVar {
	t.parent.Samp().TV(text)
	return t.parent.TV()
}

// Assigns text to the TextTag, replacing existing text.
func (t *TextTagVar) SetText(text string) *TextTagVar {
	t.text = text
	return t
}

func (t *TextTagVar) SetUnsafe(isUnsafe bool) *TextTagVar {
	t.isUnsafe = isUnsafe
	return t
}

func (t *TextTagVar) Small(text string) *TextTagVar {
	t.parent.Small().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Strong(text string) *TextTagVar {
	t.parent.Strong().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) T(text ...string) *TextTag {
	return t.parent.T(text...)
}

func (t *TextTagVar) TUnsafe(text ...string) *TextTag {
	return t.parent.TUnsafe(text...)
}

func (t *TextTagVar) TV(text ...string) *TextTagVar {
	if len(text) == 0 {
		return t
	}
	t.text += text[0]
	return t
}

// Returns the TextTagVar's text.
func (t *TextTagVar) Text() string {
	return t.text
}

func (t *TextTagVar) U(text string) *TextTagVar {
	t.parent.U().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) Var(text string) *TextTagVar {
	t.parent.Var().TV(text)
	return t.parent.TV()
}

func (t *TextTagVar) write(writer io.Writer, optEnv ...Environment) (count int, err error) {
	var env Environment
	if len(optEnv) > 0 {
		env = optEnv[0]
	} else {
		env = make(map[string]Value)
	}

	offset := 0
	for _, v := range t.vars {
		text := t.text[offset:v.startOffset]
		if !t.isUnsafe {
			text = html.EscapeString(text)
		}

		if n, writeErr := io.WriteString(writer, text); writeErr != nil {
			err = writeErr
			return
		} else {
			count += n
		}

		varName := t.text[v.startOffset : v.startOffset+v.length]

		// Strip out the leading variable delimiter '$'.
		varName = varName[1:]

		// Look for an escaped '$', which is written as "$$".
		if varName == "$" {
			if n, writeErr := io.WriteString(writer, "$"); writeErr != nil {
				err = writeErr
				return
			} else {
				count += n
			}
		} else {
			value, ok := env[varName]
			if !ok {
				value = undefinedVarValue
			}

			valueStr := value.String()
			if !t.isUnsafe {
				valueStr = html.EscapeString(valueStr)
			}

			if n, writeErr := io.WriteString(writer, valueStr); writeErr != nil {
				err = writeErr
				return
			} else {
				count += n
			}
		}

		offset = v.startOffset + v.length
	}

	// Write the remainder.
	text := t.text[offset:]
	if !t.isUnsafe {
		text = html.EscapeString(text)
	}

	n, err := io.WriteString(writer, text)
	if err != nil {
		return
	}

	count += n
	return
}

func (t *TextTagVar) writePretty(writer io.Writer, indent int, env ...Environment) (int, error) {
	count := 0
	if n, err := writeIndent(writer, indent); err != nil {
		return count, err
	} else {
		count += n
	}
	if n, err := t.write(writer, env...); err != nil {
		return count, err
	} else {
		count += n
	}
	return count, nil
}

func parseVars(text string) (vars []Var) {
	offset := 0
	for offset < len(text) {
		currText := text[offset:]

		index := strings.IndexRune(currText, '$')
		if index == -1 {
			// We're done.
			return
		}

		match := reVar.FindString(currText[index+1:])
		// Include the leading '$'.
		varLen := 1 + len(match)
		vars = append(vars, Var{startOffset: offset + index, length: varLen})

		offset += index + varLen
	}

	return
}

type Var struct {
	// The position of the '$' character.
	startOffset int

	// The length of the full variable name, including the initial '$'.
	length int
}

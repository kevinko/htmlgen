// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// The base fields for all tags.
type baseTag struct {
	htmlGen

	tagType int

	// Tag attributes.
	attrs map[int]string

	customAttrs map[string]string

	children []tagWriter
	parent   Tag

	// By default (false), a cache is dirty.
	isCacheClean bool
	// Caches the rendered open tag, including attributes.
	cacheOpen string

	// Hides the tag during the rendering process if true.
	hidden bool
}

func newBaseTag(tagType int) *baseTag {
	return &baseTag{
		tagType:     tagType,
		attrs:       make(map[int]string),
		customAttrs: make(map[string]string),
		children:    make([]tagWriter, 0),
	}
}

// Adds newTag as a child of t and returns t.
func (t *baseTag) AddChild(newTag Tag) Tag {
	t.children = append(t.children, newTag)
	return t
}

func (t *baseTag) AddChildText(tag *TextTag) Tag {
	t.children = append(t.children, tag)
	return t
}

func (t *baseTag) AddClass(classes ...string) Tag {
	return t.AddClasses(classes)
}

func (t *baseTag) AddClasses(classes []string) Tag {
	if classes == nil || len(classes) == 0 {
		return t
	} else {
		classesStr := strings.Join(classes, " ")

		oldClassesStr := t.attrs[kAttrClass]
		if len(oldClassesStr) > 0 {
			t.attrs[kAttrClass] += " " + classesStr
		} else {
			t.attrs[kAttrClass] = classesStr
		}
	}

	t.isCacheClean = false
	return t
}

func (t *baseTag) Assign(tag *Tag) Tag {
	*tag = t
	return t
}

func (t *baseTag) A(href string) Tag {
	return addChild(t, t.htmlGen.A(href))
}

func (t *baseTag) Abbr(title ...string) Tag {
	return addChild(t, t.htmlGen.Abbr(title...))
}

func (t *baseTag) Address() Tag {
	return addChild(t, t.htmlGen.Address())
}

func (t *baseTag) B() Tag {
	return addChild(t, t.htmlGen.B())
}

func (t *baseTag) Blockquote() Tag {
	return addChild(t, t.htmlGen.Blockquote())
}

func (t *baseTag) Body() *BodyTag {
	newTag := t.htmlGen.Body()
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) Br() Tag {
	return addChild(t, t.htmlGen.Br())
}

func (t *baseTag) Button() Tag {
	return addChild(t, t.htmlGen.Button())
}

func (t *baseTag) Canvas(options ...*CanvasOptions) Tag {
	return addChild(t, t.htmlGen.Canvas(options...))
}

func (t *baseTag) Caption() Tag {
	return addChild(t, t.htmlGen.Caption())
}

func (t *baseTag) CheckedInput(inputType CheckedInputType, options ...*InputOptions) *CheckedInputTag {
	newTag := t.htmlGen.CheckedInput(inputType, options...)
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) CheckedInputTypeNameValue(inputType CheckedInputType, name, value string) *CheckedInputTag {
	newTag := t.htmlGen.CheckedInputTypeNameValue(inputType, name, value)
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) Cite() Tag {
	return addChild(t, t.htmlGen.Cite())
}

func (t *baseTag) Code() Tag {
	return addChild(t, t.htmlGen.Code())
}

func (t *baseTag) Comment() Tag {
	return addChild(t, t.htmlGen.Comment())
}

func (t *baseTag) Copy() Tag {
	if !t.isCacheClean {
		tagStr := tagTypeStringMap[t.tagType]
		var err error
		if t.cacheOpen, err = t.renderCacheOpen(tagStr); err == nil {
			t.isCacheClean = true
		}
		// Otherwise, leave the cache dirty.
	}

	return &baseTag{
		tagType:      t.tagType,
		attrs:        copyAttrs(t.attrs),
		customAttrs:  copyCustomAttrs(t.customAttrs),
		children:     make([]tagWriter, 0),
		isCacheClean: t.isCacheClean,
		cacheOpen:    t.cacheOpen,
	}
}

func (t *baseTag) Datalist() Tag {
	return addChild(t, t.htmlGen.Datalist())
}

func (t *baseTag) Dfn() Tag {
	return addChild(t, t.htmlGen.Dfn())
}

func (t *baseTag) Div() Tag {
	return addChild(t, t.htmlGen.Div())
}

func (t *baseTag) DivClasses(classes ...string) Tag {
	return addChild(t, t.htmlGen.DivClasses(classes...))
}

func (t *baseTag) DivId(id string) Tag {
	return addChild(t, t.htmlGen.DivId(id))
}

func (t *baseTag) DivIdClasses(id string, classes ...string) Tag {
	return addChild(t, t.htmlGen.DivIdClasses(id, classes...))
}

func (t *baseTag) Dd() Tag {
	return addChild(t, t.htmlGen.Dd())
}

func (t *baseTag) Dl() Tag {
	return addChild(t, t.htmlGen.Dl())
}

func (t *baseTag) Dt() Tag {
	return addChild(t, t.htmlGen.Dt())
}

func (t *baseTag) Em() Tag {
	return addChild(t, t.htmlGen.Em())
}

func (t *baseTag) Footer() Tag {
	return addChild(t, t.htmlGen.Footer())
}

func (t *baseTag) Form(options ...*FormOptions) Tag {
	return addChild(t, t.htmlGen.Form(options...))
}

func (t *baseTag) getChildren() []tagWriter {
	if t.children == nil {
		return []tagWriter{}
	}
	return t.children
}

func (t *baseTag) H1() Tag {
	return addChild(t, t.htmlGen.H1())
}

func (t *baseTag) H2() Tag {
	return addChild(t, t.htmlGen.H2())
}

func (t *baseTag) H3() Tag {
	return addChild(t, t.htmlGen.H3())
}

func (t *baseTag) H4() Tag {
	return addChild(t, t.htmlGen.H4())
}

func (t *baseTag) H5() Tag {
	return addChild(t, t.htmlGen.H5())
}

func (t *baseTag) H6() Tag {
	return addChild(t, t.htmlGen.H6())
}

func (t *baseTag) Head() Tag {
	return addChild(t, t.htmlGen.Head())
}

func (t *baseTag) Hide(isHidden bool) {
	t.hidden = isHidden
}

func (t *baseTag) Hr() Tag {
	return addChild(t, t.htmlGen.Hr())
}

func (t *baseTag) I() Tag {
	return addChild(t, t.htmlGen.I())
}

func (t *baseTag) Id() string {
	// "" returned by default.
	return t.attrs[kAttrId]
}

func (t *baseTag) Img(src, alt string, options ...*ImgOptions) Tag {
	return addChild(t, t.htmlGen.Img(src, alt, options...))
}

func (t *baseTag) Input(inputType InputType, options ...*InputOptions) *InputTag {
	newTag := t.htmlGen.Input(inputType, options...)
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) InputTypeNameValue(inputType InputType, name, value string) *InputTag {
	newTag := t.htmlGen.InputTypeNameValue(inputType, name, value)
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) isHidden() bool {
	return t.hidden
}

func (t *baseTag) Kbd() Tag {
	return addChild(t, t.htmlGen.Kbd())
}

func (t *baseTag) Label(options ...*LabelOptions) Tag {
	return addChild(t, t.htmlGen.Label(options...))
}

func (t *baseTag) Li() Tag {
	return addChild(t, t.htmlGen.Li())
}

func (t *baseTag) Link(rel string, options ...*LinkOptions) Tag {
	return addChild(t, t.htmlGen.Link(rel, options...))
}

func (t *baseTag) Meta(name, content string, options ...*MetaOptions) Tag {
	return addChild(t, t.htmlGen.Meta(name, content, options...))
}

func (t *baseTag) NoScript() Tag {
	return addChild(t, t.htmlGen.NoScript())
}

func (t *baseTag) Ol() Tag {
	return addChild(t, t.htmlGen.Ol())
}

func (t *baseTag) Option(options ...*OptionOptions) *OptionTag {
	newTag := t.htmlGen.Option(options...)
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) P() Tag {
	return addChild(t, t.htmlGen.P())
}

func (t *baseTag) Parent() Tag {
	return t.parent
}

func (t *baseTag) Pre() Tag {
	return addChild(t, t.htmlGen.Pre())
}

func (t *baseTag) RemoveAttribute(key string) Tag {
	delete(t.customAttrs, key)
	t.isCacheClean = false
	return t
}

func (t *baseTag) RemoveAttributes(keys ...string) Tag {
	for _, key := range keys {
		delete(t.customAttrs, key)
	}
	t.isCacheClean = false
	return t
}

func (t *baseTag) RemoveChildren() Tag {
	t.children = t.children[:0]
	return t
}

func (t *baseTag) RemoveParent() Tag {
	// Remove t from parent's children slice.
	parent := t.parent
	parentChildren := parent.getChildren()

	// Search in reverse order to optimize for the most recently added
	// child.
	for ii := len(parentChildren) - 1; ii >= 0; ii-- {
		if parentChildren[ii] != t {
			continue
		}

		// Splice the sections before and after ii.
		newChildren := parentChildren[:ii]
		newChildren = append(newChildren, parentChildren[ii+1:]...)
		parent.setChildren(newChildren)
		break
	}

	t.parent = nil

	return t
}

// Renders the opening tag with attributes to a string.
func (t *baseTag) renderCacheOpen(tagStr string) (result string, err error) {
	tmp := new(bytes.Buffer)
	if _, err = t.writeOpenTagLead(tmp, tagStr); err != nil {
		return
	}
	// Finish the open tag.
	if _, err = writeRune(tmp, '>'); err != nil {
		return
	}
	result = tmp.String()
	return
}

func (t *baseTag) Samp() Tag {
	return addChild(t, t.htmlGen.Samp())
}

func (t *baseTag) Script(optScriptType ...string) Tag {
	return addChild(t, t.htmlGen.Script(optScriptType...))
}

func (t *baseTag) ScriptSrc(scriptType, src string) Tag {
	return addChild(t, t.htmlGen.ScriptSrc(scriptType, src))
}

func (t *baseTag) Select(options ...*SelectOptions) *SelectTag {
	newTag := t.htmlGen.Select(options...)
	addChild(t, newTag)
	return newTag
}

// Sets the attribute specified by keyId to value.  If value is empty, the
// attribute will be cleared.  This returns t.
func (t *baseTag) setAttr(keyId int, value string) Tag {
	if len(value) == 0 {
		delete(t.attrs, keyId)
	} else {
		t.attrs[keyId] = value
	}

	t.isCacheClean = false
	return t
}

func (t *baseTag) SetAttribute(key, value string) Tag {
	t.customAttrs[key] = value
	t.isCacheClean = false
	return t
}

func (t *baseTag) SetAttributes(attrs map[string]string) Tag {
	for k, v := range attrs {
		t.customAttrs[k] = v
	}
	t.isCacheClean = false
	return t
}

func (t *baseTag) setChildren(children []tagWriter) {
	t.children = children
}

func (t *baseTag) SetClass(classes ...string) Tag {
	t.SetClasses(classes)
	return t
}

func (t *baseTag) SetClasses(classes []string) Tag {
	if classes == nil || len(classes) == 0 {
		delete(t.attrs, kAttrClass)
	} else {
		classesStr := strings.Join(classes, " ")
		t.attrs[kAttrClass] = classesStr
	}

	t.isCacheClean = false
	return t
}

func (t *baseTag) SetId(id string) Tag {
	return t.setAttr(kAttrId, id)
}

func (t *baseTag) SetOnblur(script string) Tag {
	return t.setAttr(kAttrOnblur, script)
}

func (t *baseTag) SetOnclick(script string) Tag {
	return t.setAttr(kAttrOnclick, script)
}

func (t *baseTag) SetOndblclick(script string) Tag {
	return t.setAttr(kAttrOndblclick, script)
}

func (t *baseTag) SetOnchange(script string) Tag {
	return t.setAttr(kAttrOnchange, script)
}

func (t *baseTag) SetOnfocus(script string) Tag {
	return t.setAttr(kAttrOnfocus, script)
}

func (t *baseTag) SetOnselect(script string) Tag {
	return t.setAttr(kAttrOnselect, script)
}

func (t *baseTag) SetOnsubmit(script string) Tag {
	return t.setAttr(kAttrOnsubmit, script)
}

func (t *baseTag) SetOnkeydown(script string) Tag {
	return t.setAttr(kAttrOnkeydown, script)
}

func (t *baseTag) SetOnkeypress(script string) Tag {
	return t.setAttr(kAttrOnkeypress, script)
}

func (t *baseTag) SetOnkeyup(script string) Tag {
	return t.setAttr(kAttrOnkeyup, script)
}

func (t *baseTag) SetOnmousedown(script string) Tag {
	return t.setAttr(kAttrOnmousedown, script)
}

func (t *baseTag) SetOnmousemove(script string) Tag {
	return t.setAttr(kAttrOnmousemove, script)
}

func (t *baseTag) SetOnmouseout(script string) Tag {
	return t.setAttr(kAttrOnmouseout, script)
}

func (t *baseTag) SetOnmouseover(script string) Tag {
	return t.setAttr(kAttrOnmouseover, script)
}

func (t *baseTag) SetOnmouseup(script string) Tag {
	return t.setAttr(kAttrOnmouseup, script)
}

func (t *baseTag) setParent(parent Tag) {
	t.parent = parent
}

func (t *baseTag) SetTitle(title string) Tag {
	return t.setAttr(kAttrTitle, title)
}

func (t *baseTag) Span() Tag {
	return addChild(t, t.htmlGen.Span())
}

func (t *baseTag) SpanClasses(classes ...string) Tag {
	return addChild(t, t.htmlGen.SpanClasses(classes...))
}

func (t *baseTag) SpanId(id string) Tag {
	return addChild(t, t.htmlGen.SpanId(id))
}

func (t *baseTag) SpanIdClasses(id string, classes ...string) Tag {
	return addChild(t, t.htmlGen.SpanIdClasses(id, classes...))
}

func (t *baseTag) Small() Tag {
	return addChild(t, t.htmlGen.Small())
}

func (t *baseTag) Strong() Tag {
	return addChild(t, t.htmlGen.Strong())
}

func (t *baseTag) Style(options ...*StyleOptions) Tag {
	return addChild(t, t.htmlGen.Style(options...))
}

func (t *baseTag) T(text ...string) *TextTag {
	newTag := t.htmlGen.T(text...)
	newTag.parent = t
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) Table(options ...*TableOptions) Tag {
	return addChild(t, t.htmlGen.Table(options...))
}

func (t *baseTag) Tbody() Tag {
	return addChild(t, t.htmlGen.Tbody())
}

func (t *baseTag) Td(options ...*TdOptions) Tag {
	return addChild(t, t.htmlGen.Td(options...))
}

func (t *baseTag) Textarea(rows, cols int, options ...*TextareaOptions) Tag {
	return addChild(t, t.htmlGen.Textarea(rows, cols, options...))
}

func (t *baseTag) Tfoot() Tag {
	return addChild(t, t.htmlGen.Tfoot())
}

func (t *baseTag) Th(options ...*ThOptions) Tag {
	return addChild(t, t.htmlGen.Th(options...))
}

func (t *baseTag) Thead() Tag {
	return addChild(t, t.htmlGen.Thead())
}

func (t *baseTag) Title() Tag {
	return addChild(t, t.htmlGen.Title())
}

func (t *baseTag) Tr() Tag {
	return addChild(t, t.htmlGen.Tr())
}

func (t *baseTag) TUnsafe(text ...string) *TextTag {
	newTag := t.htmlGen.TUnsafe(text...)
	newTag.parent = t
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) TV(text ...string) *TextTagVar {
	newTag := t.htmlGen.TV(text...)
	newTag.parent = t
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) TVUnsafe(text ...string) *TextTagVar {
	newTag := t.htmlGen.TVUnsafe(text...)
	newTag.parent = t
	t.children = append(t.children, newTag)
	return newTag
}

func (t *baseTag) U() Tag {
	return addChild(t, t.htmlGen.U())
}

func (t *baseTag) Ul() Tag {
	return addChild(t, t.htmlGen.Ul())
}

func (t *baseTag) Up(count ...int) Tag {
	realCount := 1
	if len(count) > 0 {
		realCount = count[0]
	}

	// Behavior always assumes >= 1 ancestor steps.
	currTag := t.parent
	for realCount > 1 {
		currTag = currTag.Parent()

		realCount--
	}
	return currTag
}

func (t *baseTag) Var() Tag {
	return addChild(t, t.htmlGen.Var())
}

func (t *baseTag) write(writer io.Writer, env ...Environment) (n int, err error) {
	if t.hidden {
		return
	}

	tagStr := tagTypeStringMap[t.tagType]

	// Write the opening tag.
	if !t.isCacheClean {
		if t.cacheOpen, err = t.renderCacheOpen(tagStr); err != nil {
			return
		}
		t.isCacheClean = true
	}
	if count, err := io.WriteString(writer, t.cacheOpen); err != nil {
		return n, err
	} else {
		n += count
	}

	// Recursively write all children.
	for _, childTag := range t.children {
		if count, err := childTag.write(writer, env...); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	// Write the closing tag.
	if count, err := io.WriteString(writer, "</"); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := io.WriteString(writer, tagStr); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := writeRune(writer, '>'); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

// NOTE: currently a newline is introduced if a child tag is hidden.
func (t *baseTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	if t.hidden {
		return
	}

	tagStr := tagTypeStringMap[t.tagType]

	// Pretty print with indentation.
	if count, err := writeIndent(writer, indent); err != nil {
		return n, err
	} else {
		n += count
	}

	// Write the opening tag.
	if count, err := t.writeOpenTagLeadSorted(writer, tagStr); err != nil {
		return n, err
	} else {
		n += count
	}
	// Finish the open tag.
	if count, err := writeRune(writer, '>'); err != nil {
		return n, err
	} else {
		n += count
	}

	newIndent := indent + kIndentSpace

	// Recursively write all children.
	for _, childTag := range t.children {
		if childTag.isHidden() {
			continue
		}

		// Pretty print with newline between children.
		if count, err := writeRune(writer, '\n'); err != nil {
			return n, err
		} else {
			n += count
		}

		if count, err := childTag.writePretty(writer, newIndent, env...); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	if len(t.children) > 0 {
		// Separate from children.
		if count, err := writeRune(writer, '\n'); err != nil {
			return n, err
		} else {
			n += count
		}
		if count, err := writeIndent(writer, indent); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	// Write the closing tag.
	s := fmt.Sprintf("</%s>", tagStr)
	if count, err := io.WriteString(writer, s); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

// Writes the leading part of the opening tag and its attributes, up until
// the closing ">".
func (t *baseTag) writeOpenTagLead(writer io.Writer, tagStr string) (n int, err error) {
	if count, err := writeRune(writer, '<'); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := io.WriteString(writer, tagStr); err != nil {
		return n, err
	} else {
		n += count
	}

	// Write any attributes.
	for key, value := range t.attrs {
		// key is a kAttr constant.
		if count, err := writeRune(writer, ' '); err != nil {
			return n, err
		} else {
			n += count
		}
		if count, err := writeKeyIdValue(writer, key, value); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	// Write custom attributes.
	for key, value := range t.customAttrs {
		if count, err := writeRune(writer, ' '); err != nil {
			return n, err
		} else {
			n += count
		}
		if count, err := writeKeyValue(writer, key, value); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	return
}

// Like writeOpenTagLead, only that attributes will be sorted.
func (t *baseTag) writeOpenTagLeadSorted(writer io.Writer, tagStr string) (n int, err error) {
	if count, err := writeRune(writer, '<'); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := io.WriteString(writer, tagStr); err != nil {
		return n, err
	} else {
		n += count
	}

	// The common case doesn't include custom attributes.
	if len(t.customAttrs) == 0 {
		sortedKeys := make([]int, len(t.attrs))
		ii := 0
		for k := range t.attrs {
			sortedKeys[ii] = k
			ii++
		}
		sort.Ints(sortedKeys)

		for _, key := range sortedKeys {
			value := t.attrs[key]
			if count, err := writeRune(writer, ' '); err != nil {
				return n, err
			} else {
				n += count
			}
			if count, err := writeKeyIdValue(writer, key, value); err != nil {
				return n, err
			} else {
				n += count
			}
		}
		return
	}

	// Otherwise, merge known and custom attrs into a single map.
	numKeys := len(t.attrs) + len(t.customAttrs)
	sortedKeys := make([]string, 0, numKeys)

	attrs := make(map[string]string, numKeys)
	for keyId, v := range t.attrs {
		key := attrStringMap[keyId]
		attrs[key] = v

		sortedKeys = append(sortedKeys, key)
	}
	for k, v := range t.customAttrs {
		attrs[k] = v
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		value := attrs[key]

		if count, err := writeRune(writer, ' '); err != nil {
			return n, err
		} else {
			n += count
		}
		if count, err := writeKeyValue(writer, key, value); err != nil {
			return n, err
		} else {
			n += count
		}
	}
	return
}

// Returns newTag.
func addChild(t *baseTag, newTag Tag) Tag {
	t.children = append(t.children, newTag)
	newTag.setParent(t)
	return newTag
}

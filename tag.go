// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"bytes"
	"io"
	"strconv"
)

type BodyTag struct {
	baseTag
}

func (t *BodyTag) SetOnload(script string) Tag {
	return t.setAttr(kAttrOnload, script)
}

type CheckedInputTag struct {
	*InputTag
}

func (t *CheckedInputTag) Checked() bool {
	_, ok := t.attrs[kAttrChecked]
	return ok
}

func (t *CheckedInputTag) SetChecked(checked bool) {
	if checked {
		t.attrs[kAttrChecked] = ""
	} else {
		delete(t.attrs, kAttrChecked)
	}
	t.isCacheClean = false
}

func (tag *CheckedInputTag) Type() CheckedInputType {
	return CheckedInputType(tag.InputTag.Type())
}

type commentTag struct {
	baseTag
}

func (t *commentTag) Copy() Tag {
	return &commentTag{}
}

func (t *commentTag) write(writer io.Writer, env ...Environment) (n int, err error) {
	if count, err := io.WriteString(writer, "<!-- "); err != nil {
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

	if count, err := io.WriteString(writer, " -->"); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

func (t *commentTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	if count, err := writeIndent(writer, indent); err != nil {
		return n, err
	} else {
		n += count
	}

	if count, err := io.WriteString(writer, "<!--"); err != nil {
		return n, err
	} else {
		n += count
	}

	newIndent := indent + kIndentSpace

	// Recursively write all children.
	for _, childTag := range t.children {
		// Pretty print with newline.
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
	} else {
		// Ensure that we have a space.
		if count, err := writeRune(writer, ' '); err != nil {
			return n, err
		} else {
			n += count
		}
	}

	if count, err := io.WriteString(writer, "-->"); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

// An htmlTag will write a DOCTYPE declaration as well.
type htmlTag struct {
	baseTag
}

func (t *htmlTag) write(writer io.Writer, env ...Environment) (n int, err error) {
	// Always prepend a DOCTYPE declaration.
	if count, err := io.WriteString(writer, "<!DOCTYPE html>"); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := t.baseTag.write(writer, env...); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

func (t *htmlTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	if count, err := writeIndent(writer, indent); err != nil {
		return n, err
	} else {
		n += count
	}
	// Always prepend a DOCTYPE declaration.
	if count, err := io.WriteString(writer, "<!DOCTYPE html>\n"); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := t.baseTag.writePretty(writer, indent, env...); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

type InputTag struct {
	singleTag
}

// Clears all options except for name and type.
func (t *InputTag) ResetOptions() {
	delete(t.attrs, kAttrAction)
	delete(t.attrs, kAttrAlt)
	delete(t.attrs, kAttrAutocomplete)
	delete(t.attrs, kAttrChecked)
	delete(t.attrs, kAttrDisabled)
	delete(t.attrs, kAttrHeight)
	delete(t.attrs, kAttrMaxlength)
	delete(t.attrs, kAttrReadonly)
	delete(t.attrs, kAttrSize)
	delete(t.attrs, kAttrSrc)
	delete(t.attrs, kAttrValue)
	delete(t.attrs, kAttrWidth)
}

// NOTE: existing options will not be reset unless explicitly assigned.
func (t *InputTag) SetOptions(o *InputOptions) {
	if len(o.Action) > 0 {
		t.attrs[kAttrAction] = o.Action
	}

	if len(o.Alt) > 0 {
		t.attrs[kAttrAlt] = o.Alt
	}

	if o.AutocompleteOff {
		t.attrs[kAttrAutocomplete] = kHtmlOff
	}

	if o.Checked {
		// NOTE: this does not turn off a checked attribute, since
		// the default zero value for InputOption's Checked field
		// is false.  Instead, use the SetChecked() setter.
		t.attrs[kAttrChecked] = ""
	}

	if o.Disabled {
		t.attrs[kAttrDisabled] = ""
	}

	if o.Height > 0 {
		t.attrs[kAttrHeight] = strconv.Itoa(o.Height)
	}

	if len(o.List) > 0 {
		t.attrs[kAttrList] = o.List
	}

	if o.Max != 0 {
		t.attrs[kAttrMax] = strconv.FormatFloat(o.Max, 'g', -1, 64)
	}

	if o.Maxlength > 0 {
		t.attrs[kAttrMaxlength] = strconv.Itoa(o.Maxlength)
	}

	if o.Min != 0 {
		t.attrs[kAttrMin] = strconv.FormatFloat(o.Min, 'g', -1, 64)
	}

	if len(o.Name) > 0 {
		t.attrs[kAttrName] = o.Name
	}

	if len(o.Pattern) > 0 {
		t.attrs[kAttrPattern] = o.Pattern
	}

	if len(o.Placeholder) > 0 {
		t.attrs[kAttrPlaceholder] = o.Placeholder
	}

	if o.Readonly {
		t.attrs[kAttrReadonly] = ""
	}

	if o.Required {
		t.attrs[kAttrRequired] = ""
	}

	if o.Size > 0 {
		t.attrs[kAttrSize] = strconv.Itoa(o.Size)
	}

	if o.Step != 0 {
		t.attrs[kAttrStep] = strconv.FormatFloat(o.Step, 'g', -1, 64)
	}

	if len(o.Src) > 0 {
		t.attrs[kAttrSrc] = o.Src
	}

	// NOTE: we do not override t's assigned type if o.Type is not
	// specified, as the type must always exist.
	if len(o.Type) > 0 {
		t.attrs[kAttrType] = o.Type
	}

	if len(o.Value) > 0 {
		t.attrs[kAttrValue] = o.Value
	}

	if o.Width > 0 {
		t.attrs[kAttrWidth] = strconv.Itoa(o.Width)
	}

	t.isCacheClean = false
}

// The empty string will clear the attribute.
func (t *InputTag) SetValue(v string) {
	if len(v) == 0 {
		delete(t.attrs, kAttrValue)
	} else {
		t.attrs[kAttrValue] = v
	}
	t.isCacheClean = false
}

func (t *InputTag) Type() InputType {
	return InputType(t.attrs[kAttrType])
}

func (t *InputTag) Value() string {
	return t.attrs[kAttrValue]
}

type nullTag struct {
	baseTag
}

func newNullTag() *nullTag {
	return &nullTag{baseTag{
		tagType:  kTagTypeNull,
		children: make([]tagWriter, 0),
	}}
}

func (t *nullTag) Copy() Tag {
	return newNullTag()
}

// This just writes the children.
func (t *nullTag) write(writer io.Writer, env ...Environment) (n int, err error) {
	// Recursively write all children.
	for _, childTag := range t.children {
		if count, err := childTag.write(writer, env...); err != nil {
			return n, err
		} else {
			n += count
		}
	}
	return
}

func (t *nullTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	for ii, childTag := range t.children {
		if count, err := childTag.writePretty(writer, 0, env...); err != nil {
			return n, err
		} else {
			n += count
		}

		if ii == len(t.children)-1 {
			// Avoid the final newline.
			continue
		}

		// Pretty print with newline.
		if count, err := writeRune(writer, '\n'); err != nil {
			return n, err
		} else {
			n += count
		}
	}
	return
}

type OptionTag struct {
	baseTag
}

func (t *OptionTag) ResetOptions() {
	delete(t.attrs, kAttrDisabled)
	delete(t.attrs, kAttrLabel)
	delete(t.attrs, kAttrSelected)
	delete(t.attrs, kAttrValue)
}

func (t *OptionTag) Selected() bool {
	_, ok := t.attrs[kAttrSelected]
	return ok
}

func (t *OptionTag) SetOptions(o *OptionOptions) {
	if o.Disabled {
		t.attrs[kAttrDisabled] = ""
	}

	if len(o.Label) > 0 {
		t.attrs[kAttrLabel] = o.Label
	}

	if o.Selected {
		t.attrs[kAttrSelected] = ""
	}

	if len(o.Value) > 0 {
		t.attrs[kAttrValue] = o.Value
	}

	t.isCacheClean = false
}

func (t *OptionTag) SetSelected(selected bool) {
	if !selected {
		delete(t.attrs, kAttrSelected)
	} else {
		t.attrs[kAttrSelected] = ""
	}
	t.isCacheClean = false
}

func (t *OptionTag) Value() string {
	return t.attrs[kAttrValue]
}

type SelectTag struct {
	baseTag

	options []*OptionTag
}

func (t *SelectTag) Option(options ...*OptionOptions) *OptionTag {
	newTag := t.htmlGen.Option(options...)
	t.children = append(t.children, newTag)
	t.options = append(t.options, newTag)
	return newTag
}

func (t *SelectTag) OptionChildren() []*OptionTag {
	return t.options
}

func (t *SelectTag) RemoveChildren() Tag {
	t.children = t.children[:0]
	t.options = t.options[:0]
	return t
}

// Returns the selected OptionTag.
func (t *SelectTag) SelectedOption() (tag *OptionTag, err error) {
	for _, optionTag := range t.options {
		if optionTag.Selected() {
			tag = optionTag
			return
		}
	}
	err = ErrNotFound
	return
}

// Tags like <br />, etc.
type singleTag struct {
	baseTag
}

func newSingleTag(tagType int) *singleTag {
	return &singleTag{baseTag{
		tagType:     tagType,
		attrs:       make(map[int]string),
		customAttrs: make(map[string]string),
		children:    nil,
	}}
}

func (t *singleTag) Copy() Tag {
	if !t.isCacheClean {
		tagStr := tagTypeStringMap[t.tagType]
		var err error
		if t.cacheOpen, err = t.renderCacheOpen(tagStr); err == nil {
			t.isCacheClean = true
		}
		// Otherwise, leave the cache dirty.
	}

	return &singleTag{baseTag{
		tagType:      t.tagType,
		attrs:        copyAttrs(t.attrs),
		customAttrs:  copyCustomAttrs(t.customAttrs),
		children:     nil,
		isCacheClean: t.isCacheClean,
		cacheOpen:    t.cacheOpen,
	}}
}

func (t *singleTag) getChildren() []tagWriter {
	return []tagWriter{}
}

func (t *singleTag) RemoveChildren() Tag {
	// singleTags never have children.
	return t
}

// Renders the opening tag with attributes to a string.
func (t *singleTag) renderCacheOpen(tagStr string) (result string, err error) {
	tmp := new(bytes.Buffer)
	if _, err = t.writeOpenTagLead(tmp, tagStr); err != nil {
		return
	}
	// Finish the tag.
	if _, err = io.WriteString(tmp, " />"); err != nil {
		return
	}
	result = tmp.String()
	return
}

func (t *singleTag) write(writer io.Writer, env ...Environment) (n int, err error) {
	tagStr := tagTypeStringMap[t.tagType]

	// Write the tag.
	if !t.isCacheClean {
		if t.cacheOpen, err = t.renderCacheOpen(tagStr); err != nil {
			return
		}
		t.isCacheClean = true
	}
	return io.WriteString(writer, t.cacheOpen)
}

func (t *singleTag) writePretty(writer io.Writer, indent int, env ...Environment) (n int, err error) {
	if count, err := writeIndent(writer, indent); err != nil {
		return n, err
	} else {
		n += count
	}

	tagStr := tagTypeStringMap[t.tagType]

	// Write the tag.
	if count, err := t.writeOpenTagLeadSorted(writer, tagStr); err != nil {
		return n, err
	} else {
		n += count
	}
	// Finish the tag.
	if count, err := io.WriteString(writer, " />"); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

// Represents an HTML tag.
type Tag interface {
	// Note that all created tags will be children of the tag.
	TagFactory

	tagWriter

	// Adds the given tag as a child of the current tag and returns the
	// current object.
	AddChild(tag Tag) Tag
	AddChildText(tag *TextTag) Tag

	// Adds classes to tags set of classes.
	AddClass(classes ...string) Tag
	AddClasses(classes []string) Tag

	// Assigns the current tag to tag and returns it.
	Assign(tag *Tag) Tag

	// Returns a copy of the current tag.  Children will not be copied.
	// The Tag's cache will be updated at the time of the copy so that
	// the cache may be shared.
	Copy() Tag

	getChildren() []tagWriter

	// Hides the tag and its children during the rendering process.
	Hide(isHidden bool)

	// The empty string will be returned if no id exists.
	Id() string

	// Returns the parent tag or nil if it has no parents.
	Parent() Tag

	// Removes all children of the given Tag and returns the Tag.
	RemoveChildren() Tag

	// Removes the tag from its parent and returns itself.  This is
	// optimized for removing the most recently added child of a parent,
	// which will operate in O(1) time.
	RemoveParent() Tag

	// Attribute setters return the tag itself.  Pass an empty string
	// (or slice) to clear.

	RemoveAttribute(key string) Tag
	RemoveAttributes(keys ...string) Tag

	// Set a custom attribute.
	SetAttribute(key, value string) Tag
	SetAttributes(attrs map[string]string) Tag

	setChildren(children []tagWriter)

	SetClass(classes ...string) Tag
	SetClasses(classes []string) Tag

	SetId(id string) Tag

	// Form events
	SetOnblur(script string) Tag
	SetOnchange(script string) Tag
	SetOnfocus(script string) Tag
	SetOnselect(script string) Tag
	SetOnsubmit(script string) Tag

	// Keyboard events
	SetOnkeydown(script string) Tag
	SetOnkeypress(script string) Tag
	SetOnkeyup(script string) Tag

	// Mouse events
	SetOnclick(script string) Tag
	SetOndblclick(script string) Tag
	SetOnmousedown(script string) Tag
	SetOnmousemove(script string) Tag
	SetOnmouseout(script string) Tag
	SetOnmouseover(script string) Tag
	SetOnmouseup(script string) Tag

	// Sets parent of the tag to parent.  This does not manipulate the
	// children slices for the old parent or new parent.
	// As such, it must only be used to faciliate tag creation and is
	// not exposed publicly.
	setParent(parent Tag)

	SetTitle(title string) Tag

	// Move up count levels from tag.  Returns nil if no more ancestors
	// exist.  count always defaults to 1 if not specified or non-positive.
	Up(count ...int) Tag
}

type TagFactory interface {
	A(href string) Tag
	Abbr(title ...string) Tag
	Address() Tag

	B() Tag

	Blockquote() Tag

	Body() *BodyTag

	Br() Tag

	Button() Tag

	Canvas(options ...*CanvasOptions) Tag

	Caption() Tag

	// options is optional, and only the first one will be used.
	CheckedInput(inputType CheckedInputType, options ...*InputOptions) *CheckedInputTag
	CheckedInputTypeNameValue(inputType CheckedInputType, name, value string) *CheckedInputTag

	Cite() Tag

	Code() Tag

	// Comment node
	Comment() Tag

	Datalist() Tag

	Dfn() Tag

	// Creates a child Div element.
	Div() Tag
	// Creates a child Div element with specified classes.
	DivClasses(classes ...string) Tag
	// Creates a child Div element with specified id.
	DivId(id string) Tag
	// Creates a child Div element with specified id and classes.
	DivIdClasses(id string, classes ...string) Tag

	Dd() Tag
	Dl() Tag
	Dt() Tag

	Em() Tag

	Footer() Tag

	// options need not be specified.  If it is, only the first will be
	// used.
	Form(options ...*FormOptions) Tag

	H1() Tag
	H2() Tag
	H3() Tag
	H4() Tag
	H5() Tag
	H6() Tag

	Head() Tag

	Hr() Tag

	I() Tag

	Img(src, alt string, options ...*ImgOptions) Tag

	// options is optional, and only the first one will be used.
	Input(inputType InputType, options ...*InputOptions) *InputTag
	InputTypeNameValue(inputType InputType, name, value string) *InputTag

	Kbd() Tag

	Label(options ...*LabelOptions) Tag

	Li() Tag

	// options need not be specified; only the first will be used.
	Link(rel string, options ...*LinkOptions) Tag

	Meta(name, content string, options ...*MetaOptions) Tag

	NoScript() Tag

	Ol() Tag

	Option(options ...*OptionOptions) *OptionTag

	P() Tag
	Pre() Tag

	Samp() Tag

	Script(scriptType ...string) Tag
	ScriptSrc(scriptType, src string) Tag

	Small() Tag

	Span() Tag
	SpanClasses(classes ...string) Tag
	SpanId(id string) Tag
	SpanIdClasses(id string, classes ...string) Tag

	Select(options ...*SelectOptions) *SelectTag

	Strong() Tag

	Style(options ...*StyleOptions) Tag

	Table(options ...*TableOptions) Tag
	Tbody() Tag
	Td(options ...*TdOptions) Tag
	Tfoot() Tag
	Th(options ...*ThOptions) Tag
	Thead() Tag
	Title() Tag
	Tr() Tag

	Textarea(rows, cols int, options ...*TextareaOptions) Tag

	U() Tag

	Ul() Tag

	Var() Tag

	// Creates a new TextTag with contents text; only the first variadic
	// text is used.  If empty, an empty TextTag will be created.
	// text will be escaped automatically.  If escaping is not desired, use
	// TUnsafe().
	T(text ...string) *TextTag

	// This is a version of T() that does not escape text.
	TUnsafe(text ...string) *TextTag

	// Creates a new TextTagVar with contents text.  Variables within text will be interpreted using the
	// environment passed to one of the write methods.
	//
	// Escaping will be performed and applies to both variable values and text.
	TV(text ...string) *TextTagVar

	// This is a version of TV() that does not escape text.
	TVUnsafe(text ...string) *TextTagVar
}

type tagWriter interface {
	// Returns true if the tag is hidden.
	isHidden() bool

	// env is an optional environment for variable substitution in any supported tags.
	write(writer io.Writer, env ...Environment) (int, error)

	// Pretty printing.  Tag will be indented by 'indent' spaces.
	// This returns the number of bytes written to the writer.
	writePretty(writer io.Writer, indent int, env ...Environment) (int, error)
}

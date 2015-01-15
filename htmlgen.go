// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This provides a simple interface for generating HTML.
package htmlgen

import (
	"errors"
	"html"
	"io"
	"strconv"
	"strings"
)

const kIndentSpace = 2

// Patterns for input fields.
const (
	PatternFloat = `([0-9]+\.|[0-9]*\.?[0-9]+)`
	PatternInt   = `[0-9]*`
)

var (
	ErrNotFound = errors.New("not found")

	// Default tag factory for creating detached nodes.
	Factory = New()
)

// Factory for creating detached Tag nodes.
var H TagFactory

type CanvasOptions struct {
	Id     string
	Height int
	Width  int
}

// Optional form attributes.  We only specify those that are supported by
// the major browsers (IE, Firefox, Chrome, Safari).
type FormOptions struct {
	AcceptCharset string
	Action        string
	// Autocomplete is on by default.
	AutocompleteOff bool
	Enctype         string
	Method          string
	Name            string
	Target          string
}

type ImgOptions struct {
	// Height in pixels (non-zero to set)
	Height int
	Ismap  bool
	Usemap string
	// Width in pixels (non-zero to set)
	Width int
}

// Input tag parameters
type InputOptions struct {
	Action string
	Alt    string
	// Autocomplete is on by default.
	AutocompleteOff bool
	Checked         bool
	Disabled        bool
	// Only applied if > 0.
	Height int
	List   string
	// Only applied if != 0.
	Max float64
	// Only applied if > 0.
	Maxlength int
	// Only applied if != 0.
	Min         float64
	Name        string
	Pattern     string
	Placeholder string
	Readonly    bool
	Required    bool
	// Only applied if > 0.
	Size int
	// Only applied if != 0.
	Step  float64
	Src   string
	Type  string
	Value string
	// Only applied if > 0.
	Width int
}

type LabelOptions struct {
	For  string
	Form string
}

// Link tag attributes
type LinkOptions struct {
	Href     string
	Hreflang string
	Media    string
	Type     string
}

type MetaOptions struct {
	Charset   string
	Content   string
	HttpEquiv string
	Name      string
}

type OptionOptions struct {
	Disabled bool
	Label    string
	Selected bool
	Value    string
}

type SelectOptions struct {
	Autofocus bool
	Disabled  bool
	Form      string
	Multiple  bool
	Name      string
	Size      int
}

type StyleOptions struct {
	Media string
	Type  string
}

type TableOptions struct {
	Border bool
}

type TdOptions struct {
	// Only Opera supports Colspan=0, so we don't support it.
	Colspan int
	Headers string
	// Only Opera supports Rowspan=0, so we don't support it.
	Rowspan int
}

type ThOptions struct {
	// Colspan=0 is not supported (and will be omitted).
	Colspan int
	Headers string
	// Rowspan=0 is not supported (and will be omitted).
	Rowspan int
	Scope   string
}

type TextareaOptions struct {
	Autofocus   bool
	Disabled    bool
	Form        string
	Maxlength   int
	Name        string
	Placeholder string
	Readonly    bool
	Required    bool
	WrapHard    bool
}

type htmlGen struct {
}

func init() {
	H = New()
}

func New() TagFactory {
	return &htmlGen{}
}

func NewRoot() Tag {
	return &htmlTag{baseTag{
		tagType:     kTagTypeHtml,
		attrs:       make(map[int]string),
		customAttrs: make(map[string]string),
		children:    make([]tagWriter, 0),
	}}
}

// Creates a new null tag.
func NewNull() Tag {
	return newNullTag()
}

func (t *htmlGen) A(href string) Tag {
	newTag := newBaseTag(kTagTypeA)
	if len(href) > 0 {
		newTag.attrs[kAttrHref] = href
	}
	return newTag
}

func (t *htmlGen) Abbr(title ...string) Tag {
	tag := newBaseTag(kTagTypeAbbr)
	if len(title) > 0 {
		tag.setAttr(kAttrTitle, title[0])
	}
	return tag
}

func (t *htmlGen) Address() Tag {
	return newBaseTag(kTagTypeAddress)
}

func (t *htmlGen) B() Tag {
	return newBaseTag(kTagTypeB)
}

func (t *htmlGen) Blockquote() Tag {
	return newBaseTag(kTagTypeBlockquote)
}

func (t *htmlGen) Body() *BodyTag {
	return &BodyTag{baseTag{
		tagType:     kTagTypeBody,
		attrs:       make(map[int]string),
		customAttrs: make(map[string]string),
		children:    make([]tagWriter, 0),
	}}
}

func (t *htmlGen) Br() Tag {
	return newSingleTag(kTagTypeBr)
}

func (t *htmlGen) Button() Tag {
	return newBaseTag(kTagTypeButton)
}

func (t *htmlGen) Canvas(options ...*CanvasOptions) Tag {
	tag := newBaseTag(kTagTypeCanvas)

	if len(options) == 0 {
		return tag
	}

	o := options[0]
	if len(o.Id) > 0 {
		tag.SetId(o.Id)
	}

	if o.Height > 0 {
		tag.attrs[kAttrHeight] = strconv.Itoa(o.Height)
	}

	if o.Width > 0 {
		tag.attrs[kAttrWidth] = strconv.Itoa(o.Width)
	}

	return tag
}

func (t *htmlGen) Caption() Tag {
	return newBaseTag(kTagTypeCaption)
}

func (t *htmlGen) CheckedInput(inputType CheckedInputType, options ...*InputOptions) *CheckedInputTag {
	return &CheckedInputTag{InputTag: t.Input(InputType(inputType), options...)}
}

func (t *htmlGen) CheckedInputTypeNameValue(inputType CheckedInputType, name, value string) *CheckedInputTag {
	return &CheckedInputTag{InputTag: t.InputTypeNameValue(InputType(inputType), name, value)}
}

func (t *htmlGen) Cite() Tag {
	return newBaseTag(kTagTypeCite)
}

func (t *htmlGen) Code() Tag {
	return newBaseTag(kTagTypeCode)
}

func (t *htmlGen) Comment() Tag {
	return &commentTag{}
}

func (t *htmlGen) Datalist() Tag {
	return newBaseTag(kTagTypeDatalist)
}

func (t *htmlGen) Dfn() Tag {
	return newBaseTag(kTagTypeDfn)
}

func (t *htmlGen) Div() Tag {
	return newBaseTag(kTagTypeDiv)
}

func (t *htmlGen) DivClasses(classes ...string) Tag {
	newTag := newBaseTag(kTagTypeDiv)
	newTag.SetClasses(classes)
	return newTag
}

func (t *htmlGen) DivId(id string) Tag {
	newTag := newBaseTag(kTagTypeDiv)
	newTag.SetId(id)
	return newTag
}

func (t *htmlGen) DivIdClasses(id string, classes ...string) Tag {
	newTag := newBaseTag(kTagTypeDiv)
	newTag.SetClasses(classes)
	newTag.SetId(id)
	return newTag
}

func (t *htmlGen) Dd() Tag {
	return newBaseTag(kTagTypeDd)
}

func (t *htmlGen) Dl() Tag {
	return newBaseTag(kTagTypeDl)
}

func (t *htmlGen) Dt() Tag {
	return newBaseTag(kTagTypeDt)
}

func (t *htmlGen) Em() Tag {
	return newBaseTag(kTagTypeEm)
}

func (t *htmlGen) Footer() Tag {
	return newBaseTag(kTagTypeFooter)
}

func (t *htmlGen) Form(options ...*FormOptions) Tag {
	newTag := newBaseTag(kTagTypeForm)
	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	if len(o.AcceptCharset) > 0 {
		newTag.attrs[kAttrAcceptCharset] = o.AcceptCharset
	}

	if len(o.Action) > 0 {
		newTag.attrs[kAttrAction] = o.Action
	}

	if o.AutocompleteOff {
		newTag.attrs[kAttrAutocomplete] = kHtmlOff
	}

	if len(o.Enctype) > 0 {
		newTag.attrs[kAttrEnctype] = o.Enctype
	}

	if len(o.Method) > 0 {
		newTag.attrs[kAttrMethod] = o.Method
	}

	if len(o.Name) > 0 {
		newTag.attrs[kAttrName] = o.Name
	}

	if len(o.Target) > 0 {
		newTag.attrs[kAttrTarget] = o.Target
	}

	return newTag
}

func (t *htmlGen) H1() Tag {
	return newBaseTag(kTagTypeH1)
}

func (t *htmlGen) H2() Tag {
	return newBaseTag(kTagTypeH2)
}

func (t *htmlGen) H3() Tag {
	return newBaseTag(kTagTypeH3)
}

func (t *htmlGen) H4() Tag {
	return newBaseTag(kTagTypeH4)
}

func (t *htmlGen) H5() Tag {
	return newBaseTag(kTagTypeH5)
}

func (t *htmlGen) H6() Tag {
	return newBaseTag(kTagTypeH6)
}

func (t *htmlGen) Head() Tag {
	return newBaseTag(kTagTypeHead)
}

func (t *htmlGen) Hr() Tag {
	return newSingleTag(kTagTypeHr)
}

func (t *htmlGen) I() Tag {
	return newBaseTag(kTagTypeI)
}

func (t *htmlGen) Img(src, alt string, options ...*ImgOptions) Tag {
	newTag := newSingleTag(kTagTypeImg)

	newTag.attrs[kAttrSrc] = src
	newTag.attrs[kAttrAlt] = alt

	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	if o.Height > 0 {
		newTag.attrs[kAttrHeight] = strconv.Itoa(o.Height)
	}

	if o.Ismap {
		newTag.attrs[kAttrIsmap] = ""
	}

	if len(o.Usemap) > 0 {
		newTag.attrs[kAttrUsemap] = o.Usemap
	}

	if o.Width > 0 {
		newTag.attrs[kAttrWidth] = strconv.Itoa(o.Width)
	}

	return newTag
}

func (t *htmlGen) Input(inputType InputType, options ...*InputOptions) *InputTag {
	newTag := &InputTag{*newSingleTag(kTagTypeInput)}

	if len(inputType) > 0 {
		newTag.attrs[kAttrType] = string(inputType)
	}

	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	newTag.SetOptions(o)

	return newTag
}

func (t *htmlGen) InputTypeNameValue(inputType InputType, name, value string) *InputTag {
	newTag := &InputTag{*newSingleTag(kTagTypeInput)}

	if len(inputType) > 0 {
		newTag.attrs[kAttrType] = string(inputType)
	}

	if len(name) > 0 {
		newTag.attrs[kAttrName] = name
	}

	if len(value) > 0 {
		newTag.attrs[kAttrValue] = value
	}

	return newTag
}

func (t *htmlGen) Kbd() Tag {
	return newBaseTag(kTagTypeKbd)
}

func (t *htmlGen) Label(options ...*LabelOptions) Tag {
	newTag := newBaseTag(kTagTypeLabel)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if len(o.For) > 0 {
		newTag.attrs[kAttrFor] = o.For
	}

	if len(o.Form) > 0 {
		newTag.attrs[kAttrForm] = o.Form
	}

	return newTag
}

func (t *htmlGen) Li() Tag {
	return newBaseTag(kTagTypeLi)
}

func (t *htmlGen) Link(rel string, options ...*LinkOptions) Tag {
	newTag := newSingleTag(kTagTypeLink)

	// rel is required.
	newTag.attrs[kAttrRel] = rel

	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	if len(o.Href) > 0 {
		newTag.attrs[kAttrHref] = o.Href
	}
	if len(o.Hreflang) > 0 {
		newTag.attrs[kAttrHreflang] = o.Hreflang
	}
	if len(o.Media) > 0 {
		newTag.attrs[kAttrMedia] = o.Media
	}
	if len(o.Type) > 0 {
		newTag.attrs[kAttrType] = o.Type
	}
	return newTag
}

func (t *htmlGen) Meta(name, content string, options ...*MetaOptions) Tag {
	newTag := newSingleTag(kTagTypeMeta)

	newTag.attrs[kAttrName] = name
	newTag.attrs[kAttrContent] = content

	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	if len(o.Charset) > 0 {
		newTag.attrs[kAttrCharset] = o.Charset
	}
	if len(o.Content) > 0 {
		newTag.attrs[kAttrContent] = o.Content
	}
	if len(o.HttpEquiv) > 0 {
		newTag.attrs[kAttrHttpEquiv] = o.HttpEquiv
	}
	if len(o.Name) > 0 {
		newTag.attrs[kAttrName] = o.Name
	}
	return newTag
}

func (t *htmlGen) NoScript() Tag {
	return newBaseTag(kTagTypeNoScript)
}

func (t *htmlGen) Ol() Tag {
	return newBaseTag(kTagTypeOl)
}

func (t *htmlGen) Option(options ...*OptionOptions) *OptionTag {
	newTag := &OptionTag{*newBaseTag(kTagTypeOption)}

	if len(options) == 0 {
		return newTag
	}

	o := options[0]
	newTag.SetOptions(o)

	return newTag
}

func (t *htmlGen) P() Tag {
	return newBaseTag(kTagTypeP)
}

func (t *htmlGen) Pre() Tag {
	return newBaseTag(kTagTypePre)
}

func (t *htmlGen) Samp() Tag {
	return newBaseTag(kTagTypeSamp)
}

func (t *htmlGen) Script(optScriptType ...string) Tag {
	scriptType := ""
	if len(optScriptType) > 0 {
		scriptType = optScriptType[0]
	}
	return t.ScriptSrc(scriptType, "")
}

func (t *htmlGen) ScriptSrc(scriptType, src string) Tag {
	newTag := newBaseTag(kTagTypeScript)
	if len(scriptType) > 0 {
		newTag.attrs[kAttrType] = scriptType
	}
	if len(src) > 0 {
		newTag.attrs[kAttrSrc] = src
	}
	return newTag
}

func (t *htmlGen) Select(options ...*SelectOptions) *SelectTag {
	newTag := &SelectTag{
		baseTag: *newBaseTag(kTagTypeSelect),
		options: make([]*OptionTag, 0),
	}

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if o.Autofocus {
		newTag.attrs[kAttrAutofocus] = ""
	}

	if o.Disabled {
		newTag.attrs[kAttrDisabled] = ""
	}

	if len(o.Form) > 0 {
		newTag.attrs[kAttrForm] = o.Form
	}

	if o.Multiple {
		newTag.attrs[kAttrMultiple] = ""
	}

	if len(o.Name) > 0 {
		newTag.attrs[kAttrName] = o.Name
	}

	if o.Size > 0 {
		newTag.attrs[kAttrSize] = strconv.Itoa(o.Size)
	}

	return newTag
}

func (t *htmlGen) Small() Tag {
	return newBaseTag(kTagTypeSmall)
}

func (t *htmlGen) Span() Tag {
	return newBaseTag(kTagTypeSpan)
}

func (t *htmlGen) SpanClasses(classes ...string) Tag {
	newTag := newBaseTag(kTagTypeSpan)
	newTag.SetClasses(classes)
	return newTag
}

func (t *htmlGen) SpanId(id string) Tag {
	newTag := newBaseTag(kTagTypeSpan)
	newTag.SetId(id)
	return newTag
}

func (t *htmlGen) SpanIdClasses(id string, classes ...string) Tag {
	newTag := newBaseTag(kTagTypeSpan)
	newTag.SetId(id)
	newTag.SetClasses(classes)
	return newTag
}

func (t *htmlGen) Strong() Tag {
	return newBaseTag(kTagTypeStrong)
}

func (t *htmlGen) Style(options ...*StyleOptions) Tag {
	newTag := newBaseTag(kTagTypeStyle)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if len(o.Media) > 0 {
		newTag.attrs[kAttrMedia] = o.Media
	}
	if len(o.Type) > 0 {
		newTag.attrs[kAttrType] = o.Type
	}

	return newTag
}

func (t *htmlGen) T(text ...string) *TextTag {
	textStr := ""
	if len(text) > 0 {
		textStr = html.EscapeString(text[0])
	}
	return &TextTag{
		text: textStr,
	}
}

func (t *htmlGen) Table(options ...*TableOptions) Tag {
	newTag := newBaseTag(kTagTypeTable)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if o.Border {
		newTag.attrs[kAttrBorder] = kHtmlBorderOn
	} else {
		newTag.attrs[kAttrBorder] = kHtmlBorderOff
	}

	return newTag
}

func (t *htmlGen) Tbody() Tag {
	return newBaseTag(kTagTypeTbody)
}

func (t *htmlGen) Td(options ...*TdOptions) Tag {
	newTag := newBaseTag(kTagTypeTd)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if o.Colspan > 0 {
		newTag.attrs[kAttrColspan] = strconv.Itoa(o.Colspan)
	}

	if len(o.Headers) > 0 {
		newTag.attrs[kAttrHeaders] = o.Headers
	}

	if o.Rowspan > 0 {
		newTag.attrs[kAttrRowspan] = strconv.Itoa(o.Rowspan)
	}

	return newTag
}

func (t *htmlGen) Textarea(rows, cols int, options ...*TextareaOptions) Tag {
	newTag := newBaseTag(kTagTypeTextarea)
	newTag.attrs[kAttrRows] = strconv.Itoa(rows)
	newTag.attrs[kAttrCols] = strconv.Itoa(cols)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if o.Autofocus {
		newTag.attrs[kAttrAutofocus] = ""
	}

	if o.Disabled {
		newTag.attrs[kAttrDisabled] = ""
	}

	if len(o.Form) > 0 {
		newTag.attrs[kAttrForm] = o.Form
	}

	if o.Maxlength > 0 {
		newTag.attrs[kAttrMaxlength] = strconv.Itoa(o.Maxlength)
	}

	if len(o.Name) > 0 {
		newTag.attrs[kAttrName] = o.Name
	}

	if len(o.Placeholder) > 0 {
		newTag.attrs[kAttrPlaceholder] = o.Placeholder
	}

	if o.Readonly {
		newTag.attrs[kAttrReadonly] = ""
	}

	if o.Required {
		newTag.attrs[kAttrRequired] = ""
	}

	if o.WrapHard {
		newTag.attrs[kAttrWrap] = kHtmlWrapHard
	}

	return newTag
}

func (t *htmlGen) Tfoot() Tag {
	return newBaseTag(kTagTypeTfoot)
}

func (t *htmlGen) Th(options ...*ThOptions) Tag {
	newTag := newBaseTag(kTagTypeTh)

	if len(options) == 0 {
		return newTag
	}

	o := options[0]

	if o.Colspan > 0 {
		newTag.attrs[kAttrColspan] = strconv.Itoa(o.Colspan)
	}

	if len(o.Headers) > 0 {
		newTag.attrs[kAttrHeaders] = o.Headers
	}

	if o.Rowspan > 0 {
		newTag.attrs[kAttrRowspan] = strconv.Itoa(o.Rowspan)
	}

	if len(o.Scope) > 0 {
		newTag.attrs[kAttrScope] = o.Scope
	}

	return newTag
}

func (t *htmlGen) Thead() Tag {
	return newBaseTag(kTagTypeThead)
}

func (t *htmlGen) Title() Tag {
	return newBaseTag(kTagTypeTitle)
}

func (t *htmlGen) Tr() Tag {
	return newBaseTag(kTagTypeTr)
}

func (t *htmlGen) TUnsafe(text ...string) *TextTag {
	textStr := ""
	if len(text) > 0 {
		textStr = text[0]
	}
	return &TextTag{
		text: textStr,
	}
}

func (t *htmlGen) TV(text ...string) *TextTagVar {
	textStr := ""
	if len(text) > 0 {
		textStr = text[0]
	}
	return NewTextTagVar(textStr)
}

func (t *htmlGen) TVUnsafe(text ...string) *TextTagVar {
	textStr := ""
	if len(text) > 0 {
		textStr = text[0]
	}
	tag := NewTextTagVar(textStr)
	tag.SetUnsafe(true)
	return tag
}

func (t *htmlGen) U() Tag {
	return newBaseTag(kTagTypeU)
}

func (t *htmlGen) Ul() Tag {
	return newBaseTag(kTagTypeUl)
}

func (t *htmlGen) Var() Tag {
	return newBaseTag(kTagTypeVar)
}

// Returns a copy of attrs.
func copyAttrs(attrs map[int]string) map[int]string {
	dest := make(map[int]string, len(attrs))
	for k, v := range attrs {
		dest[k] = v
	}
	return dest
}

func copyCustomAttrs(attrs map[string]string) map[string]string {
	dest := make(map[string]string, len(attrs))
	for k, v := range attrs {
		dest[k] = v
	}
	return dest
}

// Writes the tag tree from the given root tag.
// env is an optional environment
func Write(writer io.Writer, root Tag, env ...Environment) (int, error) {
	return root.write(writer, env...)
}

func writeIndent(writer io.Writer, indent int) (int, error) {
	indentStr := strings.Repeat(" ", indent)
	return io.WriteString(writer, indentStr)
}

func writeKeyValue(writer io.Writer, key, value string) (n int, err error) {
	if count, err := io.WriteString(writer, key); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := io.WriteString(writer, "=\""); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := io.WriteString(writer, value); err != nil {
		return n, err
	} else {
		n += count
	}
	if count, err := writeRune(writer, '"'); err != nil {
		return n, err
	} else {
		n += count
	}
	return
}

// Writes HTML 'key="value"' tag attributes.  keyId is a kAttr constant.
func writeKeyIdValue(writer io.Writer, keyId int, value string) (int, error) {
	keyStr := attrStringMap[keyId]
	return writeKeyValue(writer, keyStr, value)
}

func WritePretty(writer io.Writer, root Tag, env ...Environment) (int, error) {
	return root.writePretty(writer, 0, env...)
}

func writeRune(writer io.Writer, ch rune) (int, error) {
	return io.WriteString(writer, string(ch))
}

// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

const (
	// NOTE: SetAttribute() assigns a custom attribute, so this will not
	// overwrite a pre-existing href attribute in a A(href) definition.
	//
	// For that use case, pass an empty string to A, and then use
	// SetAttribute:
	//
	//   a := H.A("")
	//   ...
	//   a.SetAttribute(AttrHref, "foo")
	AttrHref = "href"
)

// Class attributes
// This must match the order of attrStringMap.
const (
	kAttrAcceptCharset = iota
	kAttrAction        = iota
	kAttrAlt           = iota
	kAttrAutocomplete  = iota
	kAttrAutofocus     = iota
	kAttrBorder        = iota
	kAttrCharset       = iota
	kAttrChecked       = iota
	kAttrClass         = iota
	kAttrCols          = iota
	kAttrColspan       = iota
	kAttrContent       = iota
	kAttrDisabled      = iota
	kAttrEnctype       = iota
	kAttrFor           = iota
	kAttrForm          = iota
	kAttrHeaders       = iota
	kAttrHeight        = iota
	kAttrHref          = iota
	kAttrHreflang      = iota
	kAttrHttpEquiv     = iota
	kAttrId            = iota
	kAttrIsmap         = iota
	kAttrLabel         = iota
	kAttrList          = iota
	kAttrMax           = iota
	kAttrMaxlength     = iota
	kAttrMedia         = iota
	kAttrMethod        = iota
	kAttrMin           = iota
	kAttrMultiple      = iota
	kAttrName          = iota
	kAttrOnblur        = iota
	kAttrOnchange      = iota
	kAttrOnclick       = iota
	kAttrOndblclick    = iota
	kAttrOnfocus       = iota
	kAttrOnload        = iota
	kAttrOnkeydown     = iota
	kAttrOnkeypress    = iota
	kAttrOnkeyup       = iota
	kAttrOnmousedown   = iota
	kAttrOnmousemove   = iota
	kAttrOnmouseout    = iota
	kAttrOnmouseover   = iota
	kAttrOnmouseup     = iota
	kAttrOnselect      = iota
	kAttrOnsubmit      = iota
	kAttrPattern       = iota
	kAttrPlaceholder   = iota
	kAttrRowspan       = iota
	kAttrReadonly      = iota
	kAttrRel           = iota
	kAttrRequired      = iota
	kAttrRows          = iota
	kAttrScope         = iota
	kAttrSelected      = iota
	kAttrSize          = iota
	kAttrSrc           = iota
	kAttrStep          = iota
	kAttrTarget        = iota
	kAttrTitle         = iota
	kAttrType          = iota
	kAttrUsemap        = iota
	kAttrValue         = iota
	kAttrWidth         = iota
	kAttrWrap          = iota

	// The number of attributes.
	kAttrMAXCOUNT = iota
)

// HTML constants
const (
	kHtmlOff = "off"
	kHtmlOn  = "on"

	kHtmlBorderOn  = "1"
	kHtmlBorderOff = ""

	kHtmlWrapHard = "hard"
)

// This must match the order of tagTypeStringMap.
const (
	kTagTypeA          = iota
	kTagTypeAbbr       = iota
	kTagTypeAddress    = iota
	kTagTypeB          = iota
	kTagTypeBlockquote = iota
	kTagTypeBody       = iota
	kTagTypeBr         = iota
	kTagTypeButton     = iota
	kTagTypeCanvas     = iota
	kTagTypeCaption    = iota
	kTagTypeCite       = iota
	kTagTypeCode       = iota
	kTagTypeDatalist   = iota
	kTagTypeDfn        = iota
	kTagTypeDiv        = iota
	kTagTypeDd         = iota
	kTagTypeDl         = iota
	kTagTypeDt         = iota
	kTagTypeEm         = iota
	kTagTypeFooter     = iota
	kTagTypeForm       = iota
	kTagTypeH1         = iota
	kTagTypeH2         = iota
	kTagTypeH3         = iota
	kTagTypeH4         = iota
	kTagTypeH5         = iota
	kTagTypeH6         = iota
	kTagTypeHead       = iota
	kTagTypeHr         = iota
	kTagTypeHtml       = iota
	kTagTypeI          = iota
	kTagTypeImg        = iota
	kTagTypeInput      = iota
	kTagTypeKbd        = iota
	kTagTypeLabel      = iota
	kTagTypeLi         = iota
	kTagTypeLink       = iota
	kTagTypeMeta       = iota
	kTagTypeNoScript   = iota
	kTagTypeOl         = iota
	kTagTypeOption     = iota
	kTagTypeP          = iota
	kTagTypePre        = iota
	kTagTypeSamp       = iota
	kTagTypeScript     = iota
	kTagTypeSelect     = iota
	kTagTypeSmall      = iota
	kTagTypeSpan       = iota
	kTagTypeStrong     = iota
	kTagTypeStyle      = iota
	kTagTypeTable      = iota
	kTagTypeTbody      = iota
	kTagTypeTd         = iota
	kTagTypeTextarea   = iota
	kTagTypeTfoot      = iota
	kTagTypeTh         = iota
	kTagTypeThead      = iota
	kTagTypeTitle      = iota
	kTagTypeTr         = iota
	kTagTypeU          = iota
	kTagTypeUl         = iota
	kTagTypeVar        = iota

	// The number of tag types.
	kTagTypeMAXCOUNT = iota

	// The Null tag prints nothing.
	kTagTypeNull = iota
)

type CheckedInputType string

const (
	CheckedInputTypeCheckbox = CheckedInputType(InputTypeCheckbox)
	CheckedInputTypeRadio    = CheckedInputType(InputTypeRadio)
)

type InputType string

const (
	InputTypeButton   = InputType("button")
	InputTypeCheckbox = InputType("checkbox")
	InputTypeFile     = InputType("file")
	InputTypeHidden   = InputType("hidden")
	InputTypeImage    = InputType("image")
	InputTypeNumber   = InputType("number")
	InputTypePassword = InputType("password")
	InputTypeRadio    = InputType("radio")
	InputTypeRange    = InputType("range")
	InputTypeReset    = InputType("reset")
	InputTypeSubmit   = InputType("submit")
	InputTypeText     = InputType("text")
)

// Common HTML constants.
const (
	CharsetUnicode = "UTF-8"
	CharsetLatin   = "ISO-8859-1"

	FormEnctypeData       = "multipart/form-data"
	FormEnctypePlain      = "text/plain"
	FormEnctypeUrlencoded = "application/x-www-form-urlencoded"

	FormMethodGet  = "get"
	FormMethodPost = "post"

	FormTargetBlank = "_blank"
	// The default target
	FormTargetSelf   = "_self"
	FormTargetParent = "_parent"
	FormTargetTop    = "_top"

	HeaderContentType = "Content-Type"

	LinkRelAlternate    = "alternate"
	LinkRelAuthor       = "author"
	LinkRelHelp         = "help"
	LinkRelIcon         = "icon"
	LinkRelLicense      = "license"
	LinkRelNext         = "next"
	LinkRelNofollow     = "nofollow"
	LinkRelNoreferrer   = "noreferrer"
	LinkRelPrefetch     = "prefetch"
	LinkRelPrev         = "prev"
	LinkRelSearch       = "search"
	LinkRelShortcutIcon = "shortcut icon"
	LinkRelStylesheet   = "stylesheet"

	MetaNameViewport = "viewport"

	// The default media for <style>
	StyleMediaAll        = "all"
	StyleMediaHandheld   = "handheld"
	StyleMediaPrint      = "print"
	StyleMediaProjection = "projection"
	StyleMediaScreen     = "screen"
	StyleMediaTv         = "tv"

	MimeTypeCss        = "text/css"
	MimeTypeHtml       = "text/html"
	MimeTypeJavascript = "text/javascript"

	MimeTypeGif = "image/gif"
	MimeTypePng = "image/png"
)

// Indexed by kAttr
var attrStringMap [kAttrMAXCOUNT]string

// Indexed by tagType
var tagTypeStringMap [kTagTypeMAXCOUNT]string

func init() {
	// This must match the order of the kAttrs.
	attrStringMap = [kAttrMAXCOUNT]string{
		"accept-charset",
		"action",
		"alt",
		"autocomplete",
		"autofocus",
		"border",
		"charset",
		"checked",
		"class",
		"cols",
		"colspan",
		"content",
		"disabled",
		"enctype",
		"for",
		"form",
		"headers",
		"height",
		"href",
		"hreflang",
		"http-equiv",
		"id",
		"ismap",
		"label",
		"list",
		"max",
		"maxlength",
		"media",
		"method",
		"min",
		"multiple",
		"name",
		"onblur",
		"onchange",
		"onclick",
		"ondblclick",
		"onfocus",
		"onload",
		"onkeydown",
		"onkeypress",
		"onkeyup",
		"onmousedown",
		"onmousemove",
		"onmouseout",
		"onmouseover",
		"onmouseup",
		"onselect",
		"onsubmit",
		"pattern",
		"placeholder",
		"rowspan",
		"readonly",
		"rel",
		"required",
		"rows",
		"scope",
		"selected",
		"size",
		"src",
		"step",
		"target",
		"title",
		"type",
		"usemap",
		"value",
		"width",
		"wrap",
	}

	// This must match the order of the kTagTypes.
	tagTypeStringMap = [kTagTypeMAXCOUNT]string{
		"a",
		"abbr",
		"address",
		"b",
		"blockquote",
		"body",
		"br",
		"button",
		"canvas",
		"caption",
		"cite",
		"code",
		"datalist",
		"dfn",
		"div",
		"dd",
		"dl",
		"dt",
		"em",
		"footer",
		"form",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"head",
		"hr",
		"html",
		"i",
		"img",
		"input",
		"kbd",
		"label",
		"li",
		"link",
		"meta",
		"noscript",
		"ol",
		"option",
		"p",
		"pre",
		"samp",
		"script",
		"select",
		"small",
		"span",
		"strong",
		"style",
		"table",
		"tbody",
		"td",
		"textarea",
		"tfoot",
		"th",
		"thead",
		"title",
		"tr",
		"u",
		"ul",
		"var",
	}
}

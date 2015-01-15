// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func compareHtml(root Tag, cmp string, usePretty bool, env ...Environment) error {
	buf := new(bytes.Buffer)
	var err error
	if usePretty {
		_, err = WritePretty(buf, root, env...)
	} else {
		_, err = Write(buf, root, env...)
	}
	if err != nil {
		return err
	}

	if pos, err := stringCmp(buf.String(), cmp); err != nil {
		logBuff := &bytes.Buffer{}
		logBuff.WriteString(err.Error())
		logBuff.WriteString(fmt.Sprintf("\nOutput: %s\n", buf.String()))
		logBuff.WriteString(fmt.Sprintf("\nComparison: %s\n", cmp))
		logBuff.WriteString(fmt.Sprintf("\nError: %s\n", err))
		logBuff.WriteString(buf.String()[pos:])
		return errors.New(logBuff.String())
	}

	return nil
}

func Test_Copy(t *testing.T) {
	const kCompare = `<div class="foo"></div>`
	const kCopyCompare = `<div class="foo" id="copy"></div>`
	const kCommentCompare = `<!--
  commented text
-->`
	const kCommentCopyCompare = `<!-- -->`
	const kSingleCompare = `<img alt="alt" src="src" />`
	const kSingleCopyCompare = `<img alt="alt" id="img_copy" src="src" />`
	const kTextCompare = "hello"
	const kTextCopyCompare = "world"

	r := NewNull()
	tag := r.Comment()
	// Add a child
	tag.T("commented text")

	tagCopy := tag.Copy()
	if err := compare(t, tag, kCommentCompare); err != nil {
		t.Error(err)
	}
	if err := compare(t, tagCopy, kCommentCopyCompare); err != nil {
		t.Error(err)
	}

	// Test baseTag
	tag = r.DivClasses("foo")
	tagCopy = tag.Copy()
	if err := compare(t, tag, kCompare); err != nil {
		t.Error(err)
	}
	if err := compare(t, tagCopy, kCompare); err != nil {
		t.Error(err)
	}

	// Now, dirty the copy.
	tagCopy.SetId("copy")
	if err := compare(t, tag, kCompare); err != nil {
		t.Error(err)
	}
	if err := compare(t, tagCopy, kCopyCompare); err != nil {
		t.Error(err)
	}

	// Test singleTag
	tag = r.Img("src", "alt")
	tagCopy = tag.Copy()
	if err := compare(t, tag, kSingleCompare); err != nil {
		t.Error(err)
	}
	if err := compare(t, tagCopy, kSingleCompare); err != nil {
		t.Error(err)
	}

	// Now, dirty the copy.
	tagCopy.SetId("img_copy")
	if err := compare(t, tag, kSingleCompare); err != nil {
		t.Error(err)
	}
	if err := compare(t, tagCopy, kSingleCopyCompare); err != nil {
		t.Error(err)
	}

	textRoot := NewNull()
	textTag := textRoot.T("hello")
	textTagCopy := textTag.Copy()
	// The copy should not modify the original.
	textTagCopy.SetText("world")
	if err := compare(t, textRoot, kTextCompare); err != nil {
		t.Error(err)
	}
}

func Test_Htmlgen(t *testing.T) {
	const kCompare = `<!DOCTYPE html>
<html>
  <head>
    <link media="screen" rel="search" />
    <link rel="stylesheet" title="title" />
    <style media="print" type="text/css">
      <!--
        @import url(foo)
      -->
    </style>
    <title>
      title
    </title>
  </head>
  <body onblur="onblur" onfocus="onfocus" onload="onload">
    <div></div>
    <br />
    <button></button>
    <div id="fooid"></div>
    <div class="class1 class2 class3"></div>
    <div class="fooclass1 fooclass2"></div>
    <div class="fooclass1 fooclass2" id="fooidall"></div>
    <span>
      hello,
      world
    </span>
    <!--
      hello, comment!
    -->
    <script type="text/javascript">
      var hello;
    </script>
    <script src="foo.js" type="text/javascript"></script>
    <noscript>
      hi
    </noscript>
    <h1>
      h1 bold
    </h1>
    <a href="#">
      link
    </a>
    <footer>
      foot
    </footer>
    <form>
      foo
    </form>
    <form action="action" method="method" name="name">
      foo2
    </form>
    <form>
      <input name="name" type="hidden" value="value" />
    </form>
    <form>
      <input name="name" src="src" type="image" />
    </form>
    <form>
      <input type="text" value="image" />
    </form>
    <p>
      paragraph
    </p>
    <pre></pre>
    <ol>
      <li>
        foo
      </li>
      <li>
        foo bar
      </li>
    </ol>
    <ul>
      <li>
        foo
      </li>
    </ul>
    <table>
      <thead></thead>
      <tbody></tbody>
      <tfoot></tfoot>
    </table>
    <table border="1">
      <tr>
        <td colspan="1" headers="headers" rowspan="1"></td>
      </tr>
    </table>
    <abbr></abbr>
    <abbr title="hello"></abbr>
    <address></address>
    <b></b>
    <blockquote></blockquote>
    <dd></dd>
    <dl></dl>
    <dt></dt>
    <em></em>
    <hr />
    <i></i>
    <label for="for" form="form"></label>
    <small></small>
    <strong></strong>
    <canvas></canvas>
    <canvas height="1" id="foo" width="2"></canvas>
    <caption></caption>
    <cite></cite>
    <code></code>
    <samp></samp>
    <dfn></dfn>
    <datalist></datalist>
    <var></var>
    <kbd></kbd>
    <u></u>
    <th></th>
    <textarea cols="4" disabled="" name="foo" readonly="" rows="3"></textarea>
    <select form="1,2,3" multiple="" name="foo"></select>
    <option disabled="" label="label" selected="" value="value">
      Option
    </option>
    <div>
      hello
    </div>
    <div onchange="onchange" onkeydown="onkeydown" onkeypress="onkeypress" onkeyup="onkeyup" onselect="onselect" onsubmit="onsubmit"></div>
    <div onclick="onclick"></div>
    <div ondblclick="ondblclick"></div>
    <div onmousedown="onmousedown"></div>
    <div onmousemove="onmousemove"></div>
    <div onmouseout="onmouseout"></div>
    <div onmouseover="onmouseover"></div>
    <div onmouseup="onmouseup"></div>
    <div>
      
    </div>
    <p>
      Hello, foo.
    </p>
    <p>
      Hello, a&gt;b &amp; foo.
    </p>
    <p>
      Hello, a>b & foo.
    </p>
    <div attribute="yes" custom="custom"></div>
    <div></div>
    <div attribute="yes" custom="custom" multiple="1"></div>
    <meta charset="charset" content="content" http-equiv="equiv" name="name" />
    &gt;&lt;&amp;&#34;>
    ><&"
  </body>
</html>`

	root := NewRoot()

	// Head
	head := root.Head()
	head.Link(LinkRelSearch, &LinkOptions{Media: "screen"})
	head.Link(LinkRelStylesheet).SetTitle("title")
	head.Style(&StyleOptions{Type: MimeTypeCss, Media: StyleMediaPrint}).Comment().T("@import url(foo)")
	head.Title().T("title")

	// Body
	b := root.Body()
	b.Div()
	b.Br()
	b.Button()
	b.DivId("fooid")
	b.Div().AddClass("class1").AddClass().AddClass("class2", "class3")
	b.DivClasses("fooclass1", "fooclass2")
	b.DivIdClasses("fooidall", "fooclass1", "fooclass2")
	s := b.Span()
	s.T("hello,")
	s.T("world")
	b.Comment().T("hello, comment!")
	b.Script(MimeTypeJavascript).T("var hello;")
	b.ScriptSrc(MimeTypeJavascript, "foo.js")
	b.NoScript().T("hi")
	b.H1().T("h1 bold")
	b.A("#").T("link")
	b.Footer().T("foot")
	b.Form().T("foo")
	b.Form(&FormOptions{Action: "action", Method: "method", Name: "name"}).T("foo2")
	b.Form().InputTypeNameValue(InputTypeHidden, "name", "value")
	b.Form().Input(InputTypeImage, &InputOptions{Name: "name", Src: "src"})
	b.Form().Input(InputTypeText).SetValue("image")
	b.P().T("paragraph")
	b.Pre()
	ol := b.Ol()
	ol.Li().T("foo")
	ol.Li().T("foo bar")
	b.Ul().Li().T("foo")
	tbl := b.Table()
	tbl.Thead()
	tbl.Tbody()
	tbl.Tfoot()

	tbl = b.Table(&TableOptions{Border: true})
	tbl.Tr().Td(&TdOptions{Colspan: 1, Headers: "headers", Rowspan: 1})

	b.Abbr()
	b.Abbr("hello")
	b.Address()
	b.B()
	b.Blockquote()
	b.Dd()
	b.Dl()
	b.Dt()
	b.Em()
	b.Hr()
	b.I()
	b.Label(&LabelOptions{For: "for", Form: "form"})
	b.Small()
	b.Strong()
	b.Canvas()
	b.Canvas(&CanvasOptions{Id: "foo", Height: 1, Width: 2})
	b.Caption()
	b.Cite()
	b.Code()
	b.Samp()
	b.Dfn()
	b.Datalist()
	b.Var()
	b.Kbd()
	b.U()
	b.Th()
	b.Textarea(3, 4, &TextareaOptions{
		Disabled: true,
		Name:     "foo",
		Readonly: true,
	})
	b.Select(&SelectOptions{
		Form:     "1,2,3",
		Multiple: true,
		Name:     "foo",
	})
	b.Option(&OptionOptions{
		Disabled: true,
		Label:    "label",
		Selected: true,
		Value:    "value",
	}).T("Option")

	// Test adding a child.  We'll also test the Assign() operation for
	// binding a tag variable to the current tag.
	var div Tag
	b.AddChild(H.Div().Assign(&div))
	div.AddChildText(H.T("hello"))

	b.SetOnblur("onblur")
	b.SetOnfocus("onfocus")
	b.SetOnload("onload")
	b.Div().SetOnchange("onchange").SetOnselect("onselect").SetOnsubmit("onsubmit").SetOnkeydown("onkeydown").SetOnkeypress("onkeypress").SetOnkeyup("onkeyup")
	b.Div().SetOnclick("onclick")
	b.Div().SetOndblclick("ondblclick")
	b.Div().SetOnmousedown("onmousedown")
	b.Div().SetOnmousemove("onmousemove")
	b.Div().SetOnmouseout("onmouseout")
	b.Div().SetOnmouseover("onmouseover")
	div = b.Div().SetOnmouseup("onmouseup")

	// Test removal.
	div.T("removed 1")
	div.T("removed 2")
	div.RemoveChildren()

	b.Div().T()

	b.P().TV("Hello, $name.")
	b.P().TV("Hello, $escape & $name.")
	b.P().TVUnsafe("Hello, $escape & $name.")

	b.Div().SetAttribute("custom", "custom").SetAttribute("attribute", "yes")
	b.Div().SetAttribute("custom", "custom").RemoveAttribute("custom")
	b.Div().SetAttributes(map[string]string{"custom": "custom", "attribute": "yes", "multiple": "1", "extra": "2", "foo": "bar"}).RemoveAttributes("extra", "foo")
	b.Meta("name", "content", &MetaOptions{HttpEquiv: "equiv", Charset: "charset"})

	// Test hidden
	div = b.Div()
	div.T("foo")
	div.Span().T("span")
	div.Hide(true)
	b.Div().Hide(true)
	b.T(`><&"`).TUnsafe(`>`)
	b.TUnsafe(`><&"`)

	env := Environment{
		"name":   StringValue("foo"),
		"escape": StringValue("a>b"),
	}
	if err := compareHtml(root, kCompare, true, env); err != nil {
		t.Error(err)
	}
}

func Test_Simple(t *testing.T) {
	const kCompare = `<!DOCTYPE html><html><head><link rel="search" /><link rel="stylesheet" /></head></html>`

	root := NewRoot()
	head := root.Head()
	head.Link(LinkRelSearch)
	head.Link(LinkRelStylesheet)

	if err := compareHtml(root, kCompare, false); err != nil {
		t.Error(err)
	}
}

func Test_Remove(t *testing.T) {
	const kCompare = `<!DOCTYPE html><html><head></head><body><div class="foo2"></div></body></html>`
	root := NewRoot()
	root.Head()
	body := root.Body()
	div1 := body.DivClasses("foo1")
	body.DivClasses("foo2")
	div3 := body.DivClasses("foo3")
	div3.RemoveParent()
	div1.RemoveParent()

	if err := compareHtml(root, kCompare, false); err != nil {
		t.Error(err)
	}

}

func Test_TextTag(t *testing.T) {
	const kCompare = `<!DOCTYPE html>
<html>
  <body>
    hello 
    <a href="foo">
      world
    </a>
    .
    
    <abbr title="hello">
      H.W
    </abbr>
    
    <b>
      bold
    </b>
    
    <cite>
      cite
    </cite>
    
    <dfn>
      dfn
    </dfn>
    
    <em>
      em
    </em>
    
    <i>
      i
    </i>
    
    <kbd>
      kbd
    </kbd>
    
    <pre>
      pre
    </pre>
    
    <samp>
      samp
    </samp>
    
    <small>
      small
    </small>
    
    <strong>
      strong
    </strong>
    
    <u>
      u
    </u>
    
    <var>
      var
    </var>
    
    <img alt="img" src="img.jpg" />
    
    helloworld
    <span>
      span
    </span>
    foo
    <input checked="" type="checkbox" />
  </body>
</html>`

	root := NewRoot()
	body := root.Body()
	body.T("hello ").A("foo", "world").T(".")
	body.T().Abbr("H.W", "hello").B("bold").Cite("cite").Dfn("dfn").Em("em").I("i").Kbd("kbd").Pre("pre").Samp("samp").Small("small").Strong("strong").U("u").Var("var").Img("img.jpg", "img")
	body.T("hello").T("world").Up().Span().T("span").Up(2).T(`foo`)
	body.CheckedInput(CheckedInputTypeCheckbox).SetChecked(true)

	if err := compareHtml(root, kCompare, true); err != nil {
		t.Error(err)
	}
}

// Adapted from:
// http://genshi.edgewall.org/browser/trunk/examples/bench/bigtable.py
func Benchmark_Bigtable(b *testing.B) {
	// Initialization is performed in the loop to keep the data cache cold.
	for nn := 0; nn < b.N; nn++ {
		b.StopTimer()

		// Set up the input data.
		data := [1000][10]int{}
		for ii := 0; ii < len(data); ii++ {
			for jj := 0; jj < len(data[ii]); jj++ {
				data[ii][jj] = jj
			}
		}

		b.StartTimer()

		// Set up the template.
		tmpl := H.Table()
		for ii := 0; ii < len(data); ii++ {
			tr := tmpl.Tr()
			row := data[ii]
			for jj := 0; jj < len(row); jj++ {
				tr.Td().T(strconv.Itoa(row[jj]))
			}
		}

		// Render.
		buf := new(bytes.Buffer)
		if _, err := Write(buf, tmpl); err != nil {
			b.Error(err)
		}
		b.StopTimer()
	}
}

// This variant attemts to minimize tag construction by copying.
func Benchmark_BigtableFast(b *testing.B) {
	// Initialization is performed in the loop to keep the data cache cold.
	for nn := 0; nn < b.N; nn++ {
		b.StopTimer()

		// Set up the input data.
		data := [1000][10]int{}
		for ii := 0; ii < len(data); ii++ {
			for jj := 0; jj < len(data[ii]); jj++ {
				data[ii][jj] = jj
			}
		}

		b.StartTimer()

		// Set up the template.
		tmpl := H.Table()

		// We copy base tags, so that initial attribute processing need
		// only happen once.
		trBase := H.Tr()
		tdBase := H.Td()

		for ii := 0; ii < len(data); ii++ {
			tr := trBase.Copy()
			tmpl.AddChild(tr)

			for jj := 0; jj < len(data[ii]); jj++ {
				td := tdBase.Copy()
				tr.AddChild(td)
				td.T(strconv.Itoa(data[ii][jj]))
			}
		}

		// Render.
		buf := new(bytes.Buffer)
		if _, err := Write(buf, tmpl); err != nil {
			b.Error(err)
		}
		b.StopTimer()
	}
}

// Returns the position of the first differing character.  Otherwise, ok will
// be set to true if equal.
func stringCmp(a, b string) (pos int, err error) {
	count := len(a)
	if len(a) > len(b) {
		count = len(b)
	}
	for ii := 0; ii < count; ii++ {
		if a[ii] != b[ii] {
			lineA := getLine(a, ii)
			lineB := getLine(b, ii)
			err = errors.New(fmt.Sprintf("mismatch at (line,pos): (%d, %d) %q != (%d, %d) %q",
				lineA, ii, a[ii], lineB, ii, b[ii]))
			return ii, err
		}
	}
	return 0, nil
}

func compare(t *testing.T, root Tag, compare string) error {
	buf := new(bytes.Buffer)
	// Use pretty print to order attributes deterministically.
	if _, err := WritePretty(buf, root); err != nil {
		return err
	}
	if pos, err := stringCmp(buf.String(), compare); err != nil {
		t.Log(err)
		t.Log("\nOutput:\n", buf.String())
		t.Log("\nComparison:\n", compare)
		return errors.New(buf.String()[pos:])
	}
	return nil
}

// Returns the line number for the given position in a.
func getLine(a string, pos int) (lineno int) {
	lineno = 1
	for ii, c := range a {
		if ii == pos {
			return
		}
		if c == '\n' {
			lineno++
		}
	}
	panic("bad position")
}

func Benchmark_Text(b *testing.B) {
	for ii := 0; ii < b.N; ii++ {
		b := Factory.Body()
		for jj := 0; jj < 1000; jj++ {
			b.T("hello")
		}
	}
}

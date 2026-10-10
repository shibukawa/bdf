// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlro

// tokenTests are the cases of the tokenizer tests of golang.org/x/net/html,
// whose inputs the reader is run on beside that tokenizer. golden is what
// that package expects and is kept for reference.
// https://github.com/golang/go/issues/58246
const issue58246 = `<!--[if gte mso 12]>
  <xml>
      <o:OfficeDocumentSettings>
      <o:AllowPNG/>
      <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
<![endif]-->`

type tokenTest struct {
	// A short description of the test case.
	desc string
	// The HTML to parse.
	html string
	// The string representations of the expected tokens, joined by '$'.
	golden string
}

var tokenTests = []tokenTest{
	{
		"empty",
		"",
		"",
	},
	// A single text node. The tokenizer should not break text nodes on whitespace,
	// nor should it normalize whitespace within a text node.
	{
		"text",
		"foo  bar",
		"foo  bar",
	},
	// An entity.
	{
		"entity",
		"one &lt; two",
		"one &lt; two",
	},
	// A start, self-closing and end tag. The tokenizer does not care if the start
	// and end tokens don't match; that is the job of the parser.
	{
		"tags",
		"<a>b<c/>d</e>",
		"<a>$b$<c/>$d$</e>",
	},
	// Angle brackets that aren't a tag.
	{
		"not a tag #0",
		"<",
		"&lt;",
	},
	{
		"not a tag #1",
		"</",
		"&lt;/",
	},
	{
		"not a tag #2",
		"</>",
		"<!---->",
	},
	{
		"not a tag #3",
		"a</>b",
		"a$<!---->$b",
	},
	{
		"not a tag #4",
		"</ >",
		"<!-- -->",
	},
	{
		"not a tag #5",
		"</.",
		"<!--.-->",
	},
	{
		"not a tag #6",
		"</.>",
		"<!--.-->",
	},
	{
		"not a tag #7",
		"a < b",
		"a &lt; b",
	},
	{
		"not a tag #8",
		"<.>",
		"&lt;.&gt;",
	},
	{
		"not a tag #9",
		"a<<<b>>>c",
		"a&lt;&lt;$<b>$&gt;&gt;c",
	},
	{
		"not a tag #10",
		"if x<0 and y < 0 then x*y>0",
		"if x&lt;0 and y &lt; 0 then x*y&gt;0",
	},
	{
		"not a tag #11",
		"<<p>",
		"&lt;$<p>",
	},
	// EOF in a tag name.
	{
		"tag name eof #0",
		"<a",
		"",
	},
	{
		"tag name eof #1",
		"<a ",
		"",
	},
	{
		"tag name eof #2",
		"a<b",
		"a",
	},
	{
		"tag name eof #3",
		"<a><b",
		"<a>",
	},
	{
		"tag name eof #4",
		`<a x`,
		``,
	},
	// Some malformed tags that are missing a '>'.
	{
		"malformed tag #0",
		`<p</p>`,
		`<p< p="">`,
	},
	{
		"malformed tag #1",
		`<p </p>`,
		`<p <="" p="">`,
	},
	{
		"malformed tag #2",
		`<p id`,
		``,
	},
	{
		"malformed tag #3",
		`<p id=`,
		``,
	},
	{
		"malformed tag #4",
		`<p id=>`,
		`<p id="">`,
	},
	{
		"malformed tag #5",
		`<p id=0`,
		``,
	},
	{
		"malformed tag #6",
		`<p id=0</p>`,
		`<p id="0&lt;/p">`,
	},
	{
		"malformed tag #7",
		`<p id="0</p>`,
		``,
	},
	{
		"malformed tag #8",
		`<p id="0"</p>`,
		`<p id="0" <="" p="">`,
	},
	{
		"malformed tag #9",
		`<p></p id`,
		`<p>`,
	},
	// Raw text and RCDATA.
	{
		"basic raw text",
		"<script><a></b></script>",
		"<script>$&lt;a&gt;&lt;/b&gt;$</script>",
	},
	{
		"unfinished script end tag",
		"<SCRIPT>a</SCR",
		"<script>$a&lt;/SCR",
	},
	{
		"broken script end tag",
		"<SCRIPT>a</SCR ipt>",
		"<script>$a&lt;/SCR ipt&gt;",
	},
	{
		"EOF in script end tag",
		"<SCRIPT>a</SCRipt",
		"<script>$a&lt;/SCRipt",
	},
	{
		"scriptx end tag",
		"<SCRIPT>a</SCRiptx",
		"<script>$a&lt;/SCRiptx",
	},
	{
		"' ' completes script end tag",
		"<SCRIPT>a</SCRipt ",
		"<script>$a",
	},
	{
		"'>' completes script end tag",
		"<SCRIPT>a</SCRipt>",
		"<script>$a$</script>",
	},
	{
		"self-closing script end tag",
		"<SCRIPT>a</SCRipt/>",
		"<script>$a$</script>",
	},
	{
		"nested script tag",
		"<SCRIPT>a</SCRipt<script>",
		"<script>$a&lt;/SCRipt&lt;script&gt;",
	},
	{
		"script end tag after unfinished",
		"<SCRIPT>a</SCRipt</script>",
		"<script>$a&lt;/SCRipt$</script>",
	},
	{
		"script/style mismatched tags",
		"<script>a</style>",
		"<script>$a&lt;/style&gt;",
	},
	{
		"style element with entity",
		"<style>&apos;",
		"<style>$&amp;apos;",
	},
	{
		"textarea with tag",
		"<textarea><div></textarea>",
		"<textarea>$&lt;div&gt;$</textarea>",
	},
	{
		"title with tag and entity",
		"<title><b>K&amp;R C</b></title>",
		"<title>$&lt;b&gt;K&amp;R C&lt;/b&gt;$</title>",
	},
	{
		"title with trailing '&lt;' entity",
		"<title>foobar<</title>",
		"<title>$foobar&lt;$</title>",
	},
	// DOCTYPE tests.
	{
		"Proper DOCTYPE",
		"<!DOCTYPE html>",
		"<!DOCTYPE html>",
	},
	{
		"DOCTYPE with no space",
		"<!doctypehtml>",
		"<!DOCTYPE html>",
	},
	{
		"DOCTYPE with two spaces",
		"<!doctype  html>",
		"<!DOCTYPE html>",
	},
	{
		"looks like DOCTYPE but isn't",
		"<!DOCUMENT html>",
		"<!--DOCUMENT html-->",
	},
	{
		"DOCTYPE at EOF",
		"<!DOCtype",
		"<!DOCTYPE >",
	},
	// XML processing instructions.
	{
		"XML processing instruction",
		"<?xml?>",
		"<!--?xml?-->",
	},
	// Comments. See also func TestComments.
	{
		"comment0",
		"abc<b><!-- skipme --></b>def",
		"abc$<b>$<!-- skipme -->$</b>$def",
	},
	{
		"comment1",
		"a<!-->z",
		"a$<!---->$z",
	},
	{
		"comment2",
		"a<!--->z",
		"a$<!---->$z",
	},
	{
		"comment3",
		"a<!--x>-->z",
		"a$<!--x>-->$z",
	},
	{
		"comment4",
		"a<!--x->-->z",
		"a$<!--x-&gt;-->$z",
	},
	{
		"comment5",
		"a<!>z",
		"a$<!---->$z",
	},
	{
		"comment6",
		"a<!->z",
		"a$<!----->$z",
	},
	{
		"comment7",
		"a<!---<>z",
		"a$<!---<>z-->",
	},
	{
		"comment8",
		"a<!--z",
		"a$<!--z-->",
	},
	{
		"comment9",
		"a<!--z-",
		"a$<!--z-->",
	},
	{
		"comment10",
		"a<!--z--",
		"a$<!--z-->",
	},
	{
		"comment11",
		"a<!--z---",
		"a$<!--z--->",
	},
	{
		"comment12",
		"a<!--z----",
		"a$<!--z---->",
	},
	{
		"comment13",
		"a<!--x--!>z",
		"a$<!--x-->$z",
	},
	{
		"comment14",
		"a<!--!-->z",
		"a$<!--!-->$z",
	},
	{
		"comment15",
		"a<!-- !-->z",
		"a$<!-- !-->$z",
	},
	{
		"comment16",
		"a<!--i\x00j-->z",
		"a$<!--i\uFFFDj-->$z",
	},
	{
		"comment17",
		"a<!--\x00",
		"a$<!--\uFFFD-->",
	},
	{
		"comment18",
		"a<!--<!-->z",
		"a$<!--<!-->$z",
	},
	{
		"comment19",
		"a<!--<!--",
		"a$<!--<!-->",
	},
	{
		"comment20",
		"a<!--ij--kl-->z",
		"a$<!--ij--kl-->$z",
	},
	{
		"comment21",
		"a<!--ij--kl--!>z",
		"a$<!--ij--kl-->$z",
	},
	{
		"comment22",
		"a<!--!--!<--!-->z",
		"a$<!--!--!<--!-->$z",
	},
	{
		"comment23",
		"a<!--&gt;-->z",
		"a$<!--&gt;-->$z",
	},
	{
		"comment24",
		"a<!--&gt;>x",
		"a$<!--&gt;>x-->",
	},
	{
		"comment25",
		"a<!--&gt;&gt;",
		"a$<!--&gt;>-->",
	},
	{
		"comment26",
		"a<!--&gt;&gt;-",
		"a$<!--&gt;>-->",
	},
	{
		"comment27",
		"a<!--&gt;&gt;-->z",
		"a$<!--&gt;>-->$z",
	},
	{
		"comment28",
		"a<!--&amp;&gt;-->z",
		"a$<!--&amp;>-->$z",
	},
	{
		"comment29",
		"a<!--&amp;gt;-->z",
		"a$<!--&amp;gt;-->$z",
	},
	{
		"comment30",
		"a<!--&nosuchentity;-->z",
		"a$<!--&amp;nosuchentity;-->$z",
	},
	{
		"comment31",
		"a<!--i>>j-->z",
		"a$<!--i>>j-->$z",
	},
	{
		"comment32",
		"a<!--i!>>j-->z",
		"a$<!--i!&gt;>j-->$z",
	},
	// https://stackoverflow.design/email/base/mso/#targeting-specific-outlook-versions
	// says "[For] Windows Outlook 2003 and above... conditional comments allow
	// us to add bits of HTML that are only read by the Word-based versions of
	// Outlook". These comments (with angle brackets) should pass through
	// unchanged (by this Go package) when rendering.
	//
	// We should also still escape ">" as "&gt;" when necessary.
	// https://github.com/golang/go/issues/48237
	//
	// The "your code" example below comes from that stackoverflow.design link
	// above but note that it can contain angle-bracket-rich XML.
	// https://github.com/golang/go/issues/58246
	{
		"issue48237CommentWithAmpgtsemi1",
		"a<!--<p></p>&lt;!--[video]--&gt;-->z",
		"a$<!--<p></p><!--[video]--&gt;-->$z",
	},
	{
		"issue48237CommentWithAmpgtsemi2",
		"a<!--<p></p>&lt;!--[video]--!&gt;-->z",
		"a$<!--<p></p><!--[video]--!&gt;-->$z",
	},
	{
		"issue58246MicrosoftOutlookComment1",
		"a<!--[if mso]> your code <![endif]-->z",
		"a$<!--[if mso]> your code <![endif]-->$z",
	},
	{
		"issue58246MicrosoftOutlookComment2",
		"a" + issue58246 + "z",
		"a$" + issue58246 + "$z",
	},
	// An attribute with a backslash.
	{
		"backslash",
		`<p id="a\"b">`,
		`<p id="a\" b"="">`,
	},
	// Entities, tag name and attribute key lower-casing, and whitespace
	// normalization within a tag.
	{
		"tricky",
		"<p \t\n iD=\"a&quot;B\"  foo=\"bar\"><EM>te&lt;&amp;;xt</em></p>",
		`<p id="a&#34;B" foo="bar">$<em>$te&lt;&amp;;xt$</em>$</p>`,
	},
	// A nonexistent entity. Tokenizing and converting back to a string should
	// escape the "&" to become "&amp;".
	{
		"noSuchEntity",
		`<a b="c&noSuchEntity;d">&lt;&alsoDoesntExist;&`,
		`<a b="c&amp;noSuchEntity;d">$&lt;&amp;alsoDoesntExist;&amp;`,
	},
	{
		"entity without semicolon",
		`&notit;&notin;<a b="q=z&amp=5&notice=hello&not;=world">`,
		`¬it;∉$<a b="q=z&amp;amp=5&amp;notice=hello¬=world">`,
	},
	{
		"entity with digits",
		"&frac12;",
		"½",
	},
	// Attribute tests:
	// http://dev.w3.org/html5/pf-summary/Overview.html#attributes
	{
		"Empty attribute",
		`<input disabled FOO>`,
		`<input disabled="" foo="">`,
	},
	{
		"Empty attribute, whitespace",
		`<input disabled FOO >`,
		`<input disabled="" foo="">`,
	},
	{
		"Unquoted attribute value",
		`<input value=yes FOO=BAR>`,
		`<input value="yes" foo="BAR">`,
	},
	{
		"Unquoted attribute value, spaces",
		`<input value = yes FOO = BAR>`,
		`<input value="yes" foo="BAR">`,
	},
	{
		"Unquoted attribute value, trailing space",
		`<input value=yes FOO=BAR >`,
		`<input value="yes" foo="BAR">`,
	},
	{
		"Single-quoted attribute value",
		`<input value='yes' FOO='BAR'>`,
		`<input value="yes" foo="BAR">`,
	},
	{
		"Single-quoted attribute value, trailing space",
		`<input value='yes' FOO='BAR' >`,
		`<input value="yes" foo="BAR">`,
	},
	{
		"Double-quoted attribute value",
		`<input value="I'm an attribute" FOO="BAR">`,
		`<input value="I&#39;m an attribute" foo="BAR">`,
	},
	{
		"Attribute name characters",
		`<meta http-equiv="content-type">`,
		`<meta http-equiv="content-type">`,
	},
	{
		"Mixed attributes",
		`a<P V="0 1" w='2' X=3 y>z`,
		`a$<p v="0 1" w="2" x="3" y="">$z`,
	},
	{
		"Attributes with a solitary single quote",
		`<p id=can't><p id=won't>`,
		`<p id="can&#39;t">$<p id="won&#39;t">`,
	},
	// WHATWG 13.2.5.32 equals sign before attribute name state
	{
		"equals sign before attribute name",
		`<p  =>`,
		`<p =="">`,
	},
	{
		"equals sign before attribute name, extra cruft",
		`<p  =asd>`,
		`<p =asd="">`,
	},
	{
		"forward slash before attribute name",
		`<p/=">`,
		`<p ="="">`,
	},
	{
		"forward slash before attribute name with spaces around",
		`<p / =">`,
		`<p ="="">`,
	},
	{
		"forward slash after attribute name followed by a character",
		`<p a/ ="">`,
		`<p a="" =""="">`,
	},
	{
		"slash at end of unquoted attribute value",
		`<p a="\">`,
		`<p a="\">`,
	},
	{
		"self-closing tag with attribute",
		`<p a=/>`,
		`<p a="/">`,
	},
	{
		"duplicate attributes",
		`<p foo="bar" foo="baz">`,
		`<p foo="bar">`,
	},
	{
		"duplicate attributes, different case",
		`<p FOO="bar" foo="baz">`,
		`<p foo="bar">`,
	},
	{
		"partial doctype",
		`<!doc`,
		`<!--doc-->`,
	},
	{
		"partial cdata",
		`<![CDA`,
		`<!--[CDA-->`,
	},
	{
		"partial comment",
		`<!comment`,
		`<!--comment-->`,
	},
}

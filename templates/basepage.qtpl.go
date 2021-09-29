// Code generated by qtc from "basepage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// This is a base page template. All the other template pages implement this interface.
//

//line templates/basepage.qtpl:3
package templates

//line templates/basepage.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/basepage.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/basepage.qtpl:4
type Page interface {
//line templates/basepage.qtpl:4
	Title() string
//line templates/basepage.qtpl:4
	StreamTitle(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:4
	WriteTitle(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:4
	Head() string
//line templates/basepage.qtpl:4
	StreamHead(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:4
	WriteHead(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:4
	Body() string
//line templates/basepage.qtpl:4
	StreamBody(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:4
	WriteBody(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:4
}

// Page prints a page implementing Page interface.

//line templates/basepage.qtpl:12
func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
//line templates/basepage.qtpl:12
	qw422016.N().S(`
<html lang="en">
    <head>
        <title>`)
//line templates/basepage.qtpl:15
	p.StreamTitle(qw422016)
//line templates/basepage.qtpl:15
	qw422016.N().S(`</title>
        `)
//line templates/basepage.qtpl:16
	p.StreamHead(qw422016)
//line templates/basepage.qtpl:16
	qw422016.N().S(`
        <meta property="og:image" content="/favicon.png">
        <link rel="icon" type="image/x-icon" href="/favicon.ico"/>
        <link rel="shortcut icon" type="image/x-icon" href="/favicon.ico"/>

        <link rel="stylesheet" href="https://unpkg.com/purecss@2.0.6/build/grids-responsive-min.css">
        <link rel="stylesheet" href="/css/style.css">
        <link rel="stylesheet" href="/css/style.large.css">
    </head>
    <body>
        `)
//line templates/basepage.qtpl:26
	p.StreamBody(qw422016)
//line templates/basepage.qtpl:26
	qw422016.N().S(`
    </body>
</html>
`)
//line templates/basepage.qtpl:29
}

//line templates/basepage.qtpl:29
func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
//line templates/basepage.qtpl:29
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:29
	StreamPageTemplate(qw422016, p)
//line templates/basepage.qtpl:29
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:29
}

//line templates/basepage.qtpl:29
func PageTemplate(p Page) string {
//line templates/basepage.qtpl:29
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:29
	WritePageTemplate(qb422016, p)
//line templates/basepage.qtpl:29
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:29
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:29
	return qs422016
//line templates/basepage.qtpl:29
}

// Base page implementation. Other pages may inherit from it if they need
// overriding only certain Page methods

//line templates/basepage.qtpl:33
type BasePage struct{}

//line templates/basepage.qtpl:34
func (p *BasePage) StreamTitle(qw422016 *qt422016.Writer) {
//line templates/basepage.qtpl:34
	qw422016.N().S(`Default title`)
//line templates/basepage.qtpl:34
}

//line templates/basepage.qtpl:34
func (p *BasePage) WriteTitle(qq422016 qtio422016.Writer) {
//line templates/basepage.qtpl:34
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:34
	p.StreamTitle(qw422016)
//line templates/basepage.qtpl:34
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:34
}

//line templates/basepage.qtpl:34
func (p *BasePage) Title() string {
//line templates/basepage.qtpl:34
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:34
	p.WriteTitle(qb422016)
//line templates/basepage.qtpl:34
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:34
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:34
	return qs422016
//line templates/basepage.qtpl:34
}

//line templates/basepage.qtpl:35
func (p *BasePage) StreamHead(qw422016 *qt422016.Writer) {
//line templates/basepage.qtpl:35
}

//line templates/basepage.qtpl:35
func (p *BasePage) WriteHead(qq422016 qtio422016.Writer) {
//line templates/basepage.qtpl:35
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:35
	p.StreamHead(qw422016)
//line templates/basepage.qtpl:35
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:35
}

//line templates/basepage.qtpl:35
func (p *BasePage) Head() string {
//line templates/basepage.qtpl:35
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:35
	p.WriteHead(qb422016)
//line templates/basepage.qtpl:35
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:35
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:35
	return qs422016
//line templates/basepage.qtpl:35
}

//line templates/basepage.qtpl:36
func (p *BasePage) StreamBody(qw422016 *qt422016.Writer) {
//line templates/basepage.qtpl:36
}

//line templates/basepage.qtpl:36
func (p *BasePage) WriteBody(qq422016 qtio422016.Writer) {
//line templates/basepage.qtpl:36
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:36
	p.StreamBody(qw422016)
//line templates/basepage.qtpl:36
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:36
}

//line templates/basepage.qtpl:36
func (p *BasePage) Body() string {
//line templates/basepage.qtpl:36
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:36
	p.WriteBody(qb422016)
//line templates/basepage.qtpl:36
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:36
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:36
	return qs422016
//line templates/basepage.qtpl:36
}

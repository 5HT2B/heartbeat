// Code generated by qtc from "mainpage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/mainpage.qtpl:1
package templates

//line templates/mainpage.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/mainpage.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/mainpage.qtpl:2
type MainPage struct {
	LastSeen       string
	TimeDifference string
	MissingBeat    string
	TotalBeats     string
	CurrentTime    string
	GitHash        string
	GitRepo        string
	ServerName     string
}

//line templates/mainpage.qtpl:14
func (p *MainPage) StreamTitle(qw422016 *qt422016.Writer) {
//line templates/mainpage.qtpl:14
	qw422016.N().S(`
    `)
//line templates/mainpage.qtpl:15
	qw422016.E().S(p.ServerName)
//line templates/mainpage.qtpl:15
	qw422016.N().S(`
`)
//line templates/mainpage.qtpl:16
}

//line templates/mainpage.qtpl:16
func (p *MainPage) WriteTitle(qq422016 qtio422016.Writer) {
//line templates/mainpage.qtpl:16
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/mainpage.qtpl:16
	p.StreamTitle(qw422016)
//line templates/mainpage.qtpl:16
	qt422016.ReleaseWriter(qw422016)
//line templates/mainpage.qtpl:16
}

//line templates/mainpage.qtpl:16
func (p *MainPage) Title() string {
//line templates/mainpage.qtpl:16
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/mainpage.qtpl:16
	p.WriteTitle(qb422016)
//line templates/mainpage.qtpl:16
	qs422016 := string(qb422016.B)
//line templates/mainpage.qtpl:16
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/mainpage.qtpl:16
	return qs422016
//line templates/mainpage.qtpl:16
}

//line templates/mainpage.qtpl:18
func (p *MainPage) StreamHead(qw422016 *qt422016.Writer) {
//line templates/mainpage.qtpl:18
	qw422016.N().S(`
<meta property="og:site_name" content="`)
//line templates/mainpage.qtpl:19
	qw422016.E().S(p.ServerName)
//line templates/mainpage.qtpl:19
	qw422016.N().S(`">
<meta property="og:description" content="Last seen at: `)
//line templates/mainpage.qtpl:20
	qw422016.E().S(p.LastSeen)
//line templates/mainpage.qtpl:20
	qw422016.N().S(`.
This embed was generated at `)
//line templates/mainpage.qtpl:21
	qw422016.E().S(p.CurrentTime)
//line templates/mainpage.qtpl:21
	qw422016.N().S(`.
Due to caching, you will have to check the website if the embed generation time is old."/>
<meta name="theme-color" content="#6495ED">
<script>
    window.onload = function () {
        setInterval(updateInfo, 1000)
    };

    async function updateInfo() {
        let response = await fetch("/api/info?countVisit=false");
        let data = await response.json();

        await setInfo("LastSeen", data.last_seen, "Last response time")
        await setInfo("TimeDifference", data.time_difference, "Time since last response")
        await setInfo("MissingBeat", data.missing_beat, "Longest recorded absence")
        await setInfo("TotalBeats", data.total_beats, "Total beats received")
    }

    async function setInfo(id, json, prefix) {
        document.getElementById(id).innerHTML = `)
//line templates/mainpage.qtpl:21
	qw422016.N().S("`")
//line templates/mainpage.qtpl:21
	qw422016.N().S(`${prefix}:<br>${json}`)
//line templates/mainpage.qtpl:21
	qw422016.N().S("`")
//line templates/mainpage.qtpl:21
	qw422016.N().S(`
    }
</script>
`)
//line templates/mainpage.qtpl:43
}

//line templates/mainpage.qtpl:43
func (p *MainPage) WriteHead(qq422016 qtio422016.Writer) {
//line templates/mainpage.qtpl:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/mainpage.qtpl:43
	p.StreamHead(qw422016)
//line templates/mainpage.qtpl:43
	qt422016.ReleaseWriter(qw422016)
//line templates/mainpage.qtpl:43
}

//line templates/mainpage.qtpl:43
func (p *MainPage) Head() string {
//line templates/mainpage.qtpl:43
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/mainpage.qtpl:43
	p.WriteHead(qb422016)
//line templates/mainpage.qtpl:43
	qs422016 := string(qb422016.B)
//line templates/mainpage.qtpl:43
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/mainpage.qtpl:43
	return qs422016
//line templates/mainpage.qtpl:43
}

//line templates/mainpage.qtpl:45
func (p *MainPage) StreamBody(qw422016 *qt422016.Writer) {
//line templates/mainpage.qtpl:45
	qw422016.N().S(`
    <div class="spacer"></div>
    <div class="pure-g preamble">
        <div class="pure-g-u-0 pure-u-lg-1-6"></div>
        <div class="pure-u-1 pure-u-lg-4-6">
            <p class="center">
                Welcome to `)
//line templates/mainpage.qtpl:51
	qw422016.E().S(p.ServerName)
//line templates/mainpage.qtpl:51
	qw422016.N().S(`. <br>
                This page displays the last timestamp that they have unlocked and used any of their devices. <br>
                If they have been absent for more than 48 hours, something is probably wrong. <br>
                This website is running on version <a href="`)
//line templates/mainpage.qtpl:54
	qw422016.E().S(p.GitRepo)
//line templates/mainpage.qtpl:54
	qw422016.N().S(`/tree/`)
//line templates/mainpage.qtpl:54
	qw422016.E().S(p.GitHash)
//line templates/mainpage.qtpl:54
	qw422016.N().S(`"><code>`)
//line templates/mainpage.qtpl:54
	qw422016.E().S(p.GitHash)
//line templates/mainpage.qtpl:54
	qw422016.N().S(`</code></a> of <a href="`)
//line templates/mainpage.qtpl:54
	qw422016.E().S(p.GitRepo)
//line templates/mainpage.qtpl:54
	qw422016.N().S(`">Heartbeat</a>.
            </p>
        </div>
        <div class="pure-g-u-0 pure-u-lg-1-6"></div>
    </div>
    <div class="pure-g times">
        <div class="pure-u-0 pure-u-lg-1-6"></div>
        <div class="pure-u-1 pure-u-lg-1-6">
            <p class="center" id="LastSeen">Last response time:<br>`)
//line templates/mainpage.qtpl:62
	qw422016.E().S(p.LastSeen)
//line templates/mainpage.qtpl:62
	qw422016.N().S(`</p>
        </div>
        <div class="pure-u-1 pure-u-lg-1-6">
            <p class="center" id="TimeDifference">Time since last response:<br>`)
//line templates/mainpage.qtpl:65
	qw422016.E().S(p.TimeDifference)
//line templates/mainpage.qtpl:65
	qw422016.N().S(`</p>
        </div>
        <div class="pure-u-1 pure-u-lg-1-6">
            <p class="center" id="MissingBeat">Longest recorded absence:<br>`)
//line templates/mainpage.qtpl:68
	qw422016.E().S(p.MissingBeat)
//line templates/mainpage.qtpl:68
	qw422016.N().S(`</p>
        </div>
        <div class="pure-u-1 pure-u-lg-1-6">
            <p class="center" id="TotalBeats">Total beats received:<br>`)
//line templates/mainpage.qtpl:71
	qw422016.E().S(p.TotalBeats)
//line templates/mainpage.qtpl:71
	qw422016.N().S(`</p>
        </div>
        <div class="pure-u-0 pure-u-lg-1-6"></div>
    </div>
    <div class="spacer"></div>
    <div class="pure-g links">
        <div class="pure-g-u-0 pure-u-lg-1-6"></div>
        <div class="pure-u-1 pure-u-lg-4-6">
            <p class="center">
                <a href="/stats">Stats</a> - <a href="/privacy">Privacy Policy</a>
            </p>
        </div>
        <div class="pure-g-u-0 pure-u-lg-1-6"></div>
    </div>
`)
//line templates/mainpage.qtpl:85
}

//line templates/mainpage.qtpl:85
func (p *MainPage) WriteBody(qq422016 qtio422016.Writer) {
//line templates/mainpage.qtpl:85
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/mainpage.qtpl:85
	p.StreamBody(qw422016)
//line templates/mainpage.qtpl:85
	qt422016.ReleaseWriter(qw422016)
//line templates/mainpage.qtpl:85
}

//line templates/mainpage.qtpl:85
func (p *MainPage) Body() string {
//line templates/mainpage.qtpl:85
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/mainpage.qtpl:85
	p.WriteBody(qb422016)
//line templates/mainpage.qtpl:85
	qs422016 := string(qb422016.B)
//line templates/mainpage.qtpl:85
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/mainpage.qtpl:85
	return qs422016
//line templates/mainpage.qtpl:85
}

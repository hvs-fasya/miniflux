// Code generated by go generate; DO NOT EDIT.
// 2018-07-25 17:21:53.065889348 +0300 MSK m=+0.007307151

package template

var templateNewsViewsMap = map[string]string{
	"news_home": `{{ define "content"}}
<!-- OFFICIAL NEWS -->
<div class="mdl-layout__tab-panel is-active" id="official">
</div>
<!-- MEDIA NEWS -->
{{/*{{ template "news_media" . }}*/}}
<div class="mdl-layout__tab-panel" id="media">
</div>
<!-- TRAVEL ALERTS -->
<div class="mdl-layout__tab-panel" id="travel">
</div>
<!-- SECURITY ALERTS -->
<div class="mdl-layout__tab-panel" id="security">
    <h2 style="color:#666565; text-align:center; opacity: 0.87;">Real-time Global Security alerts</h2>
    <p style="color: #dcdbda; text-align:center; font-size:18px;">Data source: US, Canada, EU government bodies</p>
    <br/>
    <div id="map" style="position: relative; width: 95%; height: 70%;"></div>
    {{/*<div id="map" style="position: relative; width: 900px; height: 300px; padding-left: 10%;"></div>*/}}
</div>
<!-- VISA CHANGES -->
<div class="mdl-layout__tab-panel" id="visa">
</div>

{{ end }}`,
	"news_media": `{{ define "news_media"}}

    <br>
    {{ if eq .countrytotal 0 }}
    <div class="mdl-grid boxes-center">
    <h4>No recent news</h4>
    </div>
    {{ end }}
{{ range $i, $e := .mediaentries }}

{{ if eq $i 2 }}
    {{/*<div class="mdl-grid boxes-center">*/}}
        {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
            {{/*<div class="mdl-card__title mdl-card--expand">*/}}
                {{/*<h2 class="mdl-card__title-text">VISA UPDATES</h2>*/}}
            {{/*</div>*/}}
            {{/*<div class="mdl-card__supporting-text">*/}}
                {{/*Get Visa Updates alerts*/}}
            {{/*</div>*/}}
            {{/*<div class="mdl-card__actions mdl-card--border">*/}}
                {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">*/}}
                    {{/*View Updates*/}}
                {{/*</a>*/}}
            {{/*</div>*/}}
        {{/*</div>*/}}
        {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
            {{/*<div class="mdl-card__title mdl-card--expand">*/}}
                {{/*<h2 class="mdl-card__title-text">SEARCH VISA</h2>*/}}
            {{/*</div>*/}}
            {{/*<div class="mdl-card__supporting-text">*/}}
                {{/*Search Global Visas*/}}
            {{/*</div>*/}}
            {{/*<div class="mdl-card__actions mdl-card--border">*/}}
                {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="http://visadb.io" target="_blank">*/}}
                    {{/*Search*/}}
                {{/*</a>*/}}
            {{/*</div>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
{{ end }}
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
                <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
                    <h4>{{ noescape $e.Title }}</h4>
                    <p class="mdl-list__item">
                        <span class="mdl-list__item-primary-content">
                            <i class="material-icons mdl-list__item-icon">schedule</i>
                        {{ elapsed "UTC" $e.Date }}
                            <span class="mdl-layout-spacer"></span>
                        {{ if gt $e.Feed.Icon.IconID 0 }}
                            <i class="mdl-list__item-icon"><img src="{{ route "feedicon" "iconID" $e.Feed.Icon.IconID }}" width="26" height="26"></i>
                        {{ end }}
                        {{ $e.Feed.Title }}
                        </span>
                    </p>
                    <a role="button" href="#" data-fetch-content-url="{{ route "newsFetchContent" "entryID" $e.ID }}">READ FULL STORY</a>
                </div>
            </div>
        </div>
    </section>
{{end}}
    <!-- PAGINATION -->
    <div class="mdl-grid boxes-center">
        {{ if gt .offset 0 }}
            <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev">
        {{ else }}
            <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev" disabled>
        {{ end }}
            <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_left</i> Prev
            </button>
            <div class="mdl-layout-spacer"></div>
        {{ if .hasNext }}
            <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next">
        {{ else}}
            <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next" disabled>
        {{ end }}
        Next <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_right</i>
            </button>
    </div>
{{ end }}`,
	"news_official": `{{ define "news_official"}}

<br>
{{ range $i, $e := .officialentries }}
{{ if eq $i 2 }}
{{/*<div class="mdl-grid boxes-center">*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">VISA UPDATES</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Get Visa Updates alerts*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">*/}}
                {{/*View Updates*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">SEARCH VISA</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Search Global Visas*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="http://visadb.io" target="_blank">*/}}
                {{/*Search*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
{{/*</div>*/}}
{{ end }}
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
    <div class="mdl-card mdl-cell mdl-cell--12-col">
        <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
            <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
                <h4>{{ noescape $e.Title }}</h4>
                <p class="mdl-list__item">
                        <span class="mdl-list__item-primary-content">
                            <i class="material-icons mdl-list__item-icon">schedule</i>
                        {{ elapsed "UTC" $e.Date }}
                            <span class="mdl-layout-spacer"></span>
                        {{ if gt $e.Feed.Icon.IconID 0 }}
                            <i class="mdl-list__item-icon"><img src="{{ route "feedicon" "iconID" $e.Feed.Icon.IconID }}" width="26" height="26"></i>
                        {{ end }}
                        {{ $e.Feed.Title }}
                        </span>
                </p>
                <a role="button" href="#" data-fetch-content-url="{{ route "newsFetchContent" "entryID" $e.ID }}">READ FULL STORY</a>
            </div>
        </div>
    </div>
</section>
{{end}}
<!-- PAGINATION -->
<div class="mdl-grid boxes-center">
    {{ if gt .officialoffset 0 }}
        <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev">
    {{ else }}
        <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev" disabled>
    {{ end }}
        <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_left</i> Prev
        </button>
        <div class="mdl-layout-spacer"></div>
    {{ if .officialHasNext }}
        <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next">
    {{ else}}
        <button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next" disabled>
    {{ end }}
        Next <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_right</i>
    </button>
</div>

{{ end }}`,
	"news_sources": `{{ define "content"}}

<div class="mdl-layout__tab-panel is-active" id="official">
<br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text">
                <div class="section__text">
                    <h4>Official News:</h4>
                    <h5>Source:</h5>
                    <div>
                        <p>Global Ministry of foreign Affairs - Global Immigration Department - Global Tourism Boards</p>
                    </div>
                </div>
            </div>
        </div>
    </section>
    {{ template "disclaimer" }}
</div>
<div class="mdl-layout__tab-panel" id="media">
    <br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text">
                <div class="section__text">
                    <h4> Media News:</h4>
                    <h5>Source:</h5>
                    <div>
                        <p>Mainstream media news outlet i.e. CNN, BBC and others and Media Migration Policy.</p>
                    </div>
                </div>
            </div>
        </div>
    </section>
{{ template "disclaimer" }}
</div>
<div class="mdl-layout__tab-panel" id="travel">
    <br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text">
                <div class="section__text">
                    <h4>Travel Alerts:</h4>
                    <h5>Source:</h5>
                    <div>
                        <p>CDC Yellow Fever and CDC Travel Alerts</p>
                    </div>
                </div>
            </div>
        </div>
    </section>
{{ template "disclaimer" }}
</div>
<div class="mdl-layout__tab-panel" id="security">
    <br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text">
                <div class="section__text">
                    <h4>Security Alerts:</h4>
                    <h5>Source:</h5>
                    <div>
                        <p></p>
                    </div>
                </div>
            </div>
        </div>
    </section>
{{ template "disclaimer" }}
</div>
<div class="mdl-layout__tab-panel" id="visa">
    <br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text">
                <div class="section__text">
                    <h4>Visa Changes:</h4>
                    <h5>Source:</h5>
                    <div>
                        <p></p>
                    </div>
                </div>
            </div>
        </div>
    </section>
{{ template "disclaimer" }}
</div>

{{ end }}`,
	"news_travel": `{{ define "news_travel"}}

<br>
{{ if eq .countrytotal 0 }}
<div class="mdl-grid boxes-center">
    <h4>No recent news</h4>
</div>
{{ end }}
{{ range $i, $e := .travelentries }}

{{ if eq $i 2 }}
{{/*<div class="mdl-grid boxes-center">*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">VISA UPDATES</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Get Visa Updates alerts*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">*/}}
                {{/*View Updates*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">SEARCH VISA</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Search Global Visas*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="http://visadb.io" target="_blank">*/}}
                {{/*Search*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
{{/*</div>*/}}
{{ end }}
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
    <div class="mdl-card mdl-cell mdl-cell--12-col">
        <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
            <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
                <h4>{{ noescape $e.Title }}</h4>
                <p class="mdl-list__item">
                        <span class="mdl-list__item-primary-content">
                            <i class="material-icons mdl-list__item-icon">schedule</i>
                        {{ elapsed "UTC" $e.Date }}
                            <span class="mdl-layout-spacer"></span>
                        {{ if gt $e.Feed.Icon.IconID 0 }}
                            <i class="mdl-list__item-icon"><img src="{{ route "feedicon" "iconID" $e.Feed.Icon.IconID }}" width="26" height="26"></i>
                        {{ end }}
                        {{ $e.Feed.Title }}
                        </span>
                </p>
                <a role="button" href="#" data-fetch-content-url="{{ route "newsFetchContent" "entryID" $e.ID }}">READ FULL STORY</a>
            </div>
        </div>
    </div>
</section>
{{end}}
<!-- PAGINATION -->
<div class="mdl-grid boxes-center">
{{ if gt .offset 0 }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev">
{{ else }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev" disabled>
{{ end }}
    <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_left</i> Prev
</button>
    <div class="mdl-layout-spacer"></div>
{{ if .hasNext }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next">
{{ else}}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next" disabled>
{{ end }}
    Next <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_right</i>
</button>
</div>
{{ end }}`,
	"news_visa": `{{ define "news_visa"}}

<br>
{{ if eq .visatotal 0 }}
<div class="mdl-grid boxes-center">
    <h4>No recent news</h4>
</div>
{{ end }}
{{ range $i, $e := .visaentries }}

{{ if eq $i 2 }}
{{/*<div class="mdl-grid boxes-center">*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">VISA UPDATES</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Get Visa Updates alerts*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">*/}}
                {{/*View Updates*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
    {{/*<div class="demo-card-square mdl-card mdl-shadow--2dp mdl-cell mdl-cell--6-col">*/}}
        {{/*<div class="mdl-card__title mdl-card--expand">*/}}
            {{/*<h2 class="mdl-card__title-text">SEARCH VISA</h2>*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__supporting-text">*/}}
            {{/*Search Global Visas*/}}
        {{/*</div>*/}}
        {{/*<div class="mdl-card__actions mdl-card--border">*/}}
            {{/*<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="http://visadb.io" target="_blank">*/}}
                {{/*Search*/}}
            {{/*</a>*/}}
        {{/*</div>*/}}
    {{/*</div>*/}}
{{/*</div>*/}}
{{ end }}
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
    <div class="mdl-card mdl-cell mdl-cell--12-col">
        <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
            <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
                <h4>{{ noescape $e.Title }}</h4>
                <p class="mdl-list__item">
                        <span class="mdl-list__item-primary-content">
                            <i class="material-icons mdl-list__item-icon">schedule</i>
                        {{ elapsed "UTC" $e.Date }}
                            <span class="mdl-layout-spacer"></span>
                        {{ if gt $e.Feed.Icon.IconID 0 }}
                            <i class="mdl-list__item-icon"><img src="{{ route "feedicon" "iconID" $e.Feed.Icon.IconID }}" width="26" height="26"></i>
                        {{ end }}
                        {{ $e.Feed.Title }}
                        </span>
                </p>
                <a role="button" href="#" data-fetch-content-url="{{ route "newsFetchContent" "entryID" $e.ID }}">READ FULL STORY</a>
            </div>
        </div>
    </div>
</section>
{{end}}
<!-- PAGINATION -->
<div class="mdl-grid boxes-center">
{{ if gt .offset 0 }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev">
{{ else }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="prev" disabled>
{{ end }}
    <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_left</i> Prev
</button>
    <div class="mdl-layout-spacer"></div>
{{ if .hasNext }}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next">
{{ else}}
<button class="mdl-button mdl-js-button mdl-js-ripple-effect" data-page="next" disabled>
{{ end }}
    Next <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_right</i>
</button>
</div>
{{ end }}`,
}

var templateNewsViewsMapChecksums = map[string]string{
	"home": "1d6af89dabab628b0e32ba7a4d2cabf6e4de3f562d7913bd736adf96b24c5e1b",
	"media": "3ade9e3ac069404656e1c8805fbc8dfb437cd32079f0b59bc9aed6cb29f3799e",
	"official": "bb42be6ef1e05030745feb46bc34235ed59bec785eb1b607802b74b2af479a14",
	"sources": "37878c8c51f13928e468057b7c88e8a0183a3c69094c8e702ac583d7edabd31d",
	"travel": "e292133fa926a1dd87c6f77922def87f3867e9e6bc182dbca2bbcba19a984f06",
	"visa": "43ef5b25e68074fb88516126610fc4a7ea51ecfaa55e01b78fa5435cdb942214",
}

// Code generated by go generate; DO NOT EDIT.
// 2018-07-28 11:35:22.282883127 +0300 MSK m=+0.010259390

package template

var templateNewsCommonMap = map[string]string{
	"news_disclaimer": `{{ define "disclaimer" }}

    <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="section__text">
                <h5>DISCLAIMER</h5>
                <div>
                    <p><b>We do not make any guarantees for data accuracy.</b></p>
                    <p>Visadb.io News is not responsible for, and expressly disclaims all liability for, damages of any kind arising out of use, reference to, or reliance on any information contained within the site. While the information contained within the site is periodically updated, no guarantee is given that the information provided in this Web site is correct, complete, and up-to-date.
                        Although the Visadb.io may include links providing direct access to other Internet resources, including Websites,visadb.io News is not responsible for the accuracy or content of information contained in these sites.responsible for the accuracy or content of information contained in these sites.
                    </p>
                </div>
                <div>
                    <p><b>Privacy policy:</b></p>
                    <p>By using visadb.io News website you agree to allow us to collect personal information about you and to track your behaviour by using cookies.</p>
                </div>
            </div>
    </div>

{{ end }}`,
	"news_layout": `{{ define "news_base" }}
<!doctype html>
<!--
  Material Design Lite
  Copyright 2015 Google Inc. All rights reserved.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License
-->
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="description" content="A front-end template that helps you build fast, modern mobile web apps.">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0">
    <title>VisaDB - NEWS</title>

    <!-- Add to homescreen for Chrome on Android -->
    <meta name="mobile-web-app-capable" content="yes">
    <link rel="icon" sizes="192x192" href="images/android-desktop.png">

    <!-- Add to homescreen for Safari on iOS -->
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <meta name="apple-mobile-web-app-title" content="Material Design Lite">
    <link rel="apple-touch-icon-precomposed" href="images/ios-desktop.png">

    <!-- Tile icon for Win8 (144x144 + tile color) -->
    <meta name="msapplication-TileImage" content="images/touch/ms-touch-icon-144x144-precomposed.png">
    <meta name="msapplication-TileColor" content="#3372DF">

    <link rel="shortcut icon" href="images/favicon.png">

    <!-- SEO: If your mobile URL is different from the desktop URL, add a canonical link to the desktop page https://developers.google.com/webmasters/smartphone-sites/feature-phones -->
    <!--
    <link rel="canonical" href="http://www.example.com/">
    -->

    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=en">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.deep_purple-pink.min.css">
    <link rel="stylesheet" type="text/css" href="{{ route "news_stylesheet" "name" "news"}}">
    <link rel="stylesheet" type="text/css" href="{{ route "news_stylesheet" "name" "getmdl-select.min"}}">
    <link rel="stylesheet" type="text/css" href="{{ route "news_stylesheet" "name" "dialog-polyfill"}}">
    {{ if .csrf }}
        <meta name="X-CSRF-Token" value="{{ .csrf }}">
    {{ end }}
    <style>
        #view-source {
            position: fixed;
            display: block;
            right: 0;
            bottom: 0;
            margin-right: 40px;
            margin-bottom: 40px;
            z-index: 900;
        }
    </style>
    <script src="https://app.brandquiz.io/embed"></script>
</head>
<body class="mdl-demo mdl-color--grey-100 mdl-color-text--grey-700 mdl-base">
<!-- DATA SOURCES DIALOG -->
<dialog class="mdl-dialog" id="sourcesContent">
    <h4 class="mdl-dialog__title">Data Sources</h4>
    <div>
        {{ template "sources" }}
    </div>
    <div class="mdl-dialog__actions">
        <button type="button" class="mdl-button close">Close</button>
    </div>
</dialog>
<!-- CONTACTS DIALOG -->
<dialog class="mdl-dialog" id="contactsContent">
    <div class="mdl-dialog__content">
        <div class="brandquiz_embed" data-embed="visadbio/news-contact"></div>
    </div>
    <div class="mdl-dialog__actions">
        <button type="button" class="mdl-button close">Close</button>
    </div>
</dialog>
<!-- SUBSCRIBE DIALOG -->
<dialog class="mdl-dialog" id="subscribeContent">
    <div class="mdl-dialog__content">
    <div class="brandquiz_embed1" data-embed="visadbio/visadb-news-signup"></div>
    </div>
    <div class="mdl-dialog__actions">
        <button type="button" class="mdl-button close">Close</button>
    </div>
</dialog>
<!-- ENTRIES FULL CONTENT DIALOG -->
<dialog class="mdl-dialog" id="entryContent">
    <div class="mdl-dialog__title" id="entryContentHeader"></div>
    <div class="mdl-dialog__content">
        <div id="entryContentBody"></div>
    </div>
    <div class="mdl-dialog__actions">
        <button type="button" class="mdl-button close">Close</button>
    </div>
</dialog>

    <div class="mdl-layout mdl-js-layout mdl-layout--fixed-header">
        <header class="mdl-layout__header mdl-layout__header--scroll mdl-color--primary">
            <div class="mdl-layout--large-screen-only mdl-layout__header-row demo-layout-transparent">
                <div class="mdl-layout__header-row">
                    <!-- Title -->
                    <span class="mdl-layout-title">VISADB.IO | NEWS <span class="red-span">LIVE</span></span>
                    <div class="mdl-layout-spacer"></div>
                    <!-- Navigation -->
                    <nav class="mdl-navigation">
                        <a role="button" class="mdl-navigation__link--transparent" href="http://visadb.io" target="_blank">SEARCH VISA</a>
                        <a role="button" class="mdl-navigation__link--transparent" href="#" id="show-sources">DATA SOURCES</a>
                        <a role="button" id="show-subscribe" class="mdl-navigation__link--transparent" href="#">SUBSCRIBE</a>
                        {{/*<form action="#" id="subscribe-header">*/}}
                            {{/*<div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">*/}}
                                {{/*<input class="mdl-textfield__input" type="text" id="subscribe" required>*/}}
                                {{/*<label class="mdl-textfield__label" for="subscribe">Your Email...</label>*/}}
                            {{/*</div>*/}}
                            {{/*<input type="submit" form="subscribe-header" class="mdl-button mdl-js-button mdl-js-ripple-effect" value="OK">*/}}
                        {{/*<button type="submit" form="subscribe-header" class="mdl-button mdl-js-button mdl-js-ripple-effect">OK</button>*/}}
                        {{/*</form>*/}}
                        <a role="button" id="show-contacts" class="mdl-navigation__link--transparent" href="#">CONTACTS</a>
                    </nav>

                </div>

            </div>

            <div class="mdl-layout--large-screen-only mdl-layout__header-row">
            </div>
            <div class="mdl-layout--large-screen-only mdl-layout__header-row">
            {{/*<h3>VISADB - News</h3>*/}}
                <div class="mdl-layout-spacer"></div>
                <div class="mdl-textfield mdl-js-textfield getmdl-select getmdl-select__fix-height worldwide">
                    <input type="text" value="" class="mdl-textfield__input" id="worldwide" readonly>
                    <input type="hidden" value="" name="worldwide">
                    <i class="mdl-icon-toggle__label material-icons">keyboard_arrow_down</i>
                    <label for="worldwide" class="mdl-textfield__label">Choose Country</label>
                    <ul for="worldwide" class="mdl-menu mdl-menu--bottom-left mdl-js-menu mdl-list">
                        <li class="mdl-menu__item" data-val="WW" data-selected="true">WORLDWIDE</li>
                        {{ range .countries}}
                            <li class="mdl-menu__item mdl-list__item" data-val="{{ .Code }}">
                                <span class="mdl-list__item-primary-content">
                                {{/*<i class="material-icons mdl-list__item-icon">tab</i>*/}}
                                {{ .Name }}
                                </span>
                            </li>
                        {{ end }}
                    </ul>
                </div>
                <div class="mdl-layout-spacer"></div>
            </div>

            <div class="mdl-layout--large-screen-only mdl-layout__header-row">
            </div>
            <div class="mdl-layout__tab-bar mdl-js-ripple-effect mdl-color--primary-dark">
                <div class="mdl-layout-spacer"></div>
                <a href="#official" class="mdl-layout__tab is-active">Official News</a>
                <a href="#media" class="mdl-layout__tab">Media News</a>
                <a href="#travel" class="mdl-layout__tab">Travel Alerts</a>
                <a href="#security" class="mdl-layout__tab">Security Alerts</a>
                <a href="#visa" class="mdl-layout__tab">Visa Changes</a>
                <div class="mdl-layout-spacer"></div>
            </div>
        </header>
        <main class="mdl-layout__content">
        {{template "content" . }}
        {{ template "ticker" .ticker }}
        </main>

    <script src="https://code.getmdl.io/1.3.0/material.min.js"></script>
    <script type="text/javascript" src="{{ route "javascript" }}"></script>
    <script type="text/javascript" src="{{ route "news_mdlselect" }}"></script>
    <script type="text/javascript" src="{{ route "news_dialog-polyfill" }}"></script>
    <script type="text/javascript" src="{{ route "news_news" }}"></script>
        <script>
            var dialog = document.querySelector('#contactsContent');
            var showDialogButton = document.querySelector('#show-contacts');
            if (! dialog.showModal) {
                dialogPolyfill.registerDialog(dialog);
            }
            showDialogButton.addEventListener('click', function() {
                dialog.showModal();
            });
            dialog.querySelector('.close').addEventListener('click', function() {
                dialog.close();
            });

            var subscribe = document.querySelector('#subscribeContent');
            var showSubscribeButton = document.querySelector('#show-subscribe');
            if (! subscribe.showModal) {
                dialogPolyfill.registerDialog(subscribe);
            }
            showSubscribeButton.addEventListener('click', function() {
                subscribe.showModal();
            });
            subscribe.querySelector('.close').addEventListener('click', function() {
                subscribe.close();
            });

            var sources = document.querySelector('#sourcesContent');
            var showSourcesButton = document.querySelector('#show-sources');
            if (! sources.showModal) {
                dialogPolyfill.registerDialog(sources);
            }
            showSourcesButton.addEventListener('click', function() {
                sources.showModal();
            });
            sources.querySelector('.close').addEventListener('click', function() {
                sources.close();
            });

            var entryDialog = document.querySelector("#entryContent");
            if (! entryDialog.showModal) {
                dialogPolyfill.registerDialog(entryDialog);
            }
            entryDialog.querySelector('.close').addEventListener('click', function() {
                let body = entryDialog.querySelector('#entryContentBody');
                body.innerHTML = "";
                entryDialog.close();
            });
        </script>
        <script>
            "use strict";var brandquizEmbed1=brandquizEmbed1||{};brandquizEmbed1.baseurl="https://app.brandquiz.io",brandquizEmbed1.init=function(){function e(){jQuery(".brandquiz_embed1").attr("data-baseurl")&&(brandquizEmbed1.baseurl=jQuery(".brandquiz_embed1").attr("data-baseurl")),brandquizEmbed1.include(brandquizEmbed1.baseurl+"/js/jquery.responsiveiframe.js",function(){jQuery(document).ready(function(){brandquizEmbed1.embed()})})}"undefined"==typeof jQuery?brandquizEmbed1.include("https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js",function(){jQuery(document).ready(function(){e()})}):e()},brandquizEmbed1.embed=function(){document.getElementsByTagName("meta").viewport?document.getElementsByTagName("meta").viewport.content='width=device-width, initial-scale=1, maximum-scale=1"':jQuery("head").append('<meta name="viewport" content="width=device-width, initial-scale=1"/>');var e=jQuery(".brandquiz_embed1").attr("data-embed");jQuery(".brandquiz_embed1").css("max-width","100%"),jQuery(".brandquiz_embed1").attr("data-width")&&jQuery(".brandquiz_embed1").css("max-width",jQuery(".brandquiz_embed1").attr("data-width")+"px");var a=jQuery('<iframe frameborder="0" style="width:100%; overflow:hidden;" src="'+brandquizEmbed1.baseurl+"/"+e+"?embed=1&src="+encodeURIComponent(document.location.href)+'">');jQuery(".brandquiz_embed1").html(a),jQuery(".brandquiz_embed1").attr("data-min-height")&&jQuery(".brandquiz_embed1 iframe").css("min-height",jQuery(".brandquiz_embed1").attr("data-min-height")+"px"),jQuery(".brandquiz_embed1 iframe").responsiveIframe({xdomain:"*",scrollToTop:!1})},brandquizEmbed1.include=function(e,a){var d=document.getElementsByTagName("head")[0],t=document.createElement("script");t.src=e,t.type="text/javascript",t.onload=t.onreadystatechange=function(){t.readyState?"complete"!==t.readyState&&"loaded"!==t.readyState||(t.onreadystatechange=null,a()):a()},d.appendChild(t)},document.addEventListener("DOMContentLoaded",function(){brandquizEmbed1.init()}),window.addEventListener("message",function(e){if("brandquiz-cookie-failed"==e.data&&navigator.cookieEnabled&&!localStorage.getItem("reload")){localStorage.setItem("reload","true");var a=brandquizEmbed1.baseurl+"/redirect-back/"+btoa("rubsalt"+encodeURIComponent(window.location.href));window.location.replace(a)}});
        </script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/d3/3.5.3/d3.min.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/topojson/1.6.9/topojson.min.js"></script>
        <script src="//datamaps.github.io/scripts/0.4.4/datamaps.world.min.js"></script>
        <script defer>
            var map = new Datamap({
                element: document.getElementById('map'),
                responsive: true,
                geographyConfig: {
                    highlightOnHover: true,
                    popupOnHover: true,
                    popupTemplate: function(geography, data) {
                        return '<div class="hoverinfo" style="white-space: pre-wrap;">' + geography.properties.name + ':\n' +  data.risk + ' '
                    },
                    fillOpacity: 0.75,
                    animate: true,
                    highlightFillColor: '#FA0FA0',
                    highlightBorderColor: 'rgba(250, 15, 160, 0.2)',
                    highlightBorderWidth: 2,
                    highlightFillOpacity: 0.85,
                    responsive: true
                },
                projection: 'mercator',
                fills: {
                    defaultFill: "#f7f6f5",
                    NormalPrecautions: "#38b801",
                    NormalPrecautionsReg: "#38b801",
                    HighCaution: "#f3f30e",
                    HighCautionReg: "#f3f30e",
                    AvoidNonEssential: "#ff8b00",
                    AvoidNonEssentialReg: "#ff8b00 ",
                    AvoidAllReg: "#ff0000",
                    AvoidAll: "#ff0000"
                },
                data:  (function(){
                    let xhr = new XMLHttpRequest();
                    var mapdata;
                    xhr.onreadystatechange = function() {
                        if (xhr.readyState == 4 && xhr.status == 200) {
                            mapdata = JSON.parse(xhr.responseText);
                        }
                    };
                    xhr.open("GET", "/news/security", false);
                    try {
                        xhr.send();
                    } catch (err) {
                        console.log(err);
                    }

                    return mapdata;
                })()
            });

            // map.legend();
            var colors = d3.scale.category10();
            window.addEventListener('resize', function() {
                map.resize();
            });
        </script>
</body>
</html>
{{end}}`,
	"news_sources": `{{ define "sources"}}

    <br>
        <div class="mdl-cell mdl-cell--12-col">
                    <h5>Official News:</h5>
                    <div>
                        <p><b>Source:</b> Global Ministry of foreign Affairs - Global Immigration Department - Global Tourism Boards</p>
                    </div>
        </div>
        <div class="mdl-cell mdl-cell--12-col">
                    <h5> Media News:</h5>
                    <div>
                        <p><b>Source:</b> Mainstream media news outlet i.e. CNN, BBC and others and Media Migration Policy.</p>
                    </div>
        </div>
        <div class="mdl-cell mdl-cell--12-col">
                    <h5>Travel Alerts:</h5>
                    <div>
                        <p><b>Source:</b> CDC Yellow Fever and CDC Travel Alerts</p>
                    </div>
        </div>
        <div class="mdl-cell mdl-cell--12-col">
                    <h5>Security Alerts:</h5>
                    <div>
                        <p><b>Source:</b> </p>
                    </div>
        </div>

{{ template "disclaimer" }}

{{ end }}`,
	"news_ticker": `{{ define "ticker"}}
        <br>
        <div class="ticker">
            <marquee behavior="scroll" direction="left">{{ . }}</marquee>
        </div>
{{ end }}`,
	"news_travel": `{{ define "travel" }}

<div class="mdl-layout__tab-panel mdl-cell mdl-cell--9-col-desktop mdl-cell--7-col-tablet mdl-cell--9-col-phone" id="travel">
    <br>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--9-col-desktop mdl-cell--6-col-tablet mdl-cell--4-col-phone">
            <div class="mdl-card__supporting-text">
                <h4>Travel Alerts</h4>
                Dolore ex deserunt aute fugiat aute nulla ea sunt aliqua nisi cupidatat eu. Nostrud in laboris labore nisi amet do dolor eu fugiat consectetur elit cillum esse.
            </div>
        </div>
    {{template "countries" .countries}}
    </section>
    <section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
        <div class="mdl-card mdl-cell mdl-cell--12-col">
            <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
                <h4 class="mdl-cell mdl-cell--12-col">News Feed</h4>
            {{ range $.entries }}
                <div class="section__feed-container mdl-cell mdl-cell--2-col mdl-cell--1-col-phone">
                    <div class="section__feed-container__feed">{{ .Feed.Title }}</div>
                </div>
                <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
                    <h5>{{ .Title }}</h5>
                    <a href="{{ route "feedEntry" "feedID" .Feed.ID "entryID" .ID }}">{{ .Title }}</a>
                </div>
            {{ end }}
            </div>
        </div>
    </section>
</div>

{{end }}`,
}

var templateNewsCommonMapChecksums = map[string]string{
	"disclaimer": "fce75c71a10a694a25c54d952094b10870353309329e1e4e8bf9370ac9caf712",
	"layout":     "b0432d934b5b11064ff3027f055a10014e5da9ffb4559ec15e5c6145f1de8219",
	"sources":    "ed72d5167279a038b60d3d2cad8347b7404b354e85c753167a298acfb7d344e3",
	"ticker":     "94cbf10646256355fa43d66badc31068199d2775f67266674ccb3ec55d75b307",
	"travel":     "55313525e6c21868800810ade5a5cf82c98d5bdc6479624c71990eddbeb3aa97",
}

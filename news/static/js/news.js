/*jshint esversion: 6 */
(function() {
    'use strict';

    // class DialogHandler {
    //     constructor() {
    //         this.dialog = document.querySelector('dialog');
    //         this.showDialogButton = document.querySelector('#show-dialog');
    //     }
    //     if (dialog.showModal) {
    //         dialogPolyfill.registerDialog(this.dialog);
    //     }
    //     this.showDialogButton.addEventListener('click', function() {
    //         dialog.showModal();
    //     });
    //     this.dialog.querySelector('.close').addEventListener('click', function() {
    //         dialog.close();
    //     });
    // }

    class DomHelper {
        static isVisible(element) {
            return element.offsetParent !== null;
        }

        static openNewTab(url) {
            let win = window.open(url, "_blank");
            win.focus();
        }

        static scrollPageTo(element) {
            let windowScrollPosition = window.pageYOffset;
            let windowHeight = document.documentElement.clientHeight;
            let viewportPosition = windowScrollPosition + windowHeight;
            let itemBottomPosition = element.offsetTop + element.offsetHeight;

            if (viewportPosition - itemBottomPosition < 0 || viewportPosition - element.offsetTop > windowHeight) {
                window.scrollTo(0, element.offsetTop - 10);
            }
        }

        static getVisibleElements(selector) {
            let elements = document.querySelectorAll(selector);
            let result = [];

            for (let i = 0; i < elements.length; i++) {
                if (this.isVisible(elements[i])) {
                    result.push(elements[i]);
                }
            }

            return result;
        }
    }

    class TouchHandler {
        constructor() {
            this.reset();
        }

        reset() {
            this.touch = {
                start: {x: -1, y: -1},
                move: {x: -1, y: -1},
                element: null
            };
        }

        calculateDistance() {
            if (this.touch.start.x >= -1 && this.touch.move.x >= -1) {
                let horizontalDistance = Math.abs(this.touch.move.x - this.touch.start.x);
                let verticalDistance = Math.abs(this.touch.move.y - this.touch.start.y);

                if (horizontalDistance > 30 && verticalDistance < 70) {
                    return this.touch.move.x - this.touch.start.x;
                }
            }

            return 0;
        }

        findElement(element) {
            if (element.classList.contains("touch-item")) {
                return element;
            }

            return DomHelper.findParent(element, "touch-item");
        }

        onTouchStart(event) {
            if (event.touches === undefined || event.touches.length !== 1) {
                return;
            }

            this.reset();
            this.touch.start.x = event.touches[0].clientX;
            this.touch.start.y = event.touches[0].clientY;
            this.touch.element = this.findElement(event.touches[0].target);
        }

        onTouchMove(event) {
            if (event.touches === undefined || event.touches.length !== 1 || this.element === null) {
                return;
            }

            this.touch.move.x = event.touches[0].clientX;
            this.touch.move.y = event.touches[0].clientY;

            let distance = this.calculateDistance();
            let absDistance = Math.abs(distance);

            if (absDistance > 0) {
                let opacity = 1 - (absDistance > 75 ? 0.9 : absDistance / 75 * 0.9);
                let tx = distance > 75 ? 75 : (distance < -75 ? -75 : distance);

                this.touch.element.style.opacity = opacity;
                this.touch.element.style.transform = "translateX(" + tx + "px)";
            }
        }

        onTouchEnd(event) {
            if (event.touches === undefined) {
                return;
            }

            if (this.touch.element !== null) {
                let distance = Math.abs(this.calculateDistance());

                if (distance > 75) {
                    EntryHandler.toggleEntryStatus(this.touch.element);
                    this.touch.element.style.opacity = 1;
                    this.touch.element.style.transform = "none";
                }
            }

            this.reset();
        }

        listen() {
            let elements = document.querySelectorAll(".touch-item");

            elements.forEach((element) => {
                element.addEventListener("touchstart", (e) => this.onTouchStart(e), false);
                element.addEventListener("touchmove", (e) => this.onTouchMove(e), false);
                element.addEventListener("touchend", (e) => this.onTouchEnd(e), false);
                element.addEventListener("touchcancel", () => this.reset(), false);
            });
        }
    }

    class FormHandler {
        static handleSubmitButtons() {
            let elements = document.querySelectorAll("form");
            elements.forEach((element) => {
                element.onsubmit = () => {
                    let button = document.querySelector("button");

                    if (button) {
                        button.innerHTML = button.dataset.labelLoading;
                        button.disabled = true;
                    }
                };
            });
        }
    }

    class MouseHandler {
        onClick(selector, callback) {
            let elements = document.querySelectorAll(selector);
            elements.forEach((element) => {
                element.onclick = (event) => {
                    event.preventDefault();
                    callback(event);
                };
            });
        }
    }

    class RequestBuilder {
        constructor(url) {
            this.callback = null;
            this.url = url;
            this.options = {
                method: "POST",
                cache: "no-cache",
                credentials: "include",
                body: null,
                headers: new Headers({
                    "Content-Type": "application/json",
                    "X-Csrf-Token": this.getCsrfToken()
                })
            };
        }

        withBody(body) {
            this.options.body = JSON.stringify(body);
            return this;
        }

        withCallback(callback) {
            this.callback = callback;
            return this;
        }

        getCsrfToken() {
            let element = document.querySelector("meta[name=X-CSRF-Token]");
            if (element !== null) {
                return element.getAttribute("value");
            }

            return "";
        }

        execute() {
            fetch(new Request(this.url, this.options)).then((response) => {
                if (this.callback) {
                    this.callback(response);
                }
            });
        }
    }

    class ModalHandler {
        static exists() {
            return document.getElementById("modal-container") !== null;
        }

        static open(fragment) {
            if (ModalHandler.exists()) {
                return;
            }

            let container = document.createElement("div");
            container.id = "modal-container";
            container.appendChild(document.importNode(fragment, true));
            document.body.appendChild(container);
            container.style.visibility = 'visible';

            let closeButton = document.querySelector("a.btn-close-modal");
            if (closeButton !== null) {
                closeButton.onclick = (event) => {
                    event.preventDefault();
                    ModalHandler.close();
                };
            }
        }

        static close() {
            let container = document.getElementById("modal-container");
            if (container !== null) {
                container.parentNode.removeChild(container);
            }
        }
    }


    class TabHandler {

        constructor(tabname, offset, country="WORLDWIDE") {
            this.tabname = tabname;
            this.offset = offset;
            this.limit = 10;
            this.country = "";
            this.PrevButton = null;
            this.NextButton = null;
        }

        async LoadTab(){
            let tab = document.getElementById(this.tabname);
            while(!document.querySelector("li.selected")) {
                await new Promise(r => setTimeout(r, 300));
            }
            this.country = document.querySelector("li.selected").textContent.trim();
            let xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function() {
                if (xhr.readyState == 4 && xhr.status == 200) {
                    tab.innerHTML = xhr.responseText;
                    let prev = tab.querySelector("[data-page='prev']");
                    let next = tab.querySelector("[data-page='next']");
                    this.PrevButton = prev;
                    this.NextButton = next;
                    this.ListenButtons();
                    this.ListenDownloads();
                }
            }.bind( this );
            xhr.open("GET", "/news/" + this.tabname + "?offset=" + this.offset + "&limit=" + this.limit + "&country=" + this.country, true);
            try {
                xhr.send();
            } catch (err) {
                tab.innerHTML = "Not Found"
            }
        }

        ListenButtons(){
            this.PrevButton.addEventListener("click",(event) => {
                event.stopPropagation();
                this.offset = this.offset - this.limit;
                this.LoadTab();
            },false);

            this.NextButton.addEventListener("click",(event) => {
                event.stopPropagation();
                this.offset = this.offset + this.limit;
                this.LoadTab();
            },false);
        }

        ListenDownloads(){
            let downloads = document.querySelectorAll("a[data-fetch-content-url]");
            downloads.forEach((element) => {
                element.addEventListener("click",(event) => {
                    event.stopPropagation();
                    let header = event.target.parentElement.cloneNode(true);
                    let link = header.querySelector("a[data-fetch-content-url]");
                    header.removeChild(link);
                    this.fetchOriginalContent(event.target, header);
                });
            });
        }

        fetchOriginalContent(element, header){
            document.querySelector("#entryContentHeader").innerHTML = header.innerHTML;
            let request = new RequestBuilder(element.dataset.fetchContentUrl);
            request.withCallback((response) => {
                response.json().then((data) => {
                    if (data.hasOwnProperty("content")) {
                        // let template = document.createTextNode(data.content);
                        let body = document.querySelector("#entryContentBody");
                        body.innerHTML = data.content;
                        let images = body.querySelectorAll("audio, canvas, iframe, img, svg, video");
                        images.forEach((img) => {
                            img.style.maxWidth = "95%";
                        });
                        document.querySelector('#entryContent').showModal();
                    } else {
                        document.querySelector("#entryContentBody").innerHTML = '<div>NO CONTENT</div>';
                        document.querySelector('#entryContent').showModal();
                    }
                });
            });
            request.execute();
        }
    }

    document.addEventListener("DOMContentLoaded", function() {
        FormHandler.handleSubmitButtons();

        let mediaHandler = new TabHandler("media", 0);
        mediaHandler.LoadTab();
        let officialHandler = new TabHandler("official", 0);
        officialHandler.LoadTab();
        let visaHandler = new TabHandler("visa", 0);
        visaHandler.LoadTab();
        let travelHandler = new TabHandler("travel", 0);
        travelHandler.LoadTab();

        let countries = document.querySelector("ul[for=worldwide]");
        countries.addEventListener("click", function(e){
            let country = e.target.textContent.trim();
            mediaHandler = new TabHandler("media", 0, country);
            mediaHandler.LoadTab();
            officialHandler = new TabHandler("official", 0, country);
            officialHandler.LoadTab();
            visaHandler = new TabHandler("visa", 0, country);
            visaHandler.LoadTab();
            travelHandler = new TabHandler("travel", 0, country);
            travelHandler.LoadTab();
        })

    });

})();

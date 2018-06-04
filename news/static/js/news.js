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

    document.addEventListener("DOMContentLoaded", function() {
        FormHandler.handleSubmitButtons();

        let touchHandler = new TouchHandler();
        touchHandler.listen();

        let mouseHandler = new MouseHandler();

        mouseHandler.onClick("a[data-save-entry]", (event) => {
            event.preventDefault();
            EntryHandler.saveEntry(event.target);
        });

        mouseHandler.onClick("a[data-toggle-bookmark]", (event) => {
            event.preventDefault();
            EntryHandler.toggleBookmark(event.target);
        });

        mouseHandler.onClick("a[data-toggle-status]", (event) => {
            event.preventDefault();

            let currentItem = DomHelper.findParent(event.target, "item");
            if (currentItem) {
                EntryHandler.toggleEntryStatus(currentItem);
            }
        });

        mouseHandler.onClick("a[data-fetch-content-entry]", (event) => {
            event.preventDefault();
            EntryHandler.fetchOriginalContent(event.target);
        });

        mouseHandler.onClick("button[data-on-click=addToFilter]", (event) => {
            event.preventDefault();
            let inputs = document.querySelectorAll("#keyword");
            if (inputs && inputs.length > 0) {
                let newWord = inputs[0].value;
                FiltersHandler.addToFilter(newWord);
                inputs[0].value = "";
                inputs[0].focus();
            }
        });

        mouseHandler.onClick("div[data-on-click=rmFromFilter]", (event) => {
            event.stopPropagation();
            FiltersHandler.rmFromFilter(event.target.parentNode);
        });

    });

})();

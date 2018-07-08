// Code generated by go generate; DO NOT EDIT.
// 2018-07-08 17:41:19.706350034 +0300 MSK m=+0.010471255

package static

var NewsStylesheets = map[string]string{
	"dialog-polyfill":   `dialog{position:absolute;left:0;right:0;width:-moz-fit-content;width:-webkit-fit-content;width:fit-content;height:-moz-fit-content;height:-webkit-fit-content;height:fit-content;margin:auto;border:solid;padding:1em;background:#fff;color:#000;display:block}dialog:not([open]){display:none}dialog+.backdrop{position:fixed;top:0;right:0;bottom:0;left:0;background:rgba(0,0,0,.1)}._dialog_overlay{position:fixed;top:0;right:0;bottom:0;left:0}dialog.fixed{position:fixed;top:50%;transform:translate(0,-50%)}`,
	"getmdl-select.min": `.getmdl-select{outline:0}.getmdl-select .mdl-textfield__input{cursor:pointer}.getmdl-select .selected{background-color:#ddd}.getmdl-select .mdl-icon-toggle__label{float:right;margin-top:-30px;color:rgba(0,0,0,.4);transform:rotate(0);transition:transform .3s}.getmdl-select.is-focused .mdl-icon-toggle__label{color:#3f51b5;transform:rotate(180deg)}.getmdl-select .mdl-menu__container{width:100%!important;margin-top:2px}.getmdl-select .mdl-menu__container .mdl-menu{width:100%}.getmdl-select .mdl-menu__container .mdl-menu .mdl-menu__item{font-size:16px}.getmdl-select__fix-height .mdl-menu__container .mdl-menu{overflow-y:auto;max-height:288px!important}.getmdl-select__fix-height .mdl-menu.mdl-menu--top-left{bottom:auto;top:0}`,
	"news":              `html,body{font-family:roboto,helvetica,sans-serif;margin:0;padding:0}.mdl-dialog{width:70%}.ticker{background-color:#673ab7;color:#fff;font-size:18px;height:50px;margin-top:20px}.demo-layout-transparent,.mdl-layout__header--transparent{background-color:#fff;color:#000}.mdl-navigation__link--transparent{color:#000;display:block;line-height:64px;padding:0 24px;text-decoration:none}span.red-span{color:red}.worldwide{color:#fff}.side-cell{padding-left:40px}.boxes-center{max-width:860px;margin-top:20px;margin-bottom:60px}.mdl-menu{height:200px;overflow-y:scroll}.demo-card-square{margin-top:50px}.mdl-demo .mdl-layout__header-row{padding-left:40px}.mdl-demo .mdl-layout.is-small-screen .mdl-layout__header-row h3{font-size:inherit}.mdl-demo .mdl-layout__tab-bar-button{display:none}.mdl-demo .mdl-layout.is-small-screen .mdl-layout__tab-bar .mdl-button{display:none}.mdl-demo .mdl-layout:not(.is-small-screen) .mdl-layout__tab-bar,.mdl-demo .mdl-layout:not(.is-small-screen) .mdl-layout__tab-bar-container{overflow:visible}.mdl-demo .mdl-layout__tab-bar-container{height:64px}.mdl-demo .mdl-layout__tab-bar{padding:0;padding-left:16px;box-sizing:border-box;height:100%;width:100%}.mdl-demo .mdl-layout__tab-bar .mdl-layout__tab{height:64px;line-height:64px}.mdl-demo .mdl-layout__tab-bar .mdl-layout__tab.is-active::after{background-color:#fff;height:4px}.mdl-demo main>.mdl-layout__tab-panel{padding:8px;padding-top:48px}.mdl-demo .mdl-card{height:auto;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex-direction:column;-ms-flex-direction:column;flex-direction:column}.mdl-demo .mdl-card>*{height:auto}.mdl-demo .mdl-card .mdl-card__supporting-text{margin:40px;-webkit-flex-grow:1;-ms-flex-positive:1;flex-grow:1;padding:0;color:inherit;width:calc(100% - 80px)}.mdl-demo.mdl-demo .mdl-card__supporting-text h4{margin-top:0;margin-bottom:20px}.mdl-demo .mdl-card__actions{margin:0;padding:4px 40px;color:inherit}.mdl-demo .mdl-card__actions a{color:#00bcd4;margin:0}.mdl-demo .mdl-card__actions a:hover,.mdl-demo .mdl-card__actions a:active{color:inherit;background-color:transparent}.mdl-demo .mdl-card__supporting-text+.mdl-card__actions{border-top:1px solid rgba(0,0,0,.12)}.mdl-demo #add{position:absolute;right:40px;top:36px;z-index:999}.mdl-demo .mdl-layout__content section:not(:last-of-type){position:relative;margin-bottom:48px}.mdl-demo section.section--center{max-width:860px}.mdl-demo #features section.section--center{max-width:620px}.mdl-demo section>header{display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-align-items:center;-ms-flex-align:center;align-items:center;-webkit-justify-content:center;-ms-flex-pack:center;justify-content:center}.mdl-demo section>.section__play-btn{min-height:200px}.mdl-demo section>header>.material-icons{font-size:3rem}.mdl-demo section>button{position:absolute;z-index:99;top:8px;right:8px}.mdl-demo section .section__circle{display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-align-items:center;-ms-flex-align:center;align-items:center;-webkit-justify-content:flex-start;-ms-flex-pack:start;justify-content:flex-start;-webkit-flex-grow:0;-ms-flex-positive:0;flex-grow:0;-webkit-flex-shrink:1;-ms-flex-negative:1;flex-shrink:1}.mdl-demo section .section__text{-webkit-flex-grow:1;-ms-flex-positive:1;flex-grow:1;-webkit-flex-shrink:0;-ms-flex-negative:0;flex-shrink:0;padding-top:8px}.mdl-demo section .section__text h5{font-size:inherit;margin:0;margin-bottom:.5em}.mdl-demo section .section__text h6{font-size:inherit;margin:0;margin-bottom:.5em}.mdl-demo section .section__text a{text-decoration:none}.mdl-demo section .section__circle-container>.section__circle-container__circle{width:64px;height:64px;margin:8px 0}.mdl-demo section .section__feed-container>.section__feed-container__feed{width:84px;height:84px;margin:8px 5px}.mdl-demo section.section--footer .section__circle--big{width:100px;height:100px;border-radius:50px;margin:8px 32px}.mdl-demo .is-small-screen section.section--footer .section__circle--big{width:50px;height:50px;border-radius:25px;margin:8px 16px}.mdl-demo section.section--footer{padding:64px 0;margin:0 -8px -8px}.mdl-demo section.section--center .section__text:not(:last-child){border-bottom:1px solid rgba(0,0,0,.13)}.mdl-demo .mdl-card .mdl-card__supporting-text>h3:first-child{margin-bottom:24px}.mdl-demo .mdl-layout__tab-panel:not(#overview){background-color:#fff}.mdl-demo #features section{margin-bottom:72px}.mdl-demo #features h4,#features h5{margin-bottom:16px}.mdl-demo .toc{border-left:4px solid #c1eef4;margin:24px;padding:0;padding-left:8px;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex-direction:column;-ms-flex-direction:column;flex-direction:column}.mdl-demo .toc h4{font-size:.9rem;margin-top:0}.mdl-demo .toc a{color:#4dd0e1;text-decoration:none;font-size:16px;line-height:28px;display:block}.mdl-demo .mdl-menu__container{z-index:99}`,
}

var NewsStylesheetsChecksums = map[string]string{
	"dialog-polyfill":   "8c06ee177c2c8700b9e92049132da55d6f48219264e44b597e28728f298d4fc1",
	"getmdl-select.min": "f4dd68045fb54a443d56d255fc825e42a7e7450de7abec4ce78d4723f27f8a38",
	"news":              "e9708479b1d9b87bd1c54897cc6236ba6e15463cd2992b12b8a7cbe7531f4c1a",
}

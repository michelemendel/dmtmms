
let isLoaded = false;

function loadLibs(callsite) {
    if (isLoaded) {
        console.log("Libraries already loaded");
        return;
    }
    console.log("Loading libraries from ", callsite);
    // isLoaded = true;

    let url = "/public/tw-elements.umd.min.js";
    let body = document.body;
    let script = document.createElement("script");
    script.type = "text/javascript";
    script.src = url;
    body.appendChild(script);
}



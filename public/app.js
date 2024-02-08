
let isLoaded = false;

function loadLibs(callsite) {
    urls = [
        // TW Elements is free under AGPL, with commercial license required for specific uses. See more details: https://tw-elements.com/license/ and contact us for queries at tailwind@mdbootstrap.com
        "/public/tw-elements.umd.min.js"
    ]
    var body = document.body;
    var script = document.createElement('script');
    script.type = 'text/javascript';
    urls.forEach((url) => {
        script.src = url;
        body.appendChild(script);
    });
}

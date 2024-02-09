
// Load libraries that has to be loaded after the main app and HTMX partials are loaded.
// Since a script can't be loaded twice, we have to remove the current script and readd it.
// Match the ids and urls arrays.

ids = [
    'twe',
];

urls = [
    // TW Elements is free under AGPL, with commercial license required for specific uses. See more details: https://tw-elements.com/license/ and contact us for queries at tailwind@mdbootstrap.com
    "/node_modules/tw-elements/dist/js/tw-elements.umd.min.js",
]

function loadLibs() {
    var body = document.body;
    var script = document.createElement('script');
    script.type = 'text/javascript';

    for (var i = 0; i < ids.length; i++) {
    };

    for (var i = 0; i < ids.length; i++) {
        // Unload libs
        var els = document.querySelectorAll(`#${ids[i]}`);
        els.forEach(function (el) {
            // console.log(`Unloading ${ids[i]}`);
            el.remove();
        });

        // Load libs
        // console.log(`Loading ${ids[i]}`);
        script.id = ids[i];
        script.src = urls[i];
        body.appendChild(script);
    };
}


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
    let body = document.body;
    let script = document.createElement('script');
    script.type = 'text/javascript';

    for (var i = 0; i < ids.length; i++) {
        // Unload libs
        let els = document.querySelectorAll(`#${ids[i]}`);
        els.forEach(function (el) {
            el.remove();
        });

        // Load libs        
        script.id = ids[i];
        script.src = urls[i];
        body.appendChild(script);
    };
}

//--------------------------------------------------------------------------------
// This is to handle the issue with Flash Of Unstyled Content

// This is called only when the page has reloaded.
// We start by settings the body to be hidden, see appRoot.templ.
document.onreadystatechange = () => {
    if (document.readyState === "complete") {
        showBody();
    }
};

function hideBody() {
    document.body.style.visibility = 'hidden';
    document.body.style.opacity = 0;
}

function showBody() {
    setTimeout(() => {
        document.body.style.visibility = 'visible';
        document.body.style.opacity = 1;
    }, 100);
}



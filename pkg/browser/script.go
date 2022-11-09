package browser

const injectionScript = `
(function () {
    window.history.pushState = function (a, b, c) {
        window.collectURL(c);
    }
    window.history.replaceState = function (a, b, c) {
        window.collectURL(c);
    }
    Object.defineProperty(window.history, 'pushState', { 'writable': false, 'configurable': false });
    Object.defineProperty(window.history, 'replaceState', { 'writable': false, 'configurable': false });
    window.addEventListener('hashchange', function () {
        window.collectURL(document.location.href);
    });
})();
`

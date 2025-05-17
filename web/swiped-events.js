/*!
 * swiped-events.js - v@version@
 * Modified to run only on /learn route
 * https://github.com/john-doherty/swiped-events
 */
(function (window, document) {
    'use strict';

    // patch CustomEvent to allow constructor creation (IE/Chrome)
    if (typeof window.CustomEvent !== 'function') {
        window.CustomEvent = function (event, params) {
            params = params || { bubbles: false, cancelable: false, detail: undefined };
            var evt = document.createEvent('CustomEvent');
            evt.initCustomEvent(event, params.bubbles, params.cancelable, params.detail);
            return evt;
        };
        window.CustomEvent.prototype = window.Event.prototype;
    }

    let listenersAttached = false;

    function enableSwipeEvents() {
        if (listenersAttached) return;
        document.addEventListener('touchstart', handleTouchStart, false);
        document.addEventListener('touchmove', handleTouchMove, { passive: false });
        document.addEventListener('touchend', handleTouchEnd, false);
        listenersAttached = true;
    }

    function disableSwipeEvents() {
        if (!listenersAttached) return;
        document.removeEventListener('touchstart', handleTouchStart, false);
        document.removeEventListener('touchmove', handleTouchMove, { passive: false });
        document.removeEventListener('touchend', handleTouchEnd, false);
        listenersAttached = false;
    }

    function checkAndUpdateSwipeListeners() {
        if (window.location.pathname.startsWith("/learn")) {
            enableSwipeEvents();
        } else {
            disableSwipeEvents();
        }
    }

    // Observe SPA route changes
    window.addEventListener('popstate', checkAndUpdateSwipeListeners);
    window.addEventListener('pushstate', checkAndUpdateSwipeListeners);
    window.addEventListener('replacestate', checkAndUpdateSwipeListeners);

    // Patch history to detect SPA nav
    (function(history) {
        const pushState = history.pushState;
        const replaceState = history.replaceState;

        history.pushState = function () {
            const result = pushState.apply(history, arguments);
            window.dispatchEvent(new Event('pushstate'));
            return result;
        };
        history.replaceState = function () {
            const result = replaceState.apply(history, arguments);
            window.dispatchEvent(new Event('replacestate'));
            return result;
        };
    })(window.history);

    // Run once on load
    checkAndUpdateSwipeListeners();

    var xDown = null;
    var yDown = null;
    var xDiff = null;
    var yDiff = null;
    var timeDown = null;
    var startEl = null;

    function handleTouchEnd(e) {
        if (startEl !== e.target) return;

        var swipeThreshold = parseInt(getNearestAttribute(startEl, 'data-swipe-threshold', '20'), 10);
        var swipeUnit = getNearestAttribute(startEl, 'data-swipe-unit', 'px');
        var swipeTimeout = parseInt(getNearestAttribute(startEl, 'data-swipe-timeout', '500'), 10);
        var timeDiff = Date.now() - timeDown;
        var eventType = '';
        var changedTouches = e.changedTouches || e.touches || [];

        if (swipeUnit === 'vh') {
            swipeThreshold = Math.round((swipeThreshold / 100) * document.documentElement.clientHeight);
        }
        if (swipeUnit === 'vw') {
            swipeThreshold = Math.round((swipeThreshold / 100) * document.documentElement.clientWidth);
        }

        if (Math.abs(xDiff) > Math.abs(yDiff)) {
            if (Math.abs(xDiff) > swipeThreshold && timeDiff < swipeTimeout) {
                eventType = xDiff > 0 ? 'swiped-left' : 'swiped-right';
            }
        } else if (Math.abs(yDiff) > swipeThreshold && timeDiff < swipeTimeout) {
            eventType = yDiff > 0 ? 'swiped-up' : 'swiped-down';
        }

        if (eventType !== '') {
            var eventData = {
                dir: eventType.replace(/swiped-/, ''),
                touchType: (changedTouches[0] || {}).touchType || 'direct',
                xStart: parseInt(xDown, 10),
                xEnd: parseInt((changedTouches[0] || {}).clientX || -1, 10),
                yStart: parseInt(yDown, 10),
                yEnd: parseInt((changedTouches[0] || {}).clientY || -1, 10)
            };

            startEl.dispatchEvent(new CustomEvent('swiped', { bubbles: true, cancelable: true, detail: eventData }));
            startEl.dispatchEvent(new CustomEvent(eventType, { bubbles: true, cancelable: true, detail: eventData }));
        }

        xDown = null;
        yDown = null;
        timeDown = null;
    }

    function handleTouchStart(e) {
        if (e.target.getAttribute('data-swipe-ignore') === 'true') return;
        startEl = e.target;
        timeDown = Date.now();
        xDown = e.touches[0].clientX;
        yDown = e.touches[0].clientY;
        xDiff = 0;
        yDiff = 0;
    }

    function handleTouchMove(e) {
        if (!xDown || !yDown) return;
        var xUp = e.touches[0].clientX;
        var yUp = e.touches[0].clientY;
        xDiff = xDown - xUp;
        yDiff = yDown - yUp;
        e.preventDefault(); // prevents scrolling
    }

    function getNearestAttribute(el, attributeName, defaultValue) {
        while (el && el !== document.documentElement) {
            var attributeValue = el.getAttribute(attributeName);
            if (attributeValue) return attributeValue;
            el = el.parentNode;
        }
        return defaultValue;
    }

})(window, document);

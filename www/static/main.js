window.prmry = window.prmry || {}

window.prmry.AddLiveEventListener = function(eventType, elementId, cb) {
    document.addEventListener(eventType, function (event) {
        var el = event.target
            , found;

        while (el && !(found = el.id === elementId)) {
            el = el.parentElement;
        }

        if (found) {
            cb.call(el, event);
        }
    });
}

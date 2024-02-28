window.prmry = window.prmry || {}

window.prmry.AddLiveEventListener = function(eventType, elementId, cb) {
    document.addEventListener(eventType, function (event) {
        let el = event.target
            , found;

        while (el && !(found = el.id === elementId)) {
            el = el.parentElement;
        }

        if (found) {
            cb.call(el, event);
        }
    });
}

function getPromptEditor() {
    return document.querySelector("#promptEditor")
}

function handleEditorInput(event) {
    //get current cursor position
    const editor = getPromptEditor();
    const sel = window.getSelection();
    const node = sel.focusNode;
    const offset = sel.focusOffset;
    const pos = getCursorPosition(editor, node, offset, { pos: 0, done: false });
    if (offset === 0) pos.pos += 0.5;

    // update the preview
    editor.innerHTML = parse(editor.innerText);

    // restore the position
    sel.removeAllRanges();
    const range = setCursorPosition(editor, document.createRange(), {
        pos: pos.pos,
        done: false,
    });
    range.collapse(true);
    sel.addRange(range);

    // update the hidden promptInput field
    document.querySelector("#promptInput").value = editor.innerText;
}

function parse(text) {
    return text.replace(/(%s)/gm, "<span class=\"highlight\">$1</span>")
}

function getCursorPosition(parent, node, offset, stat) {
    if (stat.done) return stat;

    let currentNode = null;
    if (parent.childNodes.length === 0) {
        stat.pos += parent.textContent.length;
    } else {
        for (let i = 0; i < parent.childNodes.length && !stat.done; i++) {
            currentNode = parent.childNodes[i];
            if (currentNode === node) {
                stat.pos += offset;
                stat.done = true;
                return stat;
            } else {
                getCursorPosition(currentNode, node, offset, stat);
            }
        }
    }
    return stat;
}

function setCursorPosition(parent, range, stat) {
    if (stat.done) return range;

    let currentNode;
    if (parent.childNodes.length === 0) {
        if (parent.textContent.length >= stat.pos) {
            range.setStart(parent, stat.pos);
            stat.done = true;
        } else {
            stat.pos = stat.pos - parent.textContent.length;
        }
    } else {
        for (let i = 0; i < parent.childNodes.length && !stat.done; i++) {
            currentNode = parent.childNodes[i];
            setCursorPosition(currentNode, range, stat);
        }
    }

    return range;
}

function plainTextPaste(event) {
    event.preventDefault();
    let text = event.clipboardData.getData("text/plain");

    // Replace consecutive newlines with a single newline
    text = text.replace(/\n+/g, '\n');

    const selection = window.getSelection();
    if (!selection.rangeCount) return;

    const range = selection.getRangeAt(0);
    range.deleteContents();

    const textNode = document.createTextNode(text);
    range.insertNode(textNode);

    // Move the cursor to the end of the inserted text
    range.setStartAfter(textNode);
    range.setEndAfter(textNode);
    selection.removeAllRanges();
    selection.addRange(range);

    handleEditorInput(event);
}

function debounce(func, wait, immediate) {
    var timeout;
    return function() {
        var context = this, args = arguments;
        var later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        var callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
}

function copytext(element) {
    navigator.clipboard.writeText(element.getAttribute('data-copytext'));
    element.classList.remove('fa-copy');
    element.classList.add('fa-circle-check', 'text-success');
    setTimeout(() => {
        element.classList.remove('fa-circle-check', 'text-success');
        element.classList.add('fa-copy');
    }, 2000);
}

function selectContent(element) {
    const range = document.createRange();
    range.selectNodeContents(element);
    const selection = window.getSelection();
    selection.removeAllRanges();
    selection.addRange(range);
}

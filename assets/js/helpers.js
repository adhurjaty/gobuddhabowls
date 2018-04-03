// allow user to edit content on double click
function dblclickEdit(element) {
    element.contentEditable = true;
    setTimeout(function() {
        if (document.activeElement !== element) {
          element.contentEditable = false;
        }
    }, 300);
}

function stripSlash(s) {
    return s.replace(/\/+$/, "") + "/";
}
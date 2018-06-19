// allow user to edit content on double click
export function dblclickEdit(element) {
    element.contentEditable = true;
    setTimeout(function() {
        if (document.activeElement !== element) {
          element.contentEditable = false;
        }
    }, 300);
}

export function stripSlash(s) {
    return s.replace(/\/+$/, "") + "/";
}

// replaces the id field with the id
// e.g. /purchase_orders/{purchase_order_id} -> /purchase_orders/a5382f-448bc...
export function replaceUrlId(url, id) {
    return url.replace(/\{.*\}/, id);
}

export function dateStringToISO(d) {
    var s = d.split('/');
    var month = s[0];
    var day = s[1];
    var year = s[2];

    var date = new Date();

    return year + '-' + month + '-' + day + 'T' + date.getTimezoneOffset() / 60 + ':00:00Z';
}

export function getYearFromDateString(d) {
    return parseInt(d.split('/').pop());
}

export function formatMoney(amt) {
    return '$' + amt.toFixed(2);
}

export function formatSlashDate(dateStr) {
    var date = (new Date(dateStr)).toLocaleDateString();
    return date;
}

export function unFormatMoney(s) {
    return parseFloat(s.replace('$', ''));
}

export function getPurchaseOrderCost(po) {
    return po.Items.reduce((total, item) => {
        return total + item.price * item.count;
    }, 0) + po.shipping_cost;
}
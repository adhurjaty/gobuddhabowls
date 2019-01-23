import { format } from "path";


Number.prototype.pad = function(size) {
    var s = String(this);
    while (s.length < (size || 2)) {s = "0" + s;}
    return s;
}

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

export function rfc3339(d) {
    function pad(n) {
        return n < 10 ? "0" + n : n;
    }

    function timezoneOffset(offset) {
        var sign;
        if (offset === 0) {
            return "Z";
        }
        sign = (offset > 0) ? "-" : "+";
        offset = Math.abs(offset);
        return sign + pad(Math.floor(offset / 60)) + ":" + pad(offset % 60);
    }

    return d.getFullYear() + "-" +
        pad(d.getMonth() + 1) + "-" +
        pad(d.getDate()) + "T" +
        pad(d.getHours()) + ":" +
        pad(d.getMinutes()) + ":" +
        pad(d.getSeconds()) + 
        timezoneOffset(d.getTimezoneOffset());
}

export function getYearFromDateString(d) {
    return parseInt(d.split('/').pop());
}

export function formatMoney(amt) {
    return '$' + amt.toFixed(2);
}

export function formatSlashDate(dateStr) {
    var date = new Date(dateStr);
    var month = date.getMonth() + 1;
    var day = date.getUTCDate();
    var year = date.getFullYear();
    return `${month.pad(2)}/${day.pad(2)}/${year}`;
}

export function unFormatMoney(s) {
    return parseFloat(s.replace('$', ''));
}

export function getPurchaseOrderCost(po) {
    return po.Items.reduce((total, item) => {
        return total + item.price * item.count;
    }, 0) + po.shipping_cost;
}

// categorize groups item objects by category and sums value for each category
export function categorize(items, catItems = []) {
    return genericCategorize(items, catItems, (item, existingValue = 0) => {
        return existingValue + item.price * item.count;
    });
}

export function groupByCategory(items, catItems = []) {
    return genericCategorize(items, catItems, (item, existingValue = []) => {
        return existingValue.concat([item]);
    });
}

function genericCategorize(items, catItems, combineFnc) {
    return items.reduce((categorizedItems, item) => {
        var category = categorizedItems.find((x) => x.name == item.Category.name);

        if(category) {
            category.value = combineFnc(item, category.value);
        } else {
            categorizedItems.push({
                index: item.Category.index,
                name: item.Category.name,
                value: combineFnc(item),
                background: item.Category.background
            });
        }

        return categorizedItems;
    }, catItems).sort((a, b) => {
        return a.index - b.index;
    });
}

// getDate gets the date from a string object. Returns the beginning of the day
// i.e getDate(2018-6-20 11:40) == getDate(2018-6-20 21:12)
export function getDate(dateStr) {
    var date = new Date(dateStr);
    return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

export function sortItems(items) {
    items.sort((a, b) => {
        return (a.Category.index * 1000 + a.index) - (b.Category.index * 1000 + b.index);
    });

    return items;
}

export function toJsonName(goName) {
    var name = goName.replace(/([A-Z])/g, '_$1'.toLowerCase());
    if(name.charAt(0) == '_') {
        name = name.slice(1);
    }
    return name;
}

export function toGoName(jsName) {
    var toks = jsName.split('_');
    var goName = toks.map((tok) => tok.charAt(0).toUpperCase() 
        + tok.slice(1)).join('');
    if(goName.endsWith('Id')) {
        goName = goName.slice(0, goName.length - 2) + 'ID';
    }

    return goName;
}

export function parseModelJSON(str) {
    if(str) {
        var model = JSON.parse(str);
        if(Array.isArray(model)) {
            for(var i = 0; i < model.length; i++) {
                model[i] = stripQuotes(model[i]);
            }
        } else {
            model = stripQuotes(model);
        }

        return model;
    }
}

function stripQuotes(obj) {
    for(var key in obj) {
        if(typeof obj[key] == 'string') {
            obj[key] = stripCharacter(obj[key], '"');
        }
    }

    return obj;
}

export function stripCharacter(str, char) {
    var re = new RegExp(`^${char}+`);
    str = str.replace(re, '');
    re = new RegExp(`${char}+$`);
    return str.replace(re, '');
}

export function getObjectDiff(a, b) {
    var keys = [];
    for(var key in a) {
        if(typeof a[key] == "object") {
            keys.concat(getObjectDiff(a[key], b[key]));
        } else if(a[key] != b[key]) {
            keys.push(key);
        }
    }

    return keys;
}

export function blankUUID() {
    return "00000000-0000-0000-0000-000000000000";
}
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

// categorize groups item objects by category and sums value for each category
export function categorize(items, catItems = []) {
    return genericCategorize(items, catItems, (item, existingValue = 0) => {
        return existingValue + item.price * item.count;
    });
    // return items.reduce((categorizedItems, item) => {
    //     var value = item.price * item.count;
    //     var category = categorizedItems.find((x) => x.name == item.Category.name);

    //     if(category) {
    //         category.value += value;
    //     } else {
    //         categorizedItems.push({
    //             index: item.Category.index,
    //             name: item.Category.name,
    //             value: value,
    //             background: item.Category.background
    //         });
    //     }

    //     return categorizedItems;
    // }, catItems).sort((a, b) => {
    //     return a.index - b.index;
    // });
}

export function groupByCategory(items, catItems = []) {
    return genericCategorize(items, catItems, (item, existingValue = []) => {
        return existingValue.concat([item]);
    });
}

function genericCategorize(items, catItems, combineFnc) {
    return items.reduce((categorizedItems, item) => {
        var value = item.price * item.count;
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
import { parseModelJSON, formatSlashDate, formatMoney } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';

var _categorizedOptions = {
    breakdown: false,
    datagridUpdated: onDataGridEdit
};

var _columns = [
    {
        name: 'id',
        hidden: true,
        get_column: (item) => {
            return item.id;
        }
    },
    {
        name: 'inventory_item_id',
        hidden: true,
        get_column: (item) => {
            return item.inventory_item_id;
        }
    },
    {
        name: 'index',
        hidden: true,
        get_column: (item) => {
            return item.index;
        }
    },
    {
        header: 'Name',
        get_column: (item) => {
            return item.name;
        }
    },
    {
        name: 'selected_vendor',
        header: 'Vendor',
        editable: true,
        data_type: 'selector',
        get_column: (item) => {
            return item.selected_vendor;
        },
        options_func: (item) => {
            return Object.keys(item.VendorItemMap);
        },
        set_column: (item, value) => {
            var vendorItem = item.VendorItemMap[value];
            item.purchased_unit = vendorItem.purchased_unit;
            item.price = vendorItem.price;
            item.conversion = vendorItem.conversion;
            item.selected_vendor = value;
        }
    },
    {
        name: 'purchased_unit',
        header: 'Purchased Unit',
        editable: true,
        get_column: (item) => {
            return item.purchased_unit;
        },
        set_column: (item, value) => {
            item.purchased_unit = value;
        }
    },
    {
        name: 'price',
        header: 'Purchased Price',
        editable: true,
        data_type: 'money',
        get_column: (item) => {
            return formatMoney(item.price);
        },
        set_column: (item, value) => {
            item.price = parseFloat(value);
            item.VendorItemMap[item.selected_vendor].price = item.price;
        }
    },
    {
        name: 'conversion',
        header: 'Conversion',
        editable: true,
        data_type: 'number',
        get_column: (item) => {
            return item.conversion;
        },
        set_column: (item, value) => {
            item.conversion = parseFloat(value);
            item.VendorItemMap[item.selected_vendor].conversion = item.conversion;
        }
    },
    {
        name: 'count_unit',
        header: 'Count Unit',
        editable: true,
        get_column: (item) => {
            return item.count_unit;
        },
        set_column: (item, value) => {
            item.count_unit = value;
        }
    },
    {
        header: 'Count Price',
        get_column: (item) => {
            return formatMoney(item.price / item.conversion);
        }
    },
    {
        name: 'count',
        header: 'Count',
        editable: true,
        data_type: 'number',
        get_column: (item) => {
            return item.count;
        },
        set_column: (item, value) => {
            item.count = parseFloat(value);
        },
        default: 0
    },
    {
        header: 'Extension',
        get_column: (item) => {
            return formatMoney(item.count * item.price / item.conversion);
        }
    }
];

$(() => {
    var dateInput = $('#inventory-form').find('input[name="Date"]');
    datepicker(dateInput[0]);

    initDatagrid();
});

function initDatagrid() {
    var container = $('#categorized-items-display');
    var table = new CategorizedItemsDisplay(container, _columns, null, _categorizedOptions);
}

function onDataGridEdit(item) {
    var form = $('#inventory-form');
    var itemsInput = form.find('input[name="Items"]');
    var editedItems = [item];
    if(itemsInput.val()) {
        editedItems = JSON.parse(itemsInput.val());
        var idx = editedItems.findIndex(x => x.id == item.id);
        if(idx > -1) {
            editedItems[idx] = item;
        } else {
            editedItems.push(item);
        }
    }

    itemsInput.val(JSON.stringify(editedItems));
}
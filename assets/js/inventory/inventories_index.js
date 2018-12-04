import { parseModelJSON, formatSlashDate, formatMoney, replaceUrlId } from '../helpers/_helpers';
import { sendUpdate } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';

var _categorizedOptions = {
    breakdown: false
};

var _selectedInventory = null;

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
    var invList = $('#date-list');
    invList.find('li').first().addClass('active');

    setClickInventory();
    setSelectedInventory();
    setOnSubmit();
});

function setClickInventory() {
    var listItems = $('#date-list').find('li');
    listItems.each((i, el) => {
        var $el = $(el);
        $el.click(() => {
            listItems.removeClass('active');
            $el.addClass('active');
            setSelectedInventory();
        });
    });
}

function setSelectedInventory() {
    var $form = $('#inventory-form');
    var invList = $('#date-list');
    var container = $('#categorized-items-display');
    _selectedInventory = parseModelJSON(invList.find('li.active').attr('data'));
    container.attr('data', JSON.stringify(_selectedInventory.Items));
    var dateInput = $form.find('input[name="Date"]');
    dateInput.val(formatSlashDate(_selectedInventory.time));

    datepicker(dateInput[0]);

    var table = new CategorizedItemsDisplay(container, _columns, null, _categorizedOptions);
}

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var url = replaceUrlId(form.attr('action'), _selectedInventory.id);
        debugger;
        form.attr('action', url);
        var itemsInput = form.find('input[name="Items"]');
        var datagrid = $('#categorized-items-display');
        itemsInput.val(datagrid.attr('data'));
    })
}
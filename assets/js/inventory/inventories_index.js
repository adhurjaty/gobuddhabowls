import { parseModelJSON, formatSlashDate, formatMoney } from '../helpers/_helpers';
import { sendUpdate } from '../helpers/index_helpers';
import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';

var _categorizedOptions = {
    breakdown: false
};

var _columns = [
    {
        name: 'id',
        hidden: true,
        column_func: (item) => {
            return item.id;
        }
    },
    {
        name: 'inventory_item_id',
        hidden: true,
        column_func: (item) => {
            return item.inventory_item_id;
        }
    },
    {
        name: 'index',
        hidden: true,
        column_func: (item) => {
            return item.index;
        }
    },
    {
        header: 'Name',
        column_func: (item) => {
            return item.name;
        }
    },
    {
        name: 'selected_vendor',
        header: 'Vendor',
        editable: true,
        data_type: 'selector',
        column_func: (item) => {
            return item.selected_vendor;
        },
        options_func: (item) => {
            return Object.keys(item.VendorItemMap);
        },
        selection_func: (item, option) => {
            var vendorItem = item.VendorItemMap[option];
            item.purchased_unit = vendorItem.purchased_unit;
            item.price = vendorItem.price;
            item.conversion = vendorItem.conversion;
        }
    },
    {
        name: 'purchased_unit',
        header: 'Purchased Unit',
        editable: true,
        column_func: (item) => {
            return item.purchased_unit;
        }
    },
    {
        name: 'price',
        header: 'Purchased Price',
        editable: true,
        data_type: 'money',
        column_func: (item) => {
            return formatMoney(item.price);
        }
    },
    {
        name: 'conversion',
        header: 'Conversion',
        editable: true,
        data_type: 'number',
        column_func: (item) => {
            return item.conversion;
        }
    },
    {
        name: 'count_unit',
        header: 'Count Unit',
        editable: true,
        column_func: (item) => {
            return item.count_unit;
        }
    },
    {
        header: 'Count Price',
        column_func: (item) => {
            return formatMoney(item.price / item.conversion);
        }
    },
    {
        name: 'count',
        header: 'Count',
        editable: true,
        data_type: 'number',
        column_func: (item) => {
            return item.count;
        },
        default: 0
    },
    {
        header: 'Extension',
        column_func: (item) => {
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
    var selectedInventory = parseModelJSON(invList.find('li.active').attr('data'));
    container.attr('data', JSON.stringify(selectedInventory.Items));
    $form.find('input[name="Date"]').val(formatSlashDate(selectedInventory.time));

    var table = new CategorizedItemsDisplay(container, _columns, null, _categorizedOptions);
}

function setOnSubmit() {

}
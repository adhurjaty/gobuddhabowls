import { parseModelJSON, replaceUrlId, formatMoney, toGoName, categorize } from '../helpers/_helpers';
import { sendAjax, sendUpdate } from '../helpers/index_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { CategorizedOrderingTable } from '../components/_categorized_ordering_table';

var _categorizedOptions = {
    breakdown: false
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
        name: 'name',
        header: 'Name',
        editable: true,
        get_column: (item) => {
            return item.name;
        },
        set_column: (item, name) => {
            item.name = name;
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
            if(vendorItem != null) {
                item.purchased_unit = vendorItem.purchased_unit;
                item.price = vendorItem.price;
                item.conversion = vendorItem.conversion;
                item.selected_vendor = value;
            } else {
                item.selected_vendor = "";
            }
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
        hidden: true,
        get_column: (item) => {
            return item.count;
        }
    },
    {
        name: 'dropdown',
        get_column: (item) => {
            var editPath = replaceUrlId(_editPathBase, item.id);
            var deletePath = replaceUrlId(_deletePathBase, item.id);
            return `
            <div class="dropdown show">
                <button type="button" data-toggle="dropdown"
                    aria-haspopup="true" aria-expanded="false">
                    ...
                </button>
                <div class="dropdown-menu">
                    <a href="${editPath}" class="dropdown-item">Edit</a>
                    <a href="${deletePath}" class="dropdown-item text-danger"
                        data-method="DELETE" data-confirm="Are you sure?">
                        Delete
                    </a>
                </div>
            </div>
            `
        }
    },
    {
        name: 'recipe_unit',
        hidden: true,
        get_column: (item) => {
            return item.recipe_unit;
        }
    },
    {
        name: 'recipe_unit_conversion',
        hidden: true,
        get_column: (item) => {
            return item.recipe_unit_conversion;
        }
    },
    {
        name: 'yield',
        hidden: true,
        get_column: (item) => {
            return item.yield;
        }
    }
];

var _editPathBase = "";
var _deletePathBase = "";
var _items = [];

$(() => {
    var container = $('#categorized-items-display');
    _editPathBase = container.attr('edit-path');
    _deletePathBase = container.attr('delete-path');
    createMasterDatagrid(container);
    enableChangeOrderButton();
    setupSubmitButton();
});

function createMasterDatagrid(container) {
    _categorizedOptions.datagridUpdated = onDataGridEdit;
    return new CategorizedItemsDisplay(container, _columns, null,
        _categorizedOptions);
}

function onDataGridEdit(item) {
    var form = $('#inventory-item-form');
    var gridContainer = $('#categorized-items-display');

    var allItems = parseModelJSON(gridContainer.attr('data'));
    var oldItemIdx = allItems.findIndex(x => x.id == item.id);
    
    cleanupForm(form);

    for(var key in item) {
        var value = item[key];
        if(typeof value != 'object') {
            addInput(form, key, value);
        }
    }

    submitForm(form, item.inventory_item_id);

    allItems[oldItemIdx] = item;
    gridContainer.attr('data', JSON.stringify(allItems));
}

function cleanupForm(form) {
    form.find('input').each((i, el) => {
        if(!['authenticity_token', '_method'].includes($(el).attr('name'))) {
            $(el).remove();
        }
    });
}

function addInput(form, name, value) {
    var input = $(`<input type="text" name="${toGoName(name)}" value="${value}" style="display: none;" />`);
    form.append(input);
}

function submitForm(form, id) {
    var templatePath = form.attr('action');
    var actionPath = replaceUrlId(templatePath, id);

    form.attr('action', actionPath);
    sendAjax(form);

    form.attr('action', templatePath);
}

function enableChangeOrderButton() {
    var itemsDiv = $('#categorized-items-display');
    _items = parseModelJSON(itemsDiv.attr('data'));
    
    var table = new CategorizedOrderingTable(_items);
    var container = $('#re-order-display');
    table.attach(container);

    var button = $('#change-order-button');
    button.click(() => {
        $('#re-order-section').toggle();
        itemsDiv.toggle();
    });
}

function setupSubmitButton() {
    $('#save-order-button').click(saveInvItemsOrder);
}

function saveInvItemsOrder() {
    var $form = $('#inventory-item-form');

    $('#re-order-display').find('li[name="reorder-li"] li').each(
        function(i, el) {
            var id = $(el).attr('itemid');
            var item = _items.find(x => x.id == id);
            item.index = i;

            sendUpdate($form, item, (form) => sendAjax(form, true));
        }
    );
    
    location.reload();
}
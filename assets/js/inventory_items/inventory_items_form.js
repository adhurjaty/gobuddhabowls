import { DataGrid } from "../datagrid/_datagrid";
import { parseModelJSON, formatMoney } from "../helpers/_helpers";
import { ButtonGroup } from "../components/_button_group";
import { Modal } from "../components/_modal";

var _columnInfo = [
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
        header: 'Vendor',
        get_column: (item) => {
            return item.selected_vendor;
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
    }
];

var _datagridContainer = null;
var _datagrid = null;
var _items = [];
var _modal = null;
var _buttons = null;

$(() => {
    _datagridContainer = $('#vendor-datagrid');
    var input = $('input[name="VendorItemMap"]');
    _items = parseModelJSON(input.val()) || [];
    initDatagrid();
    initAddRemoveButtons();
    initModal();
    setOnSubmit();
});

function initDatagrid() {
    _datagrid = new DataGrid(_items, _columnInfo, datagridUpdated);
    _datagridContainer.html(_datagrid.$table);
}

function datagridUpdated(item) {
    var idx = _items.findIndex(x => x.selected_vendor == item.selected_vendor);
    _items[idx] = item;
}

function initAddRemoveButtons() {
    var container = $('#add-remove-buttons');
    _buttons = new ButtonGroup();
    _buttons.enableAddButton();
    _buttons.disableRemoveButton();
    _buttons.setRemoveListener(removeItem);
    container.html(_buttons.$content);
}

function removeItem() {
    var selectedRow = _datagridContainer.find('tr.active');

    // sometimes triggers multiple times from UI
    // this check ensures this function happens once
    if(!selectedRow) {
        return;
    }

    var selectedItem = _datagrid.getItem(selectedRow);
    _items = _items.filter(x => x.selected_vendor != selectedItem.selected_vendor);
    initDatagrid();
    _modal.addItem(selectedItem);
    if(_items.length == 0) {
        _buttons.disableRemoveButton();
    }
}

function initModal() {
    var vendorItems = parseModelJSON(_datagridContainer.attr('data'));
    var input = $('input[name="VendorItemMap"]');
    var existingItems = parseModelJSON(input.val());
    var remainingItems = vendorItems.filter(x => existingItems == null ||
        Object.keys(existingItems).findIndex(key => x.selected_vendor == key));
    remainingItems.forEach(x => {
        x.name = x.selected_vendor;
        x.id = x.selected_vendor;
    });
    _modal = new Modal(remainingItems, addItem);
    _modal.setSortFn((items) => {
        return items.sort((a, b) => (a.name > b.name) - (a.name < b.name));
    });
    var modalContainer = $('div[name="modal"]');
    modalContainer.html(_modal.$content);
}

function addItem(item) {
    _items.push(item);
    initDatagrid();
    _buttons.enableRemoveButton();
    _modal.removeItem(item);
}

function setOnSubmit() {
    var form = $('#vendor-datagrid').closest('form');
    form.submit(() => {
        var vendorMap = _datagrid.rows.reduce((obj, row) => {
            var item = row.item;
            obj[item.selected_vendor] = item;
            return obj;
        }, {});

        var input = $('input[name="VendorItemMap"]');
        input.val(JSON.stringify(vendorMap));
        debugger;        
    });
}
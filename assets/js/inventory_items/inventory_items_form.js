import { DataGrid } from "../datagrid/_datagrid";
import { parseModelJSON, formatMoney, blankUUID, isEmptyOrSpaces } from "../helpers/_helpers";
import { ButtonGroup } from "../components/_button_group";
import { Modal } from "../components/_modal";
import { showError } from "../helpers/index_helpers";

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
    }
];

var _datagridContainer = null;
var _datagrid = null;
var _items = [];
var _allItems = [];
var _modal = null;
var _buttons = null;
var _orderingTable = null;

$(() => {
    _datagridContainer = $('#vendor-datagrid');
    _allItems = parseModelJSON(_datagridContainer.attr('data')) || [];

    var input = $('input[name="VendorItemMap"]');
    _items = JSON.parse(input.val());
    if(_items) {
        _items = Object.keys(_items).map(key => {
            var val = _items[key];
            replaceAllItemSelectedVendor(key, val.selected_vendor);

            val.selected_vendor = key;
            return val;
        });
        _items.forEach(item => {
            item.name = item.selected_vendor;
            item.id = item.selected_vendor;
        });
        _datagridContainer.show();
    } else {
        _items = [];
        _datagridContainer.hide();
    }

    initDatagrid();
    initAddRemoveButtons();
    initModal();
    setOnSubmit();
});

function replaceAllItemSelectedVendor(name, id) {
    var idx = _allItems.findIndex(x => x.selected_vendor == id);
    if(idx > -1) {
        _allItems[idx].selected_vendor = name;
    }
}

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
    if(_items.length < _allItems.length) {
        _buttons.enableAddButton();
    } else {
        _buttons.disableAddButton();
    }
    if(_items.length == 0) {
        _buttons.disableRemoveButton();
    } else {
        _buttons.enableRemoveButton();
    }
    _buttons.setRemoveListener(removeItem);
    container.html(_buttons.$content);
}

function removeItem() {
    var selectedRow = _datagridContainer.find('tr.active');

    // sometimes triggers multiple times from UI
    // this check ensures this function happens once
    if(selectedRow.length == 0) {
        return;
    }

    var selectedItem = _datagrid.getItem(selectedRow);
    _items = _items.filter(x => x.selected_vendor 
        != selectedItem.selected_vendor);
    initDatagrid();
    _modal.addItem(selectedItem);
    if(_items.length == 0) {
        _buttons.disableRemoveButton();
        _datagridContainer.hide();
    }
}

function initModal() {
    var remainingItems = _allItems.filter(x => 
        _items.findIndex(item => x.selected_vendor 
            == item.selected_vendor) == -1);
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
    _datagridContainer.show();
    _buttons.enableRemoveButton();
    
    _items.push(item);
    initDatagrid();

    _modal.removeItem(item);
}

function setOnSubmit() {
    var form = $('#vendor-datagrid').closest('form');
    form.submit(() => {
        if(!validateNewItem()) {
            return false;
        }

        setVendorMapInput();
    });
}

function validateNewItem() {
    if(isEmptyOrSpaces($('input[name="Name"]').val())) {
        showError('Must enter non-blank name');
        return false;
    }
    if(isEmptyOrSpaces($('input[name="CategoryID"]').val())) {
        showError('Must enter non-blank category');
        return false;
    }
    if(_datagrid.rows.length == 0) {
        showError('Must add vendors');
        return false;
    }

    return true;
}

function setVendorMapInput() {
    var vendorMap = getVendorMapFromDatagrid();

    var input = $('input[name="VendorItemMap"]');
    input.val(JSON.stringify(vendorMap));
}

function getVendorMapFromDatagrid() {
    return _datagrid.rows.reduce((obj, row) => {
        var vendor = row.item.selected_vendor;
        var idx = _allItems.findIndex(x => x.selected_vendor 
            == vendor);
        if(idx == -1) {
            return obj;
        }
        var item = _allItems[idx];
        setAttrs(item, row.item);

        obj[item.selected_vendor] = item;
        return obj;
    }, {});
}

function setAttrs(item, rowItem) {
    item.price = rowItem.price;
    item.conversion = rowItem.conversion;
    item.purchased_unit = rowItem.purchased_unit;
    if(item.id == rowItem.name) {
        item.id = blankUUID();
    }
}

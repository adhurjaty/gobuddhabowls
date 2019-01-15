import { DataGrid } from "../datagrid/_datagrid";
import { parseModelJSON, formatMoney, groupByCategory } from "../helpers/_helpers";
import { ButtonGroup } from "../components/_button_group";
import { Modal } from "../components/_modal";
import { SingleOrderingTable } from "../components/_single_ordering_table";
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
    debugger;
    _items = JSON.parse(input.val());
    if(_items) {
        _items = Object.keys(_items).map(key => _items[key]);
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
    createOrderingTable();
    setOnChangeCategoryOrName();
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
    _items = _items.filter(x => x.selected_vendor != selectedItem.selected_vendor);
    initDatagrid();
    _modal.addItem(selectedItem);
    if(_items.length == 0) {
        _buttons.disableRemoveButton();
        _datagridContainer.hide();
    }
}

function initModal() {
    var remainingItems = _allItems.filter(x => 
        _items.findIndex(item => x.selected_vendor == item.selected_vendor));
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

function createOrderingTable() {
    var container = $('#inventory-items-display');
    var invItems = parseModelJSON(container.attr('data'));
    var item = getItem();
    
    var catItems = groupByCategory(invItems);
    var selectedCat = catItems.find(x => x.name == item.category);
    if(selectedCat) {
        var selectedCatItems = selectedCat.value;
        _orderingTable = new SingleOrderingTable(selectedCatItems, item);
        _orderingTable.attach(container);
    }
}

function getItem() {
    var category = $('select[name="CategoryID"] option:selected').html();
    var name = $('input[name="Name"]').val();
    var index = parseInt($('input[name="Index"]').val());

    return {
        name: name,
        category: category,
        index: index,
        id: ''
    };
}

function setOnChangeCategoryOrName() {
    $('select[name="CategoryID"]').change((option) => {
        clearInvItemsTable();
        createOrderingTable();
    });
    $('input[name="Name"]').change(() => {
        var name = $('input[name="Name"]').val();
        _orderingTable.updateItemName(name);
    });
}

function clearInvItemsTable() {
    $('#inventory-items-display').html('');
}

function setOnSubmit() {
    var form = $('#vendor-datagrid').closest('form');
    form.submit(() => {
        if(!validateNewItem()) {
            return false;
        }

        var vendorMap = _datagrid.rows.reduce((obj, row) => {
            var vendor = row.item.selected_vendor;
            var idx = _allItems.findIndex(x => x.selected_vendor == vendor);
            if(idx == -1) {
                return obj;
            }
            var item = _allItems[idx];
            setAttrs(item, row.item);

            obj[item.selected_vendor] = item;
            return obj;
        }, {});

        var input = $('input[name="VendorItemMap"]');
        input.val(JSON.stringify(vendorMap));

        var indexInput = $('input[name="Index"]');
        var index = findItemIndex();
        indexInput.val(index);
    });
}

function validateNewItem() {
    if(_datagrid.rows.length == 0) {
        showError('Must add vendors');
        return false;
    }

    return true;
}

function setAttrs(item, rowItem) {
    item.price = rowItem.price;
    item.conversion = rowItem.conversion;
    item.purchased_unit = rowItem.purchased_unit;
}

function findItemIndex() {
    var lis = _orderingTable.ul.find('li');
    var idx = lis.toArray().findIndex(x =>  $(x).attr('itemid') == '');
    if(idx == _orderingTable.items.length) {
        return _orderingTable.items[idx - 1].index;
    }

    return _orderingTable.items[idx].index;
}
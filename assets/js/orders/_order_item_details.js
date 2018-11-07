import { formatMoney, unFormatMoney, sortItems } from '../helpers/_helpers';
import { CategorizedDatagrid } from '../datagrid/_categorized_datagrid';
import { addToRemaining, removeFromRemaining } from '../_new_item_modal';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';

var _datagrid;
var _selected_$tr;
var _$container;
var _items;

export function initOrderItemsArea(items) {
    _$container = $('#vendor-items-table');
    _items = items;
    _items.forEach((item) => {
        if(!item.count) {
            item.count = 0;
        }
    });
    writeItemsToDOM();

    initDatagrid();
    initBreakdown();
    initAddRemoveButtons();
}

export function addToDatagrid(item) {
    // need to re-initialize globals because this is called from outside
    _$container = $('#vendor-items-table');
    _items = JSON.parse(_$container.attr('data'));

    if(item.count == undefined) {
        item.count = 0;
    }

    // add item to datagrid data
    _items.push(item);
    _items = sortItems(_items);
    writeItemsToDOM();

    initDatagrid();
    initBreakdown();

    // remove item from modal remaining
    removeFromRemaining(item);
}

function removeFromDatagrid(item) {
    // remove item from datagrid data
    var idx = _items.indexOf(item);
    _items.splice(idx, 1);
    writeItemsToDOM();

    initDatagrid();
    initBreakdown();

    // add item to modal remaining
    addToRemaining(item);

    $('#add-po-item').removeAttr('disabled');
    $('#remove-po-item').attr('disabled', true);
}

function datagridUpdated(updateObj) {
    var price = parseFloat(unFormatMoney(updateObj.price));
    var count = parseFloat(updateObj.count);
    var $tr = $('.datagrid').find(`tr td:contains(${updateObj.id})`).parent();
    $tr.find('td[name="total_cost"]').text(formatMoney(price * count));

    // update items and breakdown
    var idx = _items.findIndex((x) => x.id == updateObj.id);
    _items[idx].price = price;
    _items[idx].count = count;
    writeItemsToDOM();

    initBreakdown();
}

function initAddRemoveButtons() {
    $('#remove-po-item').click(() => {
        // sometimes triggers multiple times from UI
        // this check ensures this function happens once
        if(!_selected_$tr) {
            return;
        }
        var selectedItem = _datagrid.getItem(_selected_$tr);
        _selected_$tr = null;
        removeFromDatagrid(selectedItem);
    });
}

function writeItemsToDOM() {
    _$container.attr('data', JSON.stringify(_items));
}

function initSelection() {
    _datagrid.rows.forEach((row) => {
        var $tr = row.getRow();
        $tr.click(() => {
            $('#remove-po-item').removeAttr('disabled');
            _selected_$tr = $tr;
        });
    });
}

function initBreakdown() {
    var title = 'Order Breakdown';
    var total = _items.reduce((total, item) => total + item.count * item.price, 0);
    if(total != 0) {
        $('#category-breakdown').html(horizontalPercentageChart(title, _items, total));
    } else {
        $('#category-breakdown').html('');
    }
}

function initDatagrid() {
    var columnInfo = [
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
            name: 'price',
            header: 'Price',
            editable: true,
            data_type: 'money',
            column_func: (item) => {
                return formatMoney(parseFloat(item.price));
            }
        },
        {
            name: 'count',
            header: 'Count',
            editable: true,
            data_type: 'number',
            column_func: (item) => {
                return item.count;
            }
        },
        {
            name: 'total_cost',
            header: 'Total Cost',
            column_func: (item) => {
                return formatMoney(parseFloat(item.price) * parseFloat(item.count));
            }
        }
    ];
    _datagrid = new CategorizedDatagrid(_items, columnInfo, datagridUpdated);
    _$container.html(_datagrid.getTable());

    initSelection();
}


import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { formatMoney } from '../helpers/_helpers';

$(() => {
    setupTable();
});

var _table = null;
var _columnInfo = [
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
    }
];

function setupTable() {
    var allItemsText = $('#inventory-items').attr('data');
    if(allItemsText) {
        var allItems = JSON.parse(allItemsText);
        _table = new CategorizedItemsDisplay(_columnInfo, allItems);
    }
}
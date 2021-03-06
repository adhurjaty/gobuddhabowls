import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { formatMoney } from '../helpers/_helpers';

$(() => {
    setupTable();
    setOnFormSubmit();
});

var _table = null;
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
        header: 'Name',
        get_column: (item) => {
            return item.name;
        }
    },
    {
        name: 'price',
        header: 'Price',
        editable: true,
        data_type: 'money',
        get_column: (item) => {
            return formatMoney(parseFloat(item.price));
        },
        set_column: (item, value) => {
            item.price = parseFloat(value);
        },
        default: 0
    }
];

function setupTable() {
    var allItemsText = $('#all-inv-items').attr('data');
    if(allItemsText) {
        var allItems = JSON.parse(allItemsText);
        var container = $('#categorized-items-display');
        _table = new CategorizedItemsDisplay(container, _columnInfo, allItems, null);
    }
}

function setOnFormSubmit() {
    var form = $('#vendor-form-submit').closest('form');
    form.on('submit', (event) => {
        var data = _table.items;
        form.find('input[name="Items"]').val(JSON.stringify(data));
    });
}


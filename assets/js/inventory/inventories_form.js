import { formatMoney } from '../helpers/_helpers';
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
    var container = $('#categorized-items-display');
    var table = new CategorizedItemsDisplay(container, _columns, null, _categorizedOptions);
});
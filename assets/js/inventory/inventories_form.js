import { formatMoney } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';


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
        header: 'Name',
        get_column: (item) => {
            return item.name;
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
        set_column(item, value) {
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
    var container = $('#categorized-items-display');
    var table = new CategorizedItemsDisplay(container, _columns, null, _categorizedOptions);
});
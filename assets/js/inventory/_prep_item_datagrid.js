import { formatMoney } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from "../components/_categorized_items_display";

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
        name: 'batch_recipe_id',
        hidden: true,
        get_column: (item) => {
            return item.batch_recipe_id;
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
        name: 'conversion',
        header: 'Conversion',
        get_column: (item) => {
            return item.conversion;
        }
    },
    {
        name: 'count_unit',
        header: 'Count Unit',
        get_column: (item) => {
            return item.count_unit;
        }
    },
    {
        header: 'Count Price',
        get_column: (item) => {
            return formatMoney(item.price);
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
            return formatMoney(item.count * item.price);
        }
    }
];

export function createPrepItemDatagrid(container, onDataGridEdit) {
    if(onDataGridEdit) {
        _categorizedOptions.datagridUpdated = onDataGridEdit;
    }
    return new CategorizedItemsDisplay(container, _columns, null,
        _categorizedOptions);
}
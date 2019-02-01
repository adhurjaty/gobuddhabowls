import { parseModelJSON, formatMoney } from "../helpers/_helpers";
import { CategorizedItemsDisplay } from "../components/_categorized_items_display";

var _datagridOptions = {
    breakdown: true
};

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
        header: 'Recipe Unit',
        get_column: (item) => {
            return item.recipe_unit;
        }
    },
    {
        header: 'RU Conversion',
        get_column: (item) => {
            return item.recipe_unit;
        }
    },
    {
        header: 'RU Price',
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
        }
    },
    {
        header: 'Ext',
        get_column: (item) => {
            return formatMoney(item.price * item.count);
        }
    }
];

$(() => {
    initDatagrid();
    initOrderingTable();
});

function initDatagrid() {
    var container = $('#recipe-items-display');
    var allItems = parseModelJSON(container.attr('all-items'));
    new CategorizedItemsDisplay(container, _columnInfo, allItems,
        _datagridOptions)
}

function initOrderingTable() {

}
import { parseModelJSON, formatMoney, isEmptyOrSpaces } from "../helpers/_helpers";
import { CategorizedItemsDisplay } from "../components/_categorized_items_display";
import { showError } from "../helpers/index_helpers";

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
        name: 'recipe_unit',
        header: 'Recipe Unit',
        get_column: (item) => {
            return item.recipe_unit;
        }
    },
    {
        name: 'recipe_unit_conversion',
        header: 'RU Conversion',
        get_column: (item) => {
            return item.recipe_unit_conversion;
        }
    },
    {
        header: 'RU Price',
        get_column: (item) => {
            return formatMoney(item.price);
        }
    },
    {
        name: 'measure',
        header: 'Meas.',
        editable: true,
        get_column: (item) => {
            return item.measure;
        },
        set_column: (item, measure) => {
            item.measure = measure;
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

var _itemsDisplay = null;
var _items = [];

$(() => {
    initDatagrid();
    setOnFormSubmit();
});

function initDatagrid() {
    var container = $('#recipe-items-display');
    var allItems = parseModelJSON(container.attr('all-items'));
    _itemsDisplay = new CategorizedItemsDisplay(container, _columnInfo,
        allItems, _datagridOptions);
}

function setOnFormSubmit() {
    var form = $('#recipe-items-display').closest('form');
    form.submit(() => {
        _items = _itemsDisplay.datagrid.rows.map(x => x.item)
            .filter(x => x.count > 0);
        
        if(!validateItem()) {
            return false;
        }

        setRecipeItems();
    });
}

function validateItem() {
    if(_items.length == 0) {
        showError('Must add recipe items');
        return false;
    }
    if(isEmptyOrSpaces($('#category-input').val())) {
        showError('Must enter non-blank category');
        return false;
    }

    return true;
}

function setRecipeItems() {
    $('input[name="Items"]').val(JSON.stringify(_items));
}

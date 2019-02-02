import { parseModelJSON, formatMoney, groupByCategory } from "../helpers/_helpers";
import { CategorizedItemsDisplay } from "../components/_categorized_items_display";
import { SingleOrderingTable } from "../components/_single_ordering_table";
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
var _orderingTable = null;

$(() => {
    initDatagrid();
    setOrderingTable();
    setOnChangeCategoryOrName();
    setOnFormSubmit();
});

function initDatagrid() {
    var container = $('#recipe-items-display');
    var allItems = parseModelJSON(container.attr('all-items'));
    _itemsDisplay = new CategorizedItemsDisplay(container, _columnInfo,
        allItems, _datagridOptions)
}

function setOrderingTable() {
    var container = $('#recipe-item-ordering');
    var invItems = parseModelJSON(container.attr('data'));
    var item = getItem();

    var catItems = groupByCategory(invItems);
    var selectedCat = catItems.find(x => x.name == item.category);
    if(selectedCat) {
        var selectedCatItems = selectedCat.value;
        var idx = selectedCatItems.findIndex(x => x.id == item.id);
        if(idx > -1) {
            selectedCatItems.splice(idx, 1)
        }
        _orderingTable = new SingleOrderingTable(selectedCatItems, item);
        _orderingTable.attach(container);
    }
}

function getItem() {
    var category = $('select[name="CategoryID"] option:selected').html();
    var name = $('input[name="Name"]').val();
    var index = parseInt($('input[name="Index"]').val());
    var id = $('input[name="ID"]').val();

    return {
        name: name,
        category: category,
        index: index,
        id: id
    };
}

function setOnChangeCategoryOrName() {
    $('select[name="CategoryID"]').change((option) => {
        clearInvItemsTable();
        setOrderingTable();
    });
    $('input[name="Name"]').change(() => {
        var name = $('input[name="Name"]').val();
        _orderingTable.updateItemName(name);
    });
}

function clearInvItemsTable() {
    $('#recipe-item-ordering').html('');
}

function setOnFormSubmit() {
    var form = $('#recipe-items-display').closest('form');
    form.submit(() => {
        if(!validateItem()) {
            return false;
        }

        setIndex();
        setRecipeItems();
    });
}

function validateItem() {
    if(_itemsDisplay.datagrid.rows.length == 0) {
        showError('Must add recipe items');
        return false;
    }

    return true;
}

function setIndex() {
    var idx = findItemIndex();
    $('input[name="Index"]').val(idx);
}

function findItemIndex() {
    var id = $('input[name="ID"]').val();
    var lis = _orderingTable.ul.find('li');
    var idx = lis.toArray().findIndex(x =>  $(x).attr('itemid') == id);
    if(idx == _orderingTable.items.length) {
        return _orderingTable.items[idx - 1].index;
    }

    return _orderingTable.items[idx].index;
}

function setRecipeItems() {
    var items = _itemsDisplay.datagrid.rows.map(x => x.item);
    $('input[name="Items"]').val(JSON.stringify(items));
}

import { SingleOrderingTable } from "./_single_ordering_table";
import { groupByCategory, parseModelJSON } from "../helpers/_helpers";

var _orderingTable = null;
var _batchRecipeInput = null;
var _invItemInput = null;

$(() => {
    _batchRecipeInput = $('#recipes-selector');
    _invItemInput = $('#inventory-items-selector');

    setOrderingTable();
    setOnChangeCategoryOrName();
});

function setOrderingTable() {
    var container = $('#item-ordering-table');
    var items = parseModelJSON(container.attr('data'));
    var item = getItem();

    var catItems = groupByCategory(items);
    var selectedCat = catItems.find(x => x.name == item.category);
    if(selectedCat) {
        var selectedCatItems = selectedCat.value;
        var idx = selectedCatItems.findIndex(x => x.id == item.id);
        if(idx > -1) {
            selectedCatItems.splice(idx, 1)
        }
        _orderingTable = new SingleOrderingTable(selectedCatItems, item,
            onIndexChanged);
        _orderingTable.attach(container);
    }
}

function getItem() {
    var selectedOption = getSelectedOption();
    var name = selectedOption.html();
    var catID = selectedOption.val();
    var category = getCategory(catID);
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
    _invItemInput.change((option) => {
        clearInvItemsTable();
        setOrderingTable();
    });
    _batchRecipeInput.change((option) => {
        clearInvItemsTable();
        setOrderingTable();
    });
    $('input[name="Name"]').change(() => {
        var name = $('input[name="Name"]').val();
        if(_orderingTable) {
            _orderingTable.updateItemName(name);
        }
    });
}

function clearInvItemsTable() {
    $('#item-ordering-table').html('');
}

function onIndexChanged(evt) {
    $('input[name="Index"]').val(evt.newIndex);
}

function getActiveInput() {
    var checkbox = $('#inventory-item-check');

    if(checkbox.is(':checked')) { 
        return _invItemInput;
    }
    return _batchRecipeInput;
}

function getSelectedOption() {
    return getActiveInput().find('option:selected').first();
}

function getCategory(id) {
    var selectTag = getActiveInput();

    var items = parseModelJSON(selectTag.attr('data'));
    var matchingItem = items.find(x => x.id == id);
    return matchingItem.Category.name;
}
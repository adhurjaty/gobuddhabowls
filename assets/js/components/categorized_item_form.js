import { SingleOrderingTable } from "./_single_ordering_table";
import { groupByCategory, parseModelJSON } from "../helpers/_helpers";

var _orderingTable = null;
var _categoryInput = null;

$(() => {
    _categoryInput = $('#category-input');

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
    var category = _categoryInput.val();
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
    _categoryInput.change((option) => {
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
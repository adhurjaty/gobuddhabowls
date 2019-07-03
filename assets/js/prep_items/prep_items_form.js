import { parseModelJSON } from "../helpers/_helpers";

var _recipeSelector;
var _invItemSelector;

$(() => {
    _recipeSelector = $('#recipes-selector').closest('div');
    _invItemSelector = $('#inventory-items-selector').closest('div');
    _invItemSelector.hide();

    setOnSubmit();
    setInventoryItemCheck();
    setChangeRecipeUnit();
});

function setOnSubmit() {
    var form = $('button[role="submit"]').closest('form');
    form.submit(removeFields);
}

function removeFields() {
    var invCheck = $('#inventory-item-check');
    if(invCheck.is(':checked')) {
        _recipeSelector.remove();
    } else {
        _invItemSelector.remove();
    }
    invCheck.remove();
}

function setInventoryItemCheck() {
    var checkbox = $('#inventory-item-check');
    checkbox.change((e) => {
        toggleShowRecipe();
        if($(e.target).is(':checked'))
            $('select[name="InventoryItemID"]').trigger('change');
        else
            $('select[name="BatchRecipeID"]').trigger('change');
    });
}

function toggleShowRecipe() {
    _recipeSelector.toggle();
    _invItemSelector.toggle();
}

function setChangeRecipeUnit() {
    $('select[name="BatchRecipeID"]').change((e) => {
        setOtherUnitUI($(e.target), 'Recipe Unit', item => item.recipe_unit);
    });
    $('select[name="InventoryItemID"]').change((e) => {
        setOtherUnitUI($(e.target), "Inventory Count Unit", item => item.count_unit);
    });
}

function setOtherUnitUI(el, label, getUnit) {
    var id = el.find('option:selected').val();
    var recipes = parseModelJSON(el.attr('data'));
    var item = recipes.find(x => x.id == id);
    $('#other-unit-label').html(label);
    $('#other-unit').html(getUnit(item));
}


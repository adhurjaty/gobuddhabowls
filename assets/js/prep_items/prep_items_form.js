var _recipeSelector;
var _invItemSelector;

$(() => {
    _recipeSelector = $('#recipes-selector').closest('div');
    _invItemSelector = $('#inventory-items-selector').closest('div');
    _invItemSelector.hide();

    setOnSubmit();
    setInventoryItemCheck();
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
    checkbox.change(toggleShowRecipe);
}

function toggleShowRecipe() {
    _recipeSelector.toggle();
    _invItemSelector.toggle();
}

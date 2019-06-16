$(() => {
    setInventoryItemCheck();
});

function setOnSubmit() {
    var form = $('button[role="submit"]').closest('form');
    form.submit(removeFields);
}

function removeFields() {
    var invCheck = $('#inventory-item-check');
    if(invCheck.is(':checked')) {
        $('#recipes-selector').remove();
    } else {
        $('#inventory-items-selector').remove();
    }
    invCheck.remove();
}

function setInventoryItemCheck() {
    var checkbox = $('inventory-item-check');
    checkbox.change(toggleShowRecipe);
}

function toggleShowRecipe() {
    $('#recipes-selector').toggle();
    $('#inventory-items-selector').toggle();
}
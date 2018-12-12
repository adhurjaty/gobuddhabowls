$(() => {
    setOnSubmit();
});

function setOnSubmit() {
    var form = $('#inventory-form');
    form.submit((event) => {
        var itemsInput = form.find('input[name="Items"]');
        var datagrid = $('#categorized-items-display');
        itemsInput.val(datagrid.attr('data'));
    });
}

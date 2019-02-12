import { blankUUID } from "../helpers/_helpers";

$(() => {
    var categoryInput = $('#category-input');
    categoryInput.change((evt) => {
        var newCategory = $(evt.target).val();
        var id = blankUUID();
        var option = $('#category-datalist-options option').filter((i, x) =>
            $(x).html() == newCategory).first();
        if(option.length > 0) {
            id = option.attr('data-value');
        }
        $('input[name="CategoryID"]').val(id);
    });
});
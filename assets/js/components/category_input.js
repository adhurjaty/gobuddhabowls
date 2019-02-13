import { blankUUID } from "../helpers/_helpers";
import 'spectrum-colorpicker';

$(() => {
    setCategoryChanged();
    setupColorPicker();
});

function setCategoryChanged() {
    var categoryInput = $('#category-input');
    categoryInput.change((evt) => {
        var categoryName = $(evt.target).val();
        var option = $('#category-datalist-options option').filter((i, x) =>
            $(x).html() == categoryName).first();
        if(option.length > 0) {
            insertCategoryIDValue(option);
        } else if(categoryInput.val().length > 0) {
            insertNewCategoryValue(categoryName);
        }
    });
}

function insertCategoryIDValue(option) {
    var id = option.attr('data-value');
    $('#category-color-div').hide();
    var input = $('input[name="CategoryID"]');
    if(input.length == 0) {
        input = $('input[name="Category"]');
        input.attr('name', 'CategoryID');
    }
    input.val(id);
}

function insertNewCategoryValue(categoryName) {
    $('#category-color-div').show();
    var input = $('input[name="Category"]');
    if(input.length == 0) {
        input = $('input[name="CategoryID"]');
        input.attr('name', 'Category');
    }
    input.val(JSON.stringify({
        Name: categoryName,
        ID: blankUUID(),
        Background: $('#category-color-input').val()
    }));
}

function setupColorPicker() {
    var input = $('#category-color-input');
    input.spectrum({
        hideAfterPaletteSelect: true,
        color: input.val(),
        preferredFormat: "hex"
    });
    input.on('hide.spectrum', (e, color) => {
        var input = $('input[name="Category"]');
        var category = JSON.parse(input.val());
        category.Background = color.toHexString();
        input.val(category);
    });
}
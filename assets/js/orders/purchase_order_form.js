import { daterange } from '../_datepicker';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';


$(() => {
    var vendorItemsMap = JSON.parse($('#vendor-items-map').attr('data'));

    setDateRange();
    setDropdown(vendorItemsMap);
    setFormOnSubmit();
    setDateCheckbox();
});

var _table = null;

function setDateRange() {
    var $inputs = $('#new-order-date, #new-received-date');
    daterange($inputs);
}

function setDropdown(vendorItemsMap) {
    var dropdownId = '#new-order-vendor';
    var container = $('#categorized-items-display');
    var selectedId = $(`${dropdownId} option:selected`).val();

    if(selectedId) {
        var allItems = vendorItemsMap[selectedId];
        _table = new CategorizedItemsDisplay(allItems, container);
    }
    $(dropdownId).change(function(d) {
        // remove none option
        $(`${dropdownId} option[value=""]`).remove();

        var id = $(this).val();
        var items = vendorItemsMap[id];
        container.attr('data', JSON.stringify(items));

        _table = new CategorizedItemsDisplay(null, container);
    });
}

function setFormOnSubmit() {
    $('#purchase-order-form-submit').closest('form').submit(function(event) {
        if(!isValidOptionSelected()) {
            event.preventDefault();
            showError('must select a vendor');
            return;
        }

        var data = getDataFromTable();
        if(data.length == 0) {
            event.preventDefault();
            showError('must order at least one item');
            return;
        }
        
        removeUncheckedRecDate();
        
        sendOrderItems(data)
    });
}

function isValidOptionSelected() {
    return $('#new-order-vendor option:selected').val().length > 0;
}

function getDataFromTable() {
    return _table.items.filter(x => x.count > 0);
}

function removeUncheckedRecDate() {
    if(!$('#received-order-checkbox').is(':checked')) {
        $('#received-date-input').remove();
    }
}

function sendOrderItems(data) {
    var $input = $('form input[name="Items"]');
    $input.val(JSON.stringify(data));
}

function showError(msg) {
    $('#form-errors').text(msg);
    $('#form-errors').show();
}

function setDateCheckbox() {
    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });
}

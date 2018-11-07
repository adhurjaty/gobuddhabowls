import { initOrderItemsArea } from './_order_item_details';
import { datepicker, daterange } from '../_datepicker';


$(() => {
    var vendorItemsMap = JSON.parse($('#vendor-items-map').attr('data'));

    setDateRange();
    setDropdown(vendorItemsMap);
    setFormOnSubmit();
    setDateCheckbox();
});

function setDateRange() {
    var $inputs = $('#new-order-date, #new-received-date');
    daterange($inputs);
}

function setDropdown(vendorItemsMap) {
    var dropdownId = '#new-order-vendor';
    if($(`${dropdownId} option:selected`).val()) {
        var items = JSON.parse($('#vendor-items-table').attr('data'));
        initOrderItemsArea(items)
    }
    $(dropdownId).change(function(d) {
        // remove none option
        $(`${dropdownId} option[value=""]`).remove();

        var id = $(this).val();
        var items = vendorItemsMap[id];

        // initialize grid and breakdown
        initOrderItemsArea(items);
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
    return JSON.parse($('#vendor-items-table')
               .attr('data'))
               .filter((x) => x.count > 0);
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

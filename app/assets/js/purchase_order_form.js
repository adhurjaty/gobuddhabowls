import { initOrderItemsArea } from './order_item_details';
import { datepicker, daterange } from './datepicker';

var _vendorItemsMap;

$(() => {
    _vendorItemsMap = JSON.parse($('#vendor-items-map').attr('data'));

    var $inputs = $('#new-order-date, #new-received-date');
    daterange($inputs);

    if($('#new-order-vendor option:selected').val()) {
        var items = JSON.parse($('#vendor-items-table').attr('data'));
        initOrderItemsArea(items)
    }
    $('#new-order-vendor').change(function(d) {
        // remove none option
        $('#new-order-vendor option[value=""]').remove();
        var id = $(this).val();
        
        var items = _vendorItemsMap[id];

        // initialize grid and breakdown
        initOrderItemsArea(items);
    });

    $('#purchase-order-form-submit').closest('form').submit(function(event) {
        if($('#new-order-vendor option:selected').val().length == 0) {
            event.preventDefault();
            showError('must select a vendor');
            return;
        }

        var data = JSON.parse($('#vendor-items-table').attr('data'));
        data = data.filter((x) => x.count > 0);
        if(data.length == 0) {
            event.preventDefault();
            showError('must order at least one item');
            return;
        }
        
        if(!$('#received-order-checkbox').is(':checked')) {
            $('#received-date-input').remove();
        }
        sendOrderItems(data)
    });

    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });
});

function sendOrderItems(data) {
    var $input = $('form input[name="Items"]');
    $input.val(JSON.stringify(data));
}

function showError(msg) {
    $('#form-errors').text(msg);
    $('#form-errors').show();
}

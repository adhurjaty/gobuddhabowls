import { datepicker } from './datepicker';
import { initOrderItemsArea } from './order_item_details';

var _vendorItemsMap;

$(() => {
    _vendorItemsMap = JSON.parse($('#vendor-items-map').attr('data'));

    datepicker($('#new-order-date'), {
        autoclose: 'true',
        format: 'mm/dd/yyyy'
    });

    $('#new-order-vendor').change(function(d) {
        // remove none option
        $('#new-order-vendor option[value=""]').remove();
        var id = $(this).val();
        
        var items = _vendorItemsMap[id];

        // put the existing item edits in the vendor map
        // TODO: see if this is useful
        // cacheItemValues(lastID);

        // initialize grid and breakdown
        initOrderItemsArea(items);
    });

    $('#purchase-order-form>button[role="submit"]').click(function(event) {
        if(!$('#received-order-checkbox').is(':checked')) {
            $('#received-date-input').remove();
        }
        sendOrderItems();
    });

    $('#received-order-checkbox').change(function() {
        if(this.checked) {
            $('#received-date-input').show();
        } else {
            $('#received-date-input').hide();
        }
    });
});

function sendOrderItems() {
    var $input = $('form>input[name="Items"]');
    var data = $('#vendor-items-table').find('tr[item-id]').map(function(i, el) {
        return createItemFromRow($(el));
    }).get();
    data = JSON.stringify(data);
    $input.val(data);
}

function cacheItemValues(id) {
    var items = JSON.parse($('#vendor-items-table').attr('data'));
    _vendorItemsMap[id] = items;
    $('#vendor-items-map').attr('data', JSON.stringify(items));
}

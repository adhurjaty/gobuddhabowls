import { datepicker } from './datepicker';

$(() => {
    datepicker($('#new-order-date'), {
        autoclose: 'true',
        format: 'mm/dd/yyyy'
    });

    $('#new-order-vendor').change(function(d) {
        // remove none option
        // $('#new-order-vendor option[value=""]').remove();
        // var id = $(this).val();
        // $.ajax({
        //     url: '/purchase_orders/order_vendor_changed/' + id,
        //     method: 'GET',
        //     error: function(xhr, status, err) {
        //         var errMessage = xhr.responseText;
        //         debugger;
        //     },
        //     success: function(data, status, xhr) {
        //         initDatagrid();
        //     }
        // });
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

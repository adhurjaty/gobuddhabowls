// import { datepicker } from 'bootstrap-datepicker';

$(() => {
    $('#new-order-date').datepicker({
        autoclose: 'true',
        format: 'mm/dd/yyyy'
    });

    $('#new-order-vendor').change(function(d) {
        var id = $(this).val();
        $.ajax({
            url: '/purchase_orders/order_vendor_changed/' + id,
            method: 'GET',
            error: function(xhr, status, err) {
                var errMessage = xhr.responseText;
                debugger;
            }
            // success: function(data, status, xhr) {
            //     editItem.onUpdateSuccess();
            // }
        });
    });
});
// import { datepicker } from 'bootstrap-datepicker';
import { EditItem, DataGrid } from './datagrid.js';

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

    // $('#vendor-items-table').on('DOMNodeInserted', function(event) {
    //     if(event.target.parentNode.id == 'vendor-items-table') {
    //         $.each($('.datagrid'), function(i, grid) {
    //             var dg = new DataGrid(grid);
        
    //             $.each($(this).find('td[editable="true"]'), function(j, el) {
    //                 var ei = new EditItem(dg, $(el));
    //             });
    //         });
    //     }
    // });
});
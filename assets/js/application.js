require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap-sass/assets/javascripts/bootstrap.js");
require("bootstrap-table/src/bootstrap-table.js");
require("bootstrap");
require("bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js");

$(() => {
      
});

// DATAGRID

$.each($('.datagrid').find('td[contenteditable="true"]'), function(i, el) {
    var type = $(el).attr('data-type');

    var sendUpdate = function(po_id, field, contents) {
        var data = {};
        data[field] = contents;
        $.ajax({
            url: "/purchase_orders/row_edited/" + po_id,
            data: data,
            method: "POST"
        });
    };

    switch(type) {
        case 'date':
            $(el).on('focus', function(event) {
                var date_str = $(this).text();
                var $date = $('<input data-provide="datepicker" value="' + date_str + '">');
                $(this).empty();
                $(this).append($date);
                var $td = $(this)
                $date.datepicker({
                    autoclose: 'true',
                    format: 'mm/dd/yyyy',
                    defaultViewDate: date_str
                }).on('changeDate', function(event) {
                    var po_id = $td.parent().attr('item-id');
                    var field = $td.attr('field');
                    var contents = event.format();

                    sendUpdate(po_id, field, contents);
                    
                    $td.text(contents);
                }).on('hide', function(event) {
                    $td.text(date_str);
                });
                $date.focus();
            });
            break;
        case 'money':
            break;
        case 'selector':
            break;
        default:    // type 'text'
            $(el).on('blur', function(event) {
                var po_id = $(this).parent().attr('item-id');
                var field = $(this).attr('field');
                sendUpdate(po_id, field, $(this).text());
            });
            break;
    }
    
    
    // switch(type) {
    //     case 'date':
    //         el = '<input data-provide="datepicker" format="mm/dd/yyyy" startDate="' + $td.text + '">'
    //         break;
    //     case 'money':
    //         break;
    //     case 'selector':
    //         break;
    //     default:    // type 'text'
    //         el = '<input type="text" class="form-control" value="' + prevValue + '" onfocus="this.select()">'
    //         $td.append(el);
    //         $(el).focus().select();
    // }
});
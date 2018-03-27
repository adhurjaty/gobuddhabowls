require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap-sass/assets/javascripts/bootstrap.js");
require("bootstrap-table/src/bootstrap-table.js");
require("bootstrap");
require("bootstrap-datepicker");
require("bootstrap-datagrid/js/bootstrap-datagrid.js"); // must be at the bottom
$(() => {
    
    $('#datagrid').datagrid({
        inputs: {
            select: {
                el : $('<select class="form-control datagrid-input">'),
                onShow:function(cell) {
                // Set the options
                if (!$(this).find('option').length) {
                    $(this).append($('<option disabled="disabled">Select a category</option>'))
                    $(this).append($('<option value="shoe">Shoe</option>'))
                    $(this).append($('<option value="t-shirt">T-Shirt</option>'))
                    $(this).append($('<option value="pants">Pants</option>'))
                }
        
                var inputPadding = parseInt(cell.data('padding'))-1
                $(this).css('padding', inputPadding+'px')
                $(this).css('width', '100%')
                $(this).css('height', '100%')
                $(this).css('top', cell.offset().top.toString+'px')
                $(this).css('left', cell.offset().left.toString+'px')
        
                $(this).val(cell.data('value'))
                },
                onChange:function(cell) {
                cell.data('value', $(this).val())
                cell.text($(this).find('option[value='+$(this).val()+']').text())
                },
                isChanged:function(cell) {
                return $(this).val() != cell.data('value')
                }
            },
            money: {
                el : $('<input type="text" class="form-control datagrid-input">'),
                onShow:function(cell) {
                    var inputPadding = parseInt(cell.data('padding'))-1
                    $(this).css('padding', inputPadding+'px')
                    $(this).css('width', '100%')
                    $(this).css('height', '100%')
                    $(this).css('top', cell.offset().top.toString+'px')
                    $(this).css('left', cell.offset().left.toString+'px')
            
                    $(this).val(cell.data('value'))
                },
                onChange:function(cell) {
                    cell.data('value', $(this).val())
                    cell.text('$'+$(this).val())
                },
                isChanged:function(cell) {
                    return $(this).val() != cell.data('value')
                }
            },
            date: {
                el: $('<input class="datepicker" data-date-format="mm/dd/yyyy">'),
                onShow: function(cell) {
                    $('.datepicker').datepicker({
                        startDate: "-3d"
                    })
                },
                onChange: function(cell) {
                    cell.text($(this).val())
                },
                isChanged: function(cell) {
                    return true
                }
            }
        }
    });
});


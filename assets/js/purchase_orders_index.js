var datagrid = require("./datagrid.js");
var barChart = require("./vertical_bar_chart.js");

$(() => {
    $('#datagrid-holder').on('DOMNodeInserted', function(event) {
        if(event.target.parentNode.id == 'datagrid-holder') {
            $.each($('.datagrid'), function(i, grid) {
                var dg = new datagrid.DataGrid(grid);
        
                $.each($(this).find('td[editable="true"]'), function(j, el) {
                    var ei = new datagrid.EditItem(dg, $(el));
                });
            });
        }
    });

    $('#summary-table').on('DOMNodeInserted', function(event) {
        // only run if a node has been inserted at the top level
        if(event.target.parentNode.id == 'summary-table') {
            var $el = $(this).children().first();
            if($el.length > 0) {
                var height = parseInt($el.attr('height'));
                var data = unescape($el.attr('data'));

                var verticalChart = new barChart.VerticalBarChart(height, data, 'summary-table');
            }
        }
    });
});
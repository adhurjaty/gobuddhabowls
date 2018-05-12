import { EditItem, DataGrid } from './datagrid.js';
import { VerticalBarChart } from './vertical_bar_chart.js';
import { MultilineGraph } from './multiline_graph.js';
import { replaceUrlId } from './helpers.js';

$(() => {
    $('#datagrid-holder').on('DOMNodeInserted', function(event) {
        // only run if a node has been inserted at the top level
        if(event.target.nodeName == 'DIV' && event.target.parentNode.id == 'datagrid-holder') {
            $.each($('.datagrid'), function(i, grid) {
                var dg = new DataGrid(grid, datagridUpdate);
        
                $.each($(this).find('td[editable="true"]'), function(j, el) {
                    var ei = new EditItem(dg, $(el));
                });
            });
        }
    });

    $('#summary-table').on('DOMNodeInserted', function(event) {
        if(event.target.parentNode.id == 'summary-table') {
            var $el = $(this).children().first();
            if($el.length > 0) {
                var height = parseInt($el.attr('height'));
                var data = unescape($el.attr('data'));

                var verticalChart = new VerticalBarChart(height, data, 'summary-table');
            }
        }
    });

    $('#trend-chart').on('DOMNodeInserted', function(event) {
        if(event.target.parentNode.id == 'trend-chart') {
            var $el = $(this).children().first();
            if($el.length > 0) {
                var height = parseInt($el.attr('height'));
                var data = unescape($el.attr('data'));

                var graph = new MultilineGraph(height, data, 'trend-chart');
            }
        }
    });
});

function datagridUpdate(editItem) {
    var data = {};
    // HACK: not crazy about hard-coding the path. May be a way to get
    // plush templating i.e. <%= purchaseOrdersPath() %>
    var on_change_url = '/purchase_orders/';
    data[editItem.field] = editItem.contents;
    $.ajax({
        url: on_change_url + editItem.id,
        data: data,
        method: 'PUT',
        error: function(xhr, status, err) {
            var errMessage = xhr.responseText;
            editItem.showError(errMessage);
        },
        success: function(data, status, xhr) {
            editItem.onUpdateSuccess();
        }
    });
}
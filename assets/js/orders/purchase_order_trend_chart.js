import { MultilineGraph } from '../_multiline_graph'
import { categorize, getDate } from '../_helpers'


$(() => {
    var $chartContainer = $('#trend-chart');
    var height = parseInt($chartContainer.attr('data-height'));
    var purchaseOrders = JSON.parse($chartContainer.attr('data-items'));

    var trendChartData = getLineData(purchaseOrders);

    new MultilineGraph(height, trendChartData, 'trend-chart');
});

function getLineData(purchaseOrders) {
    return purchaseOrders.reduce((arr, po) => {
        arr.push(...categorize(po.Items).map((item) => {
            item.date = getDate(po.order_date);
            return item;
        }));
        return arr;
    }, [])
    .sort((a, b) => {
        return a.index - b.index;
    });
}
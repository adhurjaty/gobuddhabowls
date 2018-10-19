
import { Chart } from 'chart.js'
import { categorize, formatMoney } from './helpers';


$(() => {
    var $el = $('#summary-table');
    var height = parseInt($el.attr('data-height'));
    $el.height(height);
    var canvasId = 'bar-chart';
    $el.html(`<canvas id="${canvasId}"></canvas>`)
    var purchaseOrders = JSON.parse($el.attr('data-items'));

    var categorizedData = purchaseOrders.reduce((cat, po) => {
        return categorize(po.Items, cat);
    }, []);
    var labels = categorizedData.map((o) => o.name);
    var data = categorizedData.map((o) => o.value);
    var backgrounds = categorizedData.map((o) => o.background);

    var ctx = document.getElementById(canvasId);
    ctx.height = height;
    var barChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: backgrounds,
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            tooltips: {
                callbacks: {
                    label: (item, chart) => {
                        var amt = parseFloat(item.yLabel);
                        return formatMoney(amt);
                    }
                }
            },
            legend: {
                display: false
            },
            layout: {
                padding: {
                    top: 30
                }
            }
        }
    });
    // new VerticalBarChart(height, categorizedData, 'summary-table');
});
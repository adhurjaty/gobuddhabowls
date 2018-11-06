import { formatMoney } from "../_helpers";
import { OrderSheetTable } from "./_order_sheet_table";
import html2canvas from 'html2canvas';

$(() => {
    var columnInfo = [
        {
            name: 'Pack Size',
            column_func: (item) => {
                return item.conversion;
            }
        },
        {
            name: 'Purchased Unit',
            column_func: (item) => {
                return item.purchased_unit;
            }
        },
        {
            name: 'category',
            column_func: (item) => {
                return item.name;
            }
        },
        {
            name: 'Order Amt',
            column_func: (item) => {
                return item.count;
            }
        },
        {
            name: 'Current Price',
            column_func: (item) => {
                return formatMoney(item.price);
            }
        },
        {
            name: 'Extension',
            column_func: (item) => {
                return formatMoney(item.count * item.price);
            }
        }
    ];

    var $container = $('#order-table');
    var data = JSON.parse($container.attr('data'));
    var table = new OrderSheetTable(data, columnInfo);
    $container.append(table.getTable());

    html2canvas(document.body).then((canvas) => {
        document.body.appendChild(canvas);
        $('#order-sheet').remove();
    });
});
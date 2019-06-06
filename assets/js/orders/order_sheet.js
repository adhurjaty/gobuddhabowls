import { formatMoney } from "../helpers/_helpers";
import { OrderSheetTable } from "./_order_sheet_table";
import html2canvas from 'html2canvas';

$(() => {
    var columnInfo = [
        {
            name: 'Pack Size',
            get_column: (item) => {
                return item.conversion;
            }
        },
        {
            name: 'Purchased Unit',
            get_column: (item) => {
                return item.purchased_unit;
            }
        },
        {
            name: 'category',
            get_column: (item) => {
                return item.name;
            }
        },
        {
            name: 'Order Amt',
            get_column: (item) => {
                return item.count;
            }
        },
        {
            name: 'Current Price',
            get_column: (item) => {
                return formatMoney(item.price);
            }
        },
        {
            name: 'Extension',
            get_column: (item) => {
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
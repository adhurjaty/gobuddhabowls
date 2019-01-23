import { parseModelJSON, formatMoney } from "../helpers/_helpers";
import { TotalDatagrid } from "../datagrid/_total_datagrid";

var _columnInfo = [
    {
        header: 'Name',
        get_column: (item) => {
            return item.name;
        }
    },
    {
        header: 'Count',
        get_column: (item) => {
            return item.count;
        }
    },
    {
        header: 'Amount Collected',
        get_column: (item) => {
            return formatMoney(item.amount);
        }
    }
]

$(() => {
    initDatagrid();
});


function initDatagrid() {
    var container = $('#sales-datagrid');
    var salesSummary = parseModelJSON(container.attr('data'));
    salesSummary.Sales.sort((a, b) => a.name < b.name ? -1 : 1);
    var datagrid = new TotalDatagrid(salesSummary.Sales, _columnInfo);
    addSummaryRows(datagrid, salesSummary);
    
    container.html(datagrid.$table);
}

function addSummaryRows(datagrid, summary) {
    datagrid.addTotalRow("Tips", formatMoney(summary.tips));
    datagrid.addTotalRow("Tax", formatMoney(summary.tax));
    datagrid.addTotalRow("Fees", formatMoney(summary.fees));
    datagrid.addTotalRow("Refunds", formatMoney(summary.refunds));

    var salesTotal = summary.Sales.reduce((sum, x) => sum + x.amount, 0)
     + summary.tips;
    datagrid.addTotalRow("Total", formatMoney(salesTotal));
    var net = salesTotal - summary.tax - summary.fees - summary.refunds;
    datagrid.addTotalRow("Net", formatMoney(net));
}


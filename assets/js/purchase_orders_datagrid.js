import { DataGrid, EditItem } from './datagrid';
import { getPurchaseOrderCost } from './helpers';
import { horizontalPercentageChart } from './horizontal_percentage_chart';

$(() => {
    var purchaseOrders = JSON.parse($('#purchase-orders-holder').val());
    var openOrders = purchaseOrders.filter(x => !x.received_date);
    var recOrders = purchaseOrders.filter(x => x.received_date);
    var table = "";

    if(openOrders.length > 0) {
        table = getDataGrid("Open Orders", openOrders);
    }
    if(recOrders.length > 0) {
        table += getDataGrid("Received Orders", recOrders);
    }

    $('#datagrid-holder').html(table);
    
    $.each($('.datagrid'), function(i, grid) {
        var dg = new DataGrid(grid);

        $.each($(this).find('td[editable="true"]'), function(j, el) {
            var ei = new EditItem(dg, $(el));
        });
    });
});

function getDataGrid(title, purchaseOrders) {
    var head = `
    <div class="row justify-content-center">
        <div class="col-6 text-center">
            <h4>${title}</h4>
        </div>
    </div>
    <div class="row">
    <table class="datagrid" update-url="<%= purchaseOrdersPath() %>">
        <thead>
            <th></th>
            <th>Vendor</th>
            <th>Order Date</th>
            <th>Cost</th>
            <th>&nbsp;</th>
        </thead>
    <tbody>`;

    var foot = `
    </tbody>
      </table>
    </div>`;
    
    var rows = purchaseOrders.map(po => {
        var total = getPurchaseOrderCost(po)
        return `
        <tr item-id="${po.id}">
            <td class="expander"><span class="fa fa-caret-right"></span></td>
            <td>${po.Vendor.name}</td>
            <td editable="true" data-type="date" field="OrderDate"><%= format_date(purchaseOrder.OrderDate) %></td>
            <td>${total}</td>
            <td>
                <div class="dropdown show">
                <button type="button"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    ...
                </button>
                <div class="dropdown-menu">
                    <span class="dropdown-item" onclick="receiveItem(${po.id})">Received</span>
                    <a href="<%= editPurchaseOrderPath({ purchase_order_id: ${po.id}}) %>" class="dropdown-item">Edit</a>
                    <span class="dropdown-item text-danger" onclick="deleteItem(${po.id})">Delete</span>
                </div>
                </div>
            </td>
        </tr>
        <tr class="items-list" style="display: none;">
            <td colspan="100">
                ${horizontalPercentageChart(title, po.Items, total)}            
            </td> 
        </tr>`;
    });
    
    return head + rows + foot;
}
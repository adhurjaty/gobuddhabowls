import { DataGrid, EditItem } from './datagrid';
import { getPurchaseOrderCost } from './helpers';

$(() => {
    var purchaseOrders = $('#purchase-orders-holder').data();
    var openOrders = purchaseOrders.filter(x => !x.ReceivedDate.Valid);
    var recOrders = purchaseOrders.filter(x => x.ReceivedDate.Valid);
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
            <h4>{0}</h4>
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
    <tbody>`.format(title);

    var foot = `
    </tbody>
      </table>
    </div>`;
    
    var rows = purchaseOrders.map(po => {
        return `
        <tr item-id="{0}">
            <td class="expander"><span class="fa fa-caret-right"></span></td>
            <td>{1}</td>
            <td editable="true" data-type="date" field="OrderDate"><%= format_date(purchaseOrder.OrderDate) %></td>
            <td>{2}</td>
            <td>
                <div class="dropdown show">
                <button type="button"  data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    ...
                </button>
                <div class="dropdown-menu">
                    <span class="dropdown-item" onclick="receiveItem({0})">Received</span>
                    <a href="<%= editPurchaseOrderPath({ purchase_order_id: {0}}) %>" class="dropdown-item">Edit</a>
                    <span class="dropdown-item text-danger" onclick="deleteItem({0})">Delete</span>
                </div>
                </div>
            </td>
        </tr>
        <tr class="items-list" style="display: none;">
            <td colspan="100">
                <%= partial("partials/horizontal_percentage_chart.html", {categoryDetails: purchaseOrder.GetCategoryCosts(), title: ""}) %>              
            </td> 
        </tr>`.format(po.id, po.vendor.name, getPurchaseOrderCost(po));
    });
    
    return head + rows + foot;
}
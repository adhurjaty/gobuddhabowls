import { CollapsibleDatagrid } from './collapsible_datagrid';
import { getPurchaseOrderCost, formatMoney, formatSlashDate, replaceUrlId } from './helpers';
import { horizontalPercentageChart } from './horizontal_percentage_chart';

$(() => {
    var $container = $('#datagrid-holder');
    // var purchaseOrders = JSON.parse($container.attr('data'));
    var editOrderPath = $container.attr('data-url');
    var $openDatagridContainer = $('#open-order-datagrid')
    var openOrders = JSON.parse($openDatagridContainer.attr('data'));
    var $recDatagridContainer = $('#rec-order-datagrid');
    var recOrders = JSON.parse($recDatagridContainer.attr('data'));

    var baseColumnObjects = [
        {
            name: 'id',
            hidden: true,
            column_func: (purchaseOrder) => {
                return purchaseOrder.id;
            }
        },
        {
            name: 'vendor',
            header: 'Vendor',
            column_func: (purchaseOrder) => {
                return purchaseOrder.Vendor.name;
            }
        },
        {
            name: 'order_date',
            header: 'Order Date',
            editable: true,
            data_type: 'date',
            column_func: (purchaseOrder) => {
                return formatSlashDate(purchaseOrder.order_date);
            }
        },
        {
            name: 'cost',
            header: 'Cost',
            column_func: (purchaseOrder) => {
                return formatMoney(getPurchaseOrderCost(purchaseOrder));
            }
        }
    ];

    var openColumnObjects = baseColumnObjects.concat([
        {
            name: 'dropdown',
            column_func: ((editOrderPath) => {
                return (purchaseOrder) => {
                    return `<div class="dropdown show">
                        <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            ...
                        </button>
                        <div class="dropdown-menu">
                            <span class="dropdown-item" onclick="receiveItem(${purchaseOrder.id})">Received</span>
                            <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                            <span class="dropdown-item text-danger" onclick="deleteItem(${purchaseOrder.id})">Delete</span>
                        </div>
                    </div>`
                }
            })(editOrderPath)
        }
    ]);

    var recColumnObjects = baseColumnObjects.slice();
    recColumnObjects.splice(3, 0, {
        name: 'received_date',
        header: 'Received Date',
        editable: true,
        data_type: 'date',
        column_func: (purchaseOrder) => {
            return formatSlashDate(purchaseOrder.received_date);
        }
    });
    recColumnObjects.push({
        name: 'dropdown',
        column_func: ((editOrderPath) => {
            return (purchaseOrder) => {
                return `<div class="dropdown show">
                    <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        ...
                    </button>
                    <div class="dropdown-menu">
                        <span class="dropdown-item" onclick="reopen(${purchaseOrder.id})">Re-open</span>
                        <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                        <span class="dropdown-item text-danger" onclick="deleteItem(${purchaseOrder.id})">Delete</span>
                    </div>
                </div>`
            }
        })(editOrderPath)
    });

    var openDatagrid = new CollapsibleDatagrid(openOrders, openColumnObjects, getHiddenRow, sendUpdate)
    var recDatagrid = new CollapsibleDatagrid(recOrders, recColumnObjects, getHiddenRow, sendUpdate);

    $openDatagridContainer.html(openDatagrid.getTable());
    $recDatagridContainer.html(recDatagrid.getTable());

    // var openOrders = purchaseOrders.filter(x => !x.received_date);
    // var recOrders = purchaseOrders.filter(x => x.received_date);
    // var tableArea = "";

    // if(openOrders.length > 0) {
    //     tableArea = getDataGrid("Open Orders", openOrders, editOrderPath);
    // }
    // if(recOrders.length > 0) {
    //     tableArea += getDataGrid("Received Orders", recOrders, editOrderPath);
    // }

    // $('#datagrid-holder').html(tableArea);
    
    // $.each($('.datagrid'), function(i, grid) {
    //     new DataGrid(grid);
    // });
});

// function getDataGrid(title, purchaseOrders, editOrderPath) {
//     var head = `
//     <div class="row justify-content-center">
//         <div class="col-6 text-center">
//             <h4>${title}</h4>
//         </div>
//     </div>
//     <div class="row">
//     <table class="datagrid" update-url="<%= purchaseOrdersPath() %>">
//         <thead>
//             <th></th>
//             <th>Vendor</th>
//             <th>Order Date</th>
//             <th>Cost</th>
//             <th>&nbsp;</th>
//         </thead>
//     <tbody>`;

//     var foot = `
//     </tbody>
//       </table>
//     </div>`;
    
//     var rows = purchaseOrders.map(po => {
//         var total = getPurchaseOrderCost(po)
//         return `
//         <tr item-id="${po.id}">
//             <td class="expander"><span class="fa fa-caret-right"></span></td>
//             <td>${po.Vendor.name}</td>
//             <td editable="true" data-type="date" field="OrderDate">${formatSlashDate(po.order_date)}</td>
//             <td>${formatMoney(total)}</td>
//             <td>
//                 <div class="dropdown show">
//                     <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
//                         ...
//                     </button>
//                     <div class="dropdown-menu">
//                         <span class="dropdown-item" onclick="receiveItem(${po.id})">Received</span>
//                         <a href="${replaceUrlId(editOrderPath, po.id)}" class="dropdown-item">Edit</a>
//                         <span class="dropdown-item text-danger" onclick="deleteItem(${po.id})">Delete</span>
//                     </div>
//                 </div>
//             </td>
//         </tr>
//         <tr class="items-list" style="display: none;">
//             <td colspan="100">
//                 ${horizontalPercentageChart(title, po.Items, total)}            
//             </td> 
//         </tr>`;
//     }).join('');
    
//     return head + rows + foot;
// }

function sendUpdate(updateObj) {
    var id = updateObj.id;
    var updateUrl = replaceUrlId($('#datagrid-holder').attr('update-url'), id);
    var data = {};
    data[updateObj.name] = updateObj.value;
    $.ajax({
        url: updateUrl,
        data: data,
        method: 'PUT',
        error: function(xhr, status, err) {
            var errMessage = xhr.responseText;
            debugger;
            updateObj.onError(errMessage);
        },
        success: function(data, status, xhr) {
            updateObj.onSuccess();
        }
    });
}

function getHiddenRow(purchaseOrder) {
    return `<tr class="items-list">
                <td colspan="100">
                    ${horizontalPercentageChart(purchaseOrder.Vendor.name + ' Order', purchaseOrder.Items, getPurchaseOrderCost(purchaseOrder))}            
                </td> 
            </tr>`;
}
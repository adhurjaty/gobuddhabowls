import { CollapsibleDatagrid } from '../datagrid/_collapsible_datagrid';
import { getPurchaseOrderCost, formatMoney, formatSlashDate, replaceUrlId, toGoName } from '../helpers/_helpers';
import { horizontalPercentageChart } from '../_horizontal_percentage_chart';
import { sendUpdate, submitUpdateForm } from '../helpers/index_helpers';

const CONTAINER_ID = '#datagrid-holder',
      OPEN_ORDER_DG_ID = '#open-order-datagrid',
      REC_ORDER_DG_ID = '#rec-order-datagrid',
      EDIT_PATH_ATTR = 'data-url',
      ORDER_SHEET_PATH_ATTR = 'order-sheet-url',
      RECEIVING_PATH_ATTR = 'receiving-list-url';

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
    
$(() => {
    var $container = $(CONTAINER_ID);

    populateOpenDatagrid($container);
    populateRecDatagrid($container);
});

function populateOpenDatagrid($container) {
    var specialOption = `<span name="receive" class="dropdown-item">Received</span>`;
    var colInfo = getColInfo($container, specialOption);
    createDatagrid(OPEN_ORDER_DG_ID, colInfo);
}

function populateRecDatagrid($container) {
    var specialOption = `<span name="reopen" class="dropdown-item">Re-open</span>`;
    var colInfo = getColInfo($container, specialOption);
    colInfo.splice(3, 0, {
        name: 'received_date',
        header: 'Received Date',
        editable: true,
        data_type: 'date',
        column_func: (purchaseOrder) => {
            return formatSlashDate(purchaseOrder.received_date);
        }
    });
    createDatagrid(REC_ORDER_DG_ID, colInfo);
}

function getColInfo($container, specialOption) {
    var editOrderPath = $container.attr(EDIT_PATH_ATTR);
    var orderSheetPath = $container.attr(ORDER_SHEET_PATH_ATTR);
    var receivingListPath = $container.attr(RECEIVING_PATH_ATTR);
    return baseColumnObjects.concat([
        {
            name: 'dropdown',
            column_func: ((editOrderPath, orderSheetPath, receivingListPath) => {
                return (purchaseOrder) => {
                    return `<div class="dropdown show">
                        <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            ...
                        </button>
                        <div class="dropdown-menu">
                            ${specialOption}
                            <a href="${replaceUrlId(editOrderPath, purchaseOrder.id)}" class="dropdown-item">Edit</a>
                            <a name="delete" class="dropdown-item text-danger" data-method="DELETE">Delete</a>
                            <a href="${replaceUrlId(orderSheetPath, purchaseOrder.id)}" class="dropdown-item">Order Sheet<a>
                            <a href="${replaceUrlId(receivingListPath, purchaseOrder.id)}" class="dropdown-item">Receiving List<a>
                        </div>
                    </div>`
                };
            })(editOrderPath, orderSheetPath, receivingListPath)
        }
    ]);
}

function createDatagrid(containerId, colInfo) {
    var dataStr = $(containerId).attr('data');
    if(dataStr) {
        var data = JSON.parse(dataStr);
        var datagrid = new CollapsibleDatagrid(data, colInfo, getHiddenRow, sendDatagridUpdate);
        $(containerId).html(datagrid.getTable());
        initDropdownActions(datagrid);
    }
}

function getHiddenRow(purchaseOrder) {
    return `<tr class="items-list">
                <td colspan="100">
                    ${horizontalPercentageChart(purchaseOrder.Vendor.name + ' Order', purchaseOrder.Items, getPurchaseOrderCost(purchaseOrder))}            
                </td> 
            </tr>`;
}

function initDropdownActions(datagrid) {
    var updatePath = $('#datagrid-holder').attr('update-url');
    var $form = $('#update-po-form');
    datagrid.rows.forEach((row) => {
        var details = row.getInfo();
        var id = details.id;
        var $row = row.getRow();

        $row.find('td[name="dropdown"] span[name="receive"]').click(() => {
            var date = new Date();
            date.setHours(0, 0, 0, 0);
            submitUpdateForm($form, id, { ReceivedDate: date.toISOString() });
        });
        $row.find('td[name="dropdown"] span[name="reopen"]').click(() => {
            submitUpdateForm($form, id, { ReceivedDate: '' });
        });
        $row.find('td[name="dropdown"] a[name="delete"]').attr('href', replaceUrlId(updatePath, id));
    });
}

function sendDatagridUpdate(updateObj) {
    sendUpdate($('#update-po-form'), updateObj);
}

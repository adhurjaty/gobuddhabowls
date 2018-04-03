$('#datagrid-holder').empty();
$('#datagrid-holder').append('<%= partial("purchase_orders/datagrid.html") %>');
$('#period-selector-component').empty();
$('#period-selector-component').append('<%= partial("partials/period_selector.html") %>');
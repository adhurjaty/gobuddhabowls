import { parseModelJSON } from "../helpers/_helpers";

// var _columnInfo = [
//     {
//         header: 'Name',
//         get_column: (item) => {
//             item.name;
//         }
//     },
//     {
//         header: ''
//     }
// ]

$(() => {
    initDatagrid();
});


function initDatagrid() {
    var sales = parseModelJSON($('#sales-datagrid').attr('data'));
    debugger;
}
import { parseModelJSON } from "../helpers/_helpers";

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
            return item.count
        }
    }
]

$(() => {
    initDatagrid();
});


function initDatagrid() {
    var sales = parseModelJSON($('#sales-datagrid').attr('data'));
    debugger;
}
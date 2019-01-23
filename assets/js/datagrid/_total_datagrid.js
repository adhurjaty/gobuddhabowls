import { DataGrid } from "./_datagrid";

export class TotalDatagrid extends DataGrid {
    constructor(data, columnInfo) {
        super(data, columnInfo);
    }

    addTotalRow(name, value) {
        var tr = $('<tr></tr>');
        var nameTd = $(`<td colspan="${this.columnInfo.length - 1}">
            ${name}
        </td>`);
        var valTd = $(`<td>${value}</td>`);
        tr.append(nameTd);
        tr.append(valTd);
        this.$table.append(tr);
    }
}
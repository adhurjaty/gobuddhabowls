import { DataGrid } from "../datagrid/_datagrid";
import { groupByCategory, formatMoney } from "../_helpers";

export class OrderSheetTable extends DataGrid {
    constructor(data, columnInfo) {
        super(data, columnInfo, (p) => {});
        this.categorizedData = groupByCategory(data);
        this.insertHeadersAndTotals();
        this.insertTotalsSection();
    }

    insertHeadersAndTotals() {
        var self = this;
        var i = 0;
        this.categorizedData.forEach((catItem) => {
            var $row = this.rows[i].getRow();
            var $headerRow = $(`<tr class="category-header" style="background-color: ${catItem.background};"></tr>`);
            self.columnInfo.forEach((info) => {
                var $td = $('<td></td>');
                if(info.name == 'category') {
                    $td.text(catItem.name);
                } else {
                    $td.text(info.name);
                }
                $td.appendTo($headerRow);
            });

            $headerRow.insertBefore($row);
            i += catItem.value.length;
            
            $row = this.rows[i-1].getRow();
            var totalCost = catItem.value.reduce((tot, item) => tot + (parseFloat(item.price) * parseFloat(item.count)), 0);
            var $totalRow = $(`<tr>
                <td colspan="${self.columnInfo.length - 3}"></td>
                <td colspan="2" style="background-color: ${catItem.background};">${catItem.name} Total</td>
                <td style="background-color: ${catItem.background};">${formatMoney(totalCost)}</td>
            </tr>`);
            $totalRow.insertAfter($row);
        });
    }

    insertTotalsSection() {
        var self = this;
        var $tbody = this.$table.find('tbody');
        var totalUnits = this.categorizedData.reduce((tot, catItem) => {
            return tot + catItem.value.reduce((t, item) => t + parseFloat(item.count), 0);
        }, 0);

        // add empty row
        $tbody.append($('<tr><td colspan="100"></td></tr>'));

        // add total units row
        $tbody.append($(`<tr>
            <td colspan="${self.columnInfo.length - 3}" style="visibility: hidden;"></td>
            <td colspan="2">Total Units</td>
            <td>${totalUnits}</td>
        </tr>`));

        this.categorizedData.forEach((catItem) => {
            var totalCost = catItem.value.reduce((tot, item) => tot + (parseFloat(item.price) * parseFloat(item.count)), 0);
            $tbody.append($(`<tr>
                <td colspan="${self.columnInfo.length - 3}" style="visibility: hidden;"></td>
                <td colspan="2" style="background-color: ${catItem.background};">${catItem.name} Total</td>
                <td style="background-color: ${catItem.background};">${formatMoney(totalCost)}</td>
            </tr>`));
        });

    }
}
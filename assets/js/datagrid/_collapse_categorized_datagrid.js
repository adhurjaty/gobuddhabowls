import { CollapsibleDatagrid } from "./_collapsible_datagrid";
import { groupByCategory } from "../helpers/_helpers";

export class CollapseCategorizedDatagrid extends CollapsibleDatagrid {
    constructor(data, columnInfo, collapseInfo, updateFn) {
        super(data, columnInfo, collapseInfo, updateFn);
        this.categorizedData = groupByCategory(data);
        this.insertCategoryLabels();
    }

    insertCategoryLabels() {
        var i = 0;
        this.categorizedData.forEach((catItem) => {
            var $row = this.rows[i].getRow();
            var $labelRow = $(`<tr class="category-header" style="background-color: ${catItem.background};">
                <td colspan="100">${catItem.name}</td>
            </tr>`);

            $labelRow.insertBefore($row);
            i += catItem.value.length;
        });
    }
}
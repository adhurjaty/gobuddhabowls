import { DataGrid } from "./_datagrid";

var collapsedCaret = 'fa-caret-right';
var expandedCaret = 'fa-caret-down';

export class CollapsibleDatagrid extends DataGrid {
    constructor(data, columnInfo, collapseInfo, updateFnc) {
        super(data, columnInfo, updateFnc);
        this.collapseInfo = collapseInfo;
        this.insertCollapers();
        this.insertCollapseRows();
        this.initCollapsers();
    }

    insertCollapers() {
        this.$table.find('thead').prepend('<th>&nbsp;</th>');
        this.rows.forEach((row) => {
            var $row = row.getRow();
            $row.prepend(`<td class="expander"><span class="fa ${collapsedCaret}"></span></td>`);
        });
    }

    insertCollapseRows() {
        for (let i = this.data.length-1; i >= 0; i--) {
            const $row = this.rows[i].getRow();
            const item = this.data[i];
            var $hiddenRow = $(this.collapseInfo(item))
            $hiddenRow.hide();
            $hiddenRow.insertAfter($row);
        }
    }

    initCollapsers() {
        this.$table.find('td.expander').click((event) => {
            var $td = $(event.currentTarget);
            var $span = $td.find('span');
            if($span.hasClass(collapsedCaret)) {
                $span.removeClass(collapsedCaret);
                $span.addClass(expandedCaret);
                $td.parent().next().show();
            } else {
                $span.removeClass(expandedCaret);
                $span.addClass(collapsedCaret);
                $td.parent().next().hide();
            }
        });
    }
}
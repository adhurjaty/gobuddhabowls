import { replaceUrlId } from "./helpers";

// require("bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js");

var collapsedCaret = 'fa-caret-right';
var expandedCaret = 'fa-caret-down';

export class EditItem {
    constructor(datagrid, $td) {
        this.datagrid = datagrid;
        this.$td = $td;
        this.type = $td.attr('data-type');
        this.id = $td.parent().attr('item-id');
        this.field = $td.attr('field');
        this.contents = $td.text();
        this.errorMessage = "";
        this.setListener();
    }

    setListener() {
        var self = this;

        switch(self.type) {
            case 'date':
                self.$td.on('focus', function(event) {
                    self.clearError($(this));
                    var $date = $('<input data-provide="datepicker" value="' + self.contents + '">');
                    $date.css('width', '85px');
                    $(this).empty();
                    $(this).append($date);
                    var startDate = self.contents ? self.contents : new Date().toLocaleDateString("en-US");
                    $date.datepicker({
                        autoclose: 'true',
                        format: 'mm/dd/yyyy',
                        defaultViewDate: startDate
                    }).on('changeDate', function(event) {
                        self.contents = event.format();
                        self.datagrid.sendUpdate(self);
                    }).on('hide', function(event) {
                        self.$td.text(self.contents);
                        if(self.errorMessage) {
                            self.showError(self.errorMessage);
                        }
                    });
                    $date.focus();
                });
                break;
            case 'money':
                break;
            case 'selector':
                break;
            default:    // type 'text'
                $(el).on('focus', function(event) {
                    self.clearError($(this));
                });
                $(el).on('blur', function(event) {
                    var id = $(this).parent().attr('item-id');
                    var field = $(this).attr('field');
                    self.datagrid.sendUpdate(self);
                });
                break;
        }
    }

    showError(msg) {
        this.errorMessage = msg;
        var $span = $('<span></span>');
        $span.text(msg);
        $span.css('color', 'red');
        this.$td.append($('<br>'));
        this.$td.append($span);
    }

    clearError() {
        this.$td.find('span').remove();
    }

    onUpdateSuccess() {
        switch(this.type) {
            case 'date':
                this.$td.text(this.contents);
                break;
            default:
                break;
        }
    }
}

export class DataGrid {
    constructor(grid) {
        this.grid = grid;
        this.on_change_url = $(grid).attr('onchange-href');
        this.initRows();

        if($(grid).find('td.expander') != undefined) {
            this.initCollapse();
        }
    }

    initRows() {
        self = this;
        $.each($(self.grid).find('tr'), function(i, tr) {
            $(tr).click(function(event) {
                if(!$(this).hasClass('active')) {
                    self.clearSelectedRow();
                    $(this).addClass('active');
                    self.setEditable($(this));
                }
            })
        });
    }

    initCollapse() {
        self = this;
        $.each($(this.grid).find('td.expander'), function(i, el) {
            $(el).click(function(event) {
                var $span = $(this).find('span');
                if($span.hasClass(collapsedCaret)) {
                    $span.removeClass(collapsedCaret);
                    $span.addClass(expandedCaret);
                    self.expandInfo($(this).parent().next());
                } else if($span.hasClass(expandedCaret)) {
                    $span.removeClass(expandedCaret);
                    $span.addClass(collapsedCaret);
                    self.collapseInfo($(this).parent().next());                    
                }
            });
        });
    }

    expandInfo($tr) {
        $tr.show();
    }

    collapseInfo($tr) {
        $tr.hide();
    }

    clearSelectedRow() {
        var row = $(this.grid).find('tr.active');
        row.removeClass('active');
        this.removeEditable(row);
    }

    setEditable(row) {
        row.find('td[editable="true"]').attr('contenteditable', true);
    }

    removeEditable(row) {
        row.find('td[editable="true"]').removeAttr('contenteditable');
    }

    sendUpdate(editItem) {
        if(this.on_change_url != undefined) {
            var data = {};
            data[editItem.field] = editItem.contents;
            $.ajax({
                url: replaceUrlId(this.on_change_url, editItem.id),
                data: data,
                method: "PUT",
                error: function(xhr, status, err) {
                    var errMessage = xhr.responseText;
                    editItem.showError(errMessage);
                },
                success: function(data, status, xhr) {
                    editItem.onUpdateSuccess();
                }
            });
        }
    }
}
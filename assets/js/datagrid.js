import { replaceUrlId } from "./helpers";

var collapsedCaret = 'fa-caret-right';
var expandedCaret = 'fa-caret-down';

// EditItem represents an editable cell in a Datagrid. 
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

    // setListener initializes the behaviors of the different types of cell data types
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
            // TODO: fill these options in
            case 'money':
                break;
            case 'selector':
                break;
            case 'number':
                self.$td.on('focus', function(event) {
                    self.clearError($(this));
                    document.execCommand('selectAll',false,null)                    
                });
                self.$td.on('blur', function(event) {
                    if(!isNaN($(this).text())) {
                        self.contents = $(this).text();
                        self.datagrid.sendUpdate(self);
                    } else {
                        $(this).text("0");
                    }
                });
                break;
            default:    // type 'text'
                self.$td.on('focus', function(event) {
                    self.clearError($(this));
                });
                self.$td.on('blur', function(event) {
                    // var id = $(this).parent().attr('item-id');
                    // var field = $(this).attr('field');
                    self.content = $(this).text();
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

// DataGrid is a class for creating a table that has editable cells that
// may update models on edit
export class DataGrid {
    constructor(grid) {
        this.grid = grid;
        this.on_change_url = $(grid).attr('onchange-href');
        this.initRows();

        if($(grid).find('td.expander') != undefined) {
            this.initCollapse();
        }
    }

    // initRows sets click highlighting for rows
    // TODO: remove highlighting when clicking off the table
    initRows() {
        self = this;
        $.each($(self.grid).find('tbody>tr'), function(i, tr) {
            $(tr).click(function(event) {
                if(!$(this).hasClass('active')) {
                    self.clearSelectedRow();
                    $(this).addClass('active');
                    self.setEditable($(this));
                }
            })
        });

        $(this.grid).on('focusout', function(event) {
            debugger;
            $(this).find('tr').each(function(i, tr) {
                $(tr).removeClass('active');
            });
        });
    }

    // initCollapse sets up tables with hidden rows to enable showing
    // row detail content. Add 'expander' class to the caret td tag
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

    // clearSelectedRow unhighlighs a row
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

    // sendUpdate updates the model if an 'on_change_url' atrribute is defined
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
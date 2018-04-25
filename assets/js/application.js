require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap-sass/assets/javascripts/bootstrap.js");
require("bootstrap-table/src/bootstrap-table.js");
require("bootstrap");
require("bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js");
window.Sortable = require("sortablejs");
require("bootstrap-colorpicker");
require("./helpers.js");
require("./period_selector.js");
require("./datagrid.js");

$(() => {
    var el = document.getElementById('categories-movable');
    if(el != undefined) {
        var sortable = Sortable.create(el, {
            group: {
                name: "components",
                pull: function(to, from, dragEl, evt) {
                  if(evt.type === 'dragstart') {
                    return false;
                  }
                  return true;
                }
            },
            animation: 150,
            handle: '.drag-handle'
        });
    }

    $('.colorpicker-component').colorpicker();
});
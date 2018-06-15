import "expose-loader?$!expose-loader?jQuery!jquery";
import "bootstrap-sass/assets/javascripts/bootstrap.js";
import "bootstrap-table";
// needed to comment this out so modals would not disappear
// right after appearing
// import "bootstrap";
import "bootstrap-datepicker";
window.Sortable = require("sortablejs");
import "bootstrap-colorpicker";
import "./helpers.js";
import "./period_selector.js";
import "./inventory_item_categories.js";
import "./purchase_orders_index.js";
import "./purchase_order_form.js";

$(() => {
    
});
$('#vendor-items-table').html('<%= partial("purchase_orders/inventory_items_datagrid.html") %>');
// When the user scrolls the page, execute myFunction 
window.onscroll = function() {

    // Get the header
    var header = document.getElementById("sticky-header");
    // Get the offset position of the navbar
    var sticky = header.offsetTop;
    debugger;

    // Add the sticky class to the header when you reach its scroll position. Remove "sticky" when you leave the scroll position
    if (window.pageYOffset >= sticky) {
        header.classList.add("sticky");
    } else {
        header.classList.remove("sticky");
    }
};

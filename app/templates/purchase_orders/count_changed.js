<%= if (categoryDetails.Total == 0.0) { %>
    var content = "";
<% } else { %>
    var content = '<%= partial("partials/horizontal_percentage_chart.html", {categoryDetails: categoryDetails, title: title}) %>';
<% } %>
$('#order-category-breakdown').html(content);
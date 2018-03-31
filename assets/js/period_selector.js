$('#period-selector').change(function() {

});
$('#week-selector').change(function() {
    
});
$('#year-selector').change(function() {
    
});

$.each($('.input-daterange input'), function(i, el) {
    $(this).datepicker({
        format: "mm/dd/yyyy",
        
    });
});
var endpoint = $('#period-selector-component').attr('onchange-href');

$('#period-selector-component').change(function() {
    $('#period-selector').change(function(event) {
        var start_date = $('option:selected', this).attr('start-date');

        var data = {}
        data['StartTime'] = start_date;
        data['Year'] = $('#year-selector').val();

        $.ajax({
            url: endpoint,
            data: data,
            method: "GET",
            error: function(xhr, status, err) {
                // var errMessage = xhr.responseText;
                // editItem.showError(errMessage);
            },
            success: function(data, status, xhr) {
                // $('#week-selector').val($('#week-selector option:first').val())
            }
        });   
    });
    $('#week-selector').change(function() {
        var start_date = $(this).attr('start-date');
        var end_date = $(this).attr('end-date');

        var data = {}
        data['StartTime'] = start_date;
        data['EndTime'] = end_date;
        data['Year'] = $('#year-selector').val();

        $.ajax({
            url: endpoint,
            data: data,
            method: "GET",
            error: function(xhr, status, err) {
                debugger;
            },
            success: function(data, status, xhr) {
                // editItem.onUpdateSuccess();
            }
        });
    });
    $('#year-selector').change(function() {
        var date = {};
        data['Year'] = $('#year-selector').val();
        $.ajax({
            url: endpoint,
            data: data,
            method: "GET",
            error: function(xhr, status, err) {
                debugger;
            },
            success: function(data, status, xhr) {
                // $('#period-selector').val($('#period-selector option:first').val())
            }
        });  
    });

    $.each($('.input-daterange input'), function(i, el) {
        $(this).datepicker({
            autoclose: true,
            format: "mm/dd/yyyy",
        }).on('datechange', function(event) {

        });
    });

});

$(() => {
    var data = {}
    var $selectedWeek = $('#week-selector option:selected');
    data['StartTime'] = $selectedWeek.attr('start-date');
    data['EndTime'] = $selectedWeek.attr('end-date');
    data['Year'] = $('#year-selector option:selected').val();
    $.ajax({
        url: endpoint,
        data: data,
        method: 'GET',
        error: function(xhr, status, err) {
            debugger;
        },
        success: function(data, status, xhr) {

        }
    });
})

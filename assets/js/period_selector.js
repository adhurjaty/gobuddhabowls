// import "bootstrap-datepicker";
import { datepicker } from "./datepicker";
import { dateStringToISO, getYearFromDateString } from './helpers.js';

var endpoint = $('#period-selector-component').attr('onchange-href');

$('#period-selector').change(function(event) {
    debugger;

    var start_date = $('option:selected', this).attr('start-date');

    var data = {}
    data['StartTime'] = start_date;
    data['Year'] = $('#year-selector').val();

    $.ajax({
        url: endpoint,
        data: data,
        method: "GET",
        error: function(xhr, status, err) {
            debugger;
        },
        success: function(data, status, xhr) {
            // $('#week-selector').val($('#week-selector option:first').val())
        }
    });   
});
$('#week-selector').change(function() {
    var start_date = $('option:selected', this).attr('start-date');

    var data = {}
    data['StartTime'] = start_date;
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
    var data = {};
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

datepicker($('.input-daterange'), {
    autoclose: true,
    format: "mm/dd/yyyy",
});
$.each($('.input-daterange'), function(i, d) {
    $(this).on('changeDate', function(event) {
        var startDate = dateStringToISO($('#daterange-start').val());
        var endDate = dateStringToISO($('#daterange-end').val());
        if($(this) == $('#daterange-start')) {
            startDate = event.date.toISOString();
        }
        if($(this) == $('#daterange-end')) {
            endDate = event.date.toISOString();
        }
        var data = {};
        data['Year'] = event.date.getFullYear();
        data['StartTime'] = startDate;
        data['EndTime'] = endDate;

        $.ajax({
            url: endpoint,
            data: data,
            method: "GET",
            error: function(xhr, status, err) {
                debugger;
            },
            success: function(data, status, xhr) {
                
            }
        });
    });
});

$(() => {
    var data = {}
    var $selectedWeek = $('#week-selector option:selected');
    if($selectedWeek.hasClass('empty-option')) {
        data['StartTime'] = dateStringToISO($('#daterange-start').val());
        data['EndTime'] = dateStringToISO($('#daterange-end').val());
        data['Year'] = getYearFromDateString($('#daterange-start').val());
    } else {
        data['StartTime'] = $selectedWeek.attr('start-date');
        data['Year'] = $('#year-selector option:selected').val();
    }
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

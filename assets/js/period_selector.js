import { datepicker, daterange } from './_datepicker';

$(() => {

    var $form = $('#period-selector-form');
    var $periodSelector = $form.find('select[name="Period"]')
    var $weekSelector = $form.find('select[name="Week"]');
    var $yearSelector = $form.find('select[name="Year"]');

    $periodSelector.change(function(event) {
        var date = getFormattedTime($periodSelector.val());
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(date);
        $form.submit();
    });
    $weekSelector.change(function() {
        var date = getFormattedTime($weekSelector.val());
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(date);
        $form.submit();
    });
    $yearSelector.change(function() {
        var year = parseInt($yearSelector.val());
        var date = new Date(year, 1).toISOString();
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(date);
        $form.submit();
    });

    var range = daterange($('.input-daterange input'), () => {
        var startTime = $form.find('input[name="StartTime"]').val();
        var endTime = $form.find('input[name="EndTime"]').val();
        startTime = new Date(startTime).toISOString();
        endTime = new Date(endTime).toISOString();

        clearForm($form);
        $form.find('input[name="StartTime"]').val(startTime);
        $form.find('input[name="EndTime"]').val(endTime);

        $form.submit();
        $('.input-daterange').remove();
    });
});

function clearForm($form) {
    $form.find('select[name="Week"]').val(null);
    $form.find('select[name="Period"]').val(null);
    $form.find('select[name="Year"]').val(null);
}

function clearFormAndEnd($form) {
    $form.find('input[name="EndTime"]').remove();
    clearForm($form);
}

function getFormattedTime(dateStr) {
    return (new Date(dateStr)).toISOString()
}

function setMinEndDateRange(input) {
    if ($(input).attr('name') == 'EndTime') {
        var minDate = new Date($('form input[name="StartTime"]').val());
        minDate.setDate(minDate.getDate() + 1)

        return {
            minDate: minDate
        };
    }
}

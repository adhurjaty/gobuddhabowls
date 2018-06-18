// import "bootstrap-datepicker";
import { datepicker } from "./datepicker";

$(() => {

    var $form = $('#period-selector-form');
    var $periodSelector = $form.find('select[name="Period"]')
    var $weekSelector = $form.find('select[name="Week"]');
    var $yearSelector = $form.find('select[name="Year"]');

    $periodSelector.change(function(event) {
        var unixDate = $periodSelector.val();
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(unixDate);
        $form.submit();
    });
    $weekSelector.change(function() {
        var unixDate = $weekSelector.val();
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(unixDate);
        $form.submit();
    });
    $yearSelector.change(function() {
        var year = parseInt($yearSelector.val());
        var date = new Date(year, 1)
        var unixDate = date.getTime() / 1000;
        clearFormAndEnd($form);
        $form.find('input[name="StartTime"]').val(unixDate);
        $form.submit();
    });

    datepicker($('.input-daterange'), {
        autoclose: true,
        format: "mm/dd/yyyy",
    });
    $.each($('.input-daterange'), function(i, d) {
        $(this).on('changeDate', function(event) {
            var startTime = $form.find('input[name="StartTime"]').val();
            var endTime = $form.find('input[name="StartTime"]').val();
            startTime = new Date(startTime).getTime() / 1000;
            endTime = new Date(endTime).getTime() / 1000;

            clearForm($form);
            $form.find('input[name="StartTime"]').val(startTime);
            $form.find('input[name="EndTime"]').val(endTime);

            $form.submit();

        });
    });
});

function clearForm($form) {
    $form.find('input[name="StartTime"]').val(null);
    $form.find('select[name="Week"]').remove();
    $form.find('select[name="Period"]').remove();
    $form.find('select[name="Year"]').remove();
}

function clearFormAndEnd($form) {
    $form.find('input[name="EndTime"]').remove();
    clearForm($form);
}

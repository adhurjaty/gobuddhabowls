import { datepicker } from '../_datepicker';

$(() => {
    var dateInput = $('#inventory-form').find('input[name="Date"]');
    datepicker(dateInput[0]);
});
import { parseModelJSON, formatSlashDate, formatMoney } from '../helpers/_helpers';
import { CategorizedItemsDisplay } from '../components/_categorized_items_display';
import { datepicker } from '../_datepicker';

$(() => {
    var dateInput = $('#inventory-form').find('input[name="Date"]');
    datepicker(dateInput[0]);
});
import Pikaday from 'pikaday';
import { formatSlashDate } from './helpers/_helpers';

export function datepicker(el, onDateSelected = (d) => {}) {
    disableAutocomplete(el);
    return new Pikaday({
        field: el,
        format: 'YYYY-MM-DD',
        onSelect: onDateSelected
    });
}

export function daterange($inputs, onDateSelected = function(d) {}) {
    if($inputs.length < 2) {
        return datepicker($inputs.eq(0), onDateSelected);
    }

    $inputs.each((i, el) => {
        disableAutocomplete(el);
    });

    var $startInput = $inputs.eq(0);
    var $endInput = $inputs.eq(1);

    var onSelectStart = (date) => {
        var endDate = new Date($endInput.val());
        if(date > endDate) {
            $endInput.val(formatSlashDate(date));
        }

        onDateSelected(date)
    };

    var start = new Pikaday({
        field: $startInput.get(0),
        format: 'YYYY-MM-DD',
        onSelect: onSelectStart
    });

    var startDate = new Date($startInput.val());


    var end = new Pikaday({
        field: $endInput.get(0),
        format: 'YYYY-MM-DD',
        onSelect: onDateSelected,
        minDate: startDate,
    });

    return [start, end];
}


function disableAutocomplete(el) {
    $(el).attr('autocomplete', 'off');
}

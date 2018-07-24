import Pikaday from 'pikaday';
import { formatSlashDate } from './helpers';

export function datepicker(el, onDateSelected) {
    return new Pikaday({
        field: el,
        format: 'MM/DD/YYYY',
        onSelect: onDateSelected
    });
}

export function daterange($inputs, onDateSelected) {
    if($inputs.length < 2) {
        return datepicker($inputs.eq(0), onDateSelected);
    }

    var $startInput = $inputs.eq(0);
    var $endInput = $inputs.eq(1);

    var onSelectStart = (date) => {
        var endDate = new Date($endInput.val());
        if(date > endDate) {
            $endInput.val(formatSlashDate(date));
        }

        onDateSelected()
    };

    var start = new Pikaday({
        field: $startInput.get(0),
        format: 'MM/DD/YYYY',
        onSelect: onSelectStart
    });

    var startDate = new Date($startInput.val());

    var end = new Pikaday({
        field: $endInput.get(0),
        format: 'MM/DD/YYYY',
        onSelect: onDateSelected,
        minDate: startDate,
    });

    return [start, end];
}

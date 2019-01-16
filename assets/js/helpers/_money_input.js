// code from https://stackoverflow.com/questions/14650932/set-value-to-currency-in-input-type-number
$.fn.currencyInput = function() {
    this.each(function() {
        var wrapper = $("<div class='currency-input' />");
        $(this).wrap(wrapper);
        $(this).before("<span class='currency-symbol'>$</span>");
        $(this).val(this.valueAsNumber.toFixed(2));
        $(this).change(function() {
            var min = parseFloat($(this).attr("min"));
            var max = parseFloat($(this).attr("max"));
            var value = this.valueAsNumber;
            if(value < min)
                value = min;
            else if(value > max)
                value = max;
            $(this).val(value.toFixed(2));
        });
    });
};

$(() => {
    $('input.currency').attr('type', 'number');
    $('input.currency').attr('step', .01);
    $('input.currency').currencyInput();
});
export function textSelector(defaultText, options, id, name) {
    return `<input type="text"
        value="${defaultText}"
        name="${name}"
        list="${id}-datalist-options"/>
    <datalist id="${id}-datalist-options">
        ${options.reduce((s, option) => {
            return `${s}\n<option value="${option}">
                ${option}
            </option>`;
        }, "", this)}
    </datalist>`;
}

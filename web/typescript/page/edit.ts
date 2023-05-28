export function checkMaxLength() {
    type input = HTMLTextAreaElement | HTMLInputElement;

    const validate = (el: input, max: number) => {
        el.setCustomValidity(el.value.length > max
            ? `Your input must be up to ${max} characters.`
            : '');
    }

    [...document.querySelectorAll('[maxlength]')].forEach((el: input) => {
        const max = parseInt(el.getAttribute('maxlength'), 10);
        el.removeAttribute('maxlength');
        validate(el, max);

        el.addEventListener('input', () => validate(el, max), { passive: true });
    });
}

export function checkMaxLength() {
    type input = HTMLTextAreaElement | HTMLInputElement;

    [...document.querySelectorAll('[maxlength]')].forEach((element: input) => {
        const maxlength = parseInt(element.getAttribute('maxlength'), 10);
        element.removeAttribute('maxlength');
        element.addEventListener('input', () => {
            element.setCustomValidity(element.value.length > maxlength ? `Your input must be up to ${maxlength} characters.` : '');
        }, {passive: true});
    });
}

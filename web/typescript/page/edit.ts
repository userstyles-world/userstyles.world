export function checkMaxLength() {
    [...document.querySelectorAll('textarea[maxlength]')].forEach((element) => {
        const maxlength = element.getAttribute('maxlength');
        element.removeAttribute('maxlength');
        element.addEventListener('input', () => {
            element.setCustomValidity(element.value.length > maxlength ? `Your input must be up to ${maxlength} characters.` : '');
        }, {passive: true});
    });
}

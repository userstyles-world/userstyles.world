export function checkMaxLength() {
    [...document.querySelectorAll('textarea[maxlength]')].forEach((element) => {
        const maxlength = element.getAttribute('maxlength');
        element.removeAttribute('maxlength');
        element.addEventListener('input', () => {
            element.setCustomValidity(maxlength >= element.value.length ? '' : 'Your imput must be shorter than ' + maxlength + ' characters.');
        }, {passive: true});
    });
}

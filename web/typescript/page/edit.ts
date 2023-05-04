export function checkMaxLength() {
    [...document.querySelectorAll('textarea[maxlength]')].forEach((element) => {
        const maxlength = element.getAttribute('maxlength');
        console.log(maxlength)
        element.removeAttribute('maxlength');
        element.addEventListener('input', () => {
            console.log(element.value.length)
            element.setCustomValidity(maxlength >= element.value.length ? '' : 'Your imput must be shorter than ' + maxlength + ' characters.');
        }, {passive: true});
    });
}

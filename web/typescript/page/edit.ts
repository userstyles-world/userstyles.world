export function editLimit(selectorQery, maxLen) {
    const inputElement = document.querySelector(selectorQery);
    console.log(inputElement);
    inputElement && inputElement.addEventListener('input', () => {
        var len = inputElement.value.length;
        if(len >= maxLen) {
            inputElement.setCustomValidity('toolong');
        }
        else {
            inputElement.setCustomValidity('');
        }
    });
}

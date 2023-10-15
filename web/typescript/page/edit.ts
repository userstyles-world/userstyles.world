export function checkMaxLength() {
    type input = HTMLTextAreaElement | HTMLInputElement;
    type error = HTMLParagraphElement;

    const validate = (el: input, max: number) => {
        const curr = el.value.length;
        if (curr > max) {
            el.setCustomValidity(`Your input must be up to ${max} characters.`);
            message(el, `Input is too long. Characters used: ${curr}/${max}`);
        } else {
            el.setCustomValidity('');
            message(el, '');
        }
    }

    const message = (el: input, msg: string) => {
        const e = document.querySelector(`.danger.${el.name}`) as error;
        if (e) e.innerText = msg;
    }

    [...document.querySelectorAll('[maxlength]')].forEach((el: input) => {
        const max = parseInt(el.getAttribute('maxlength'), 10);
        el.removeAttribute('maxlength');
        validate(el, max);

        el.addEventListener('input', () => validate(el, max), { passive: true });
    });
}

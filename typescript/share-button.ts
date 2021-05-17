export function ShareButton() {
    const i = document.getElementById('share') as HTMLInputElement;
    const shareButton = document.getElementById('btn-share') as HTMLButtonElement;
    shareButton.addEventListener('click', () => {
        i.select();
        document.execCommand('copy');
        i.blur();
        shareButton.classList.add('copied');
    });
}

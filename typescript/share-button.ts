export function ShareButton() {
    const parentElement = document.getElementById('share') as HTMLInputElement;
    const shareButton = document.getElementById('btn-share') as HTMLButtonElement;
    if (!shareButton) {
        return;
    }
    shareButton.addEventListener('click', () => {
        parentElement.select();
        document.execCommand('copy');
        parentElement.blur();
        shareButton.classList.add('copied');
    });
}

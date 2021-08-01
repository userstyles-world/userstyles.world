export function ShareButton() {
    const parentElement = document.getElementById('share') as HTMLSpanElement;
    const shareButton = document.getElementById('btn-share') as HTMLButtonElement;
    if (!shareButton) {
        return;
    }
    shareButton.addEventListener('click', () => {
        const selection = window.getSelection();
        const range = document.createRange();
        range.selectNodeContents(parentElement);
        selection.removeAllRanges();
        selection.addRange(range);
        // add to clipboard.
        document.execCommand('copy');

        shareButton.classList.add('copied');
    });
}

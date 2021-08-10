import {doDomOperationProxy} from 'utils/dom';

export const initViewStyle = () => {
    doDomOperationProxy(() => {
        shareButton();
        checkIfStyleInstalled();
    });
};

function shareButton() {
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

function checkIfStyleInstalled() {
    const onMessage = (ev: MessageEvent<any>) => {
        if (!ev.data || !ev.data.type) {
            return;
        }
        const {type, data} = ev.data;
        if ('usw-style-info-response' === type && 'installed' === data.requestType) {
            window.removeEventListener('message', onMessage);
            if (data.installed) {
                const installButton: HTMLAnchorElement = document.querySelector('#install');
                if (installButton) {
                    // Need to user innerHTML to preserve the icon.
                    installButton.innerHTML = installButton.innerHTML.replace('Install', 'Reinstall');
                }
            }
        }
    };
    addEventListener('message', onMessage);
    // example URL: /style/1/any-slug
    const styleID = location.pathname.split('/')[2];
    dispatchEvent(new MessageEvent('message', {
        data: {type: 'usw-style-info-request', requestType: 'installed', styleID},
        origin: 'https://userstyles.world'
    }));
}

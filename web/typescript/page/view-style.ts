import {doDomOperation} from 'utils/dom';

export const initViewStyle = () => doDomOperation(() => {
    shareButton();
    checkIfStyleInstalled();
});

function shareButton() {
    const urlValue = document.getElementById('share').textContent;
    const shareButton = document.getElementById('btn-share') as HTMLButtonElement;
    if (!shareButton) {
        return;
    }
    shareButton.removeAttribute("hidden");
    shareButton.addEventListener('click', () => {
        navigator.clipboard.writeText(urlValue).then(() => {
            shareButton.classList.add('copied');
        }, () => {
            shareButton.classList.add('copied-failed');
        });
    });
}

function checkIfStyleInstalled() {
    const onMessage = (ev: MessageEvent<any>) => {
        const {type, data} = ev.data;
        if (!data || !type) {
            return;
        }
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

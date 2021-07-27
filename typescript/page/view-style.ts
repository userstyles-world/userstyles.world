import {doDomOperationProxy} from 'utils/dom';

export const checkIfStyleInstalled = doDomOperationProxy(() => {
    const onMessage = (ev: MessageEvent<any>) => {
        if (!ev.data || !ev.data.type) {
            return;
        }
        const {type, data} = ev.data;
        if ('usw-style-info-response' === type && 'installed' === data.requestType) {
            window.removeEventListener('message', onMessage);
            if (data.installed) {
                const writeReviewButton: HTMLAnchorElement = document.querySelector('#write-review');
                if (writeReviewButton) {
                    writeReviewButton.style.display = '';
                }
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
});

import {storeNewSettings} from 'utils/storage';

export function checkRedirect(redirect: string) {
    if (redirect) {
        storeNewSettings({redirect: ''});
        window.location.href = window.location.origin + redirect;
    }
}

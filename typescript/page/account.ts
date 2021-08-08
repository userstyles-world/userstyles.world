import {storeNewSettings} from 'utils/storage';

export function checkRedirect(redirect: string) {
    if (redirect) {
        storeNewSettings({redirect: ''});
        if (window.location.href !== redirect) {
            window.location.href = redirect;
        }
    }
}

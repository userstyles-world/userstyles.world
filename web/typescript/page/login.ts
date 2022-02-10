import {storeNewSettings} from 'utils/storage';

export function saveRedirect() {
    const redirect = new URLSearchParams(location.search).get('r');
    if (redirect) {
        storeNewSettings({redirect});
    }
}

import {storeNewSettings} from 'utils/storage';

export function saveRedirect() {
    const params = location.search + location.hash;
    const redirect = new URLSearchParams(params).get('r');
    if (redirect) {
        storeNewSettings({redirect});
    }
}

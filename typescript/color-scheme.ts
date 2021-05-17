import {isMatchMediaChangeEventListenerSupported} from './utils/platform';
import type {UserSettings} from './utils/storage';

const setColorSchemeAttribute = (value: string) => {
    document.documentElement.setAttribute('data-color-scheme', value);
};

const darkScheme = matchMedia('(prefers-color-scheme: dark)');
const handleColorScheme = () => {
    if (darkScheme.matches) {
        setColorSchemeAttribute('dark');
    } else {
        setColorSchemeAttribute('light');
    }
};

export function InitalizeColorScheme(colorScheme: UserSettings['colorScheme']) {
    switch (colorScheme) {
        case 'follow-system': {
            // By default it should be light the site. So if said browser
            // doesn't have this media query it will matches to false and
            // set the site to a light color-scheme.
            handleColorScheme();
            // As it follows the system we should listen for any changes.
            if (isMatchMediaChangeEventListenerSupported) {
                darkScheme.addEventListener('change', handleColorScheme);
            } else {
                darkScheme.addListener(handleColorScheme);
            }
            break;
        }
        default:
            setColorSchemeAttribute(colorScheme);
            // Make sure to remove the event listener.
            if (isMatchMediaChangeEventListenerSupported) {
                darkScheme.removeEventListener('change', handleColorScheme);
            } else {
                darkScheme.removeListener(handleColorScheme);
            }
            break;
    }
}

import {isMatchMediaChangeEventListenerSupported} from './utils/platform';
import type {UserSettings} from './utils/storage';

const setColorSchemeAttribute = (value: string) => {
    document.documentElement.setAttribute('data-color-scheme', value);
};

// By default it should be dark the site. So if said browser
// doesn't have this media query it will matches to false and
// set the site to a dark color-scheme.
const lightScheme = matchMedia('(prefers-color-scheme: light)');
const handleColorScheme = () => {
    if (lightScheme.matches) {
        setColorSchemeAttribute('light');
    } else {
        setColorSchemeAttribute('dark');
    }
};

export function InitalizeColorScheme(colorScheme: UserSettings['colorScheme']) {
    switch (colorScheme) {
        case 'follow-system': {
            handleColorScheme();
            // As it follows the system we should listen for any changes.
            if (isMatchMediaChangeEventListenerSupported) {
                lightScheme.addEventListener('change', handleColorScheme);
            } else {
                lightScheme.addListener(handleColorScheme);
            }
            break;
        }
        default:
            setColorSchemeAttribute(colorScheme);
            // Make sure to remove the event listener.
            if (isMatchMediaChangeEventListenerSupported) {
                lightScheme.removeEventListener('change', handleColorScheme);
            } else {
                lightScheme.removeListener(handleColorScheme);
            }
            break;
    }
}

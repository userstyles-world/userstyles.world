import {changeEntriesBehavior} from 'page/modlog';
import {InitalizeColorScheme as initalizeColorScheme} from './color-scheme';
import {ShareButton} from './share-button';
import {BroadcastReady} from './third-party';
import {SaveUserSettingsButton, SetValues} from './user-settings';
import {doDomOperation} from './utils/dom';
import type {UserSettings} from './utils/storage';
import {getSettings} from './utils/storage';

const WhenDOMReady = () => {
    ShareButton();
    BroadcastReady();
    SaveUserSettingsButton(onSettingsUpdate);
    SetValues(getSettings());
};

// WhenDOMReady contains code that only should be handle
// when the DOM is ready to go.
// Any other code shouldn't depend on this setup function.
doDomOperation(WhenDOMReady);

// Once settings update we should reinstalize any functionallity.
// That relies on this settings.
const onSettingsUpdate = () => {
    const settings = getSettings();
    initalizeColorScheme(settings.colorScheme);
};

// Initalize functions that requires settings and don't depend on the DOM.
// Note that we don't save getSettings() result, as this initalize is a 1 time thing
// And having it sit in the memory is kinda useless.
initalizeColorScheme(getSettings().colorScheme);

function pageSpecificFunctions(settings: UserSettings) {
    switch (location.pathname) {
        case '/modlog':
            changeEntriesBehavior(settings.entriesBehavior);
            break;
    }
}

pageSpecificFunctions(getSettings());

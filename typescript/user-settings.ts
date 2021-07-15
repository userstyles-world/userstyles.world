import type {UserSettings} from './utils/storage';
import {storeNewSettings} from './utils/storage';

const PREFIX = 'usr-settings';

export function SetValues(settings: UserSettings) {
    if (!window.location.pathname.startsWith('/account')) {
        return;
    }
    (document.getElementById(`${PREFIX}--color-scheme`) as HTMLSelectElement).value = settings.colorScheme;
    (document.getElementById(`${PREFIX}--entries-behavior`) as HTMLSelectElement).value = settings.entriesBehavior;
}

export function SaveUserSettingsButton(onSettingsUpdate: () => void) {
    const saveButton = document.getElementById(`${PREFIX}--save`) as HTMLButtonElement;
    if (!saveButton) {
        return;
    }
    saveButton.addEventListener('click', () => {
        const newSettings: Partial<UserSettings> = {};

        newSettings.colorScheme =
            (document.getElementById(`${PREFIX}--color-scheme`) as HTMLSelectElement).value as any;
        newSettings.entriesBehavior =
            (document.getElementById(`${PREFIX}--entries-behavior`) as HTMLSelectElement).value as any;

        storeNewSettings(newSettings);
        onSettingsUpdate();
    }, {passive: true});
}

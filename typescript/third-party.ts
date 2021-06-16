import {removeElement} from './utils/dom';

const fillInformationOnForm = (key: string, value: string) => {
    if (!value) {
        return;
    }
    const element: HTMLInputElement = document.body.querySelector(`form [name="${key}"`);
    if (!element) {
        return;
    }
    element.value = value;
};

const onMessage = (ev: MessageEvent<any>) => {
    if (!ev.data || !ev.data.type) {
        return;
    }
    const {type, data} = ev.data;
    switch (type) {
        case 'usw-remove-stylus-button': {
            removeElement(document.querySelector('a#stylus'));
        }
        case 'usw-fill-new-style': {
            if ('/api/oauth/authorize_style/new' !== window.location.pathname || !data) {
                return;
            }
            // TODO figure out which fields we can use more.
            fillInformationOnForm('name', data['name']);
            fillInformationOnForm('description', data['description']);
            fillInformationOnForm('code', data['sourceCode']);

        }
    }
};

window.addEventListener('message', onMessage);

export function BroadcastReady() {
    window.dispatchEvent(new MessageEvent('message', {
        data: {type: 'usw-ready'},
        origin: 'https://userstyles.world'
    }));
}

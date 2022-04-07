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

const DEFAULT_USERSTYLE_META = (name: string) => `/* ==UserStyle==
@name           ${name || 'A new userstyle!'}
@namespace      userstyles.world
@version        1.0.0
==/UserStyle== */\n`;

const handleSourceCode = (sourceCode: string, name: string) => {
    if (!sourceCode) {
        return;
    }
    if (!/\/\* *?==UserStyle==/g.test(sourceCode)) {
        sourceCode = `${DEFAULT_USERSTYLE_META(name)}${sourceCode}`;
    }
    fillInformationOnForm('code', sourceCode);
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
            if ('/api/oauth/style/new' !== location.pathname || !data) {
                return;
            }
            fillInformationOnForm('name', data['name']);
            const metaData = data['metadata'];
            if (metaData) {
                fillInformationOnForm('description', metaData['description']);
                fillInformationOnForm('license', metaData['license']);
                fillInformationOnForm('homepage', metaData['homepage']);
            }
            handleSourceCode(data['sourceCode'], data['name']);
        }
    }
};

addEventListener('message', onMessage);

export function broadcastReady() {
    dispatchEvent(new MessageEvent('message', {
        data: {type: 'usw-ready'},
        origin: 'https://userstyles.world'
    }));
}

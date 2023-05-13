const fillInformationOnForm = (key: string, value: string) => {
    if (!value) {
        return;
    }
    const element: HTMLInputElement = document.querySelector(`form [name="${key}"`);
    if (!element) {
        return;
    }
    element.value = value;
};

const handleSourceCode = (sourceCode: string) => {
    if (!sourceCode) {
        return;
    }
    fillInformationOnForm('code', sourceCode);
};

const onMessage = (ev: MessageEvent<any>) => {
    const {type, data} = ev.data;
    if (!type) {
        return;
    }
    switch (type) {
        case 'usw-remove-stylus-button': {
            const el = document.querySelector('a#stylus');
            el && el.remove();
        }
        case 'usw-fill-new-style': {
            if (!data || '/api/oauth/style/new' !== location.pathname) {
                return;
            }
            fillInformationOnForm('name', data['name']);
            fillInformationOnForm('description', data['description']);
            const uc = data['usercssData'];
            if (uc) {
                fillInformationOnForm('homepage', uc['homepageURL']);
                fillInformationOnForm('license', uc['license']);
            }
            handleSourceCode(data['sourceCode']);
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

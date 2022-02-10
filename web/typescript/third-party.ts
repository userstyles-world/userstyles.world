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
    if (!data || !type) {
        return;
    }
    switch (type) {
        case 'usw-remove-stylus-button': {
            const el = document.querySelector('a#stylus');
            el && el.remove();
        }
        case 'usw-fill-new-style': {
            if ('/api/oauth/style/new' !== location.pathname) {
                return;
            }
            fillInformationOnForm('name', data['name']);
            const metaData = data['metadata'];
            if (metaData) {
                fillInformationOnForm('description', metaData['description']);
                fillInformationOnForm('license', metaData['license']);
                fillInformationOnForm('homepage', metaData['homepage']);
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

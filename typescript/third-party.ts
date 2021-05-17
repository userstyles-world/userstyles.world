const onMessage = (ev: MessageEvent<any>) => {
    if (!ev.data || !ev.data.type) {
        return;
    }
    switch (ev.data.type) {
        case "usw-remove-stylus-button": {
            document.querySelector('a#stylus').remove();
        }
    }
}

window.addEventListener("message", onMessage)

export function BroadcastReady() {
    window.dispatchEvent(new MessageEvent('message', {
        data: {type: "usw-ready"},
        origin: "https://userstyles.world"
    }));
}
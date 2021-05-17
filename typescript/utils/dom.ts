const readyStateListeners = new Set<() => void>();

export const isDOMReady = () => 'complete' === document.readyState || 'interactive' === document.readyState;
export const addDOMReadyListener = (listener: () => void) => readyStateListeners.add(listener);

if (!isDOMReady()) {
    const onReadyStateChange = () => {
        if (isDOMReady()) {
            document.removeEventListener('readystatechange', onReadyStateChange);
            readyStateListeners.forEach((listener) => listener());
            readyStateListeners.clear();
        }
    };
    document.addEventListener('readystatechange', onReadyStateChange);
}

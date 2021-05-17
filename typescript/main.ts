import {ShareButton} from './share-button';
import {BroadcastReady} from './third-party';
import {addDOMReadyListener, isDOMReady} from './utils/dom';

const setup = () => {
    ShareButton();
    BroadcastReady();
};

// Setup contains code that only should be handle
// when the DOM is ready to go.
// Any other code shouldn't depend on this setup function.
if (isDOMReady()) {
    setup();
} else {
    addDOMReadyListener(setup);
}

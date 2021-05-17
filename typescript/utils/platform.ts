export const isMatchMediaChangeEventListenerSupported = (
    'function' === typeof MediaQueryList &&
    'function' === typeof MediaQueryList.prototype.addEventListener
);

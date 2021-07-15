import {doDomOperation} from 'utils/dom';
import type {UserSettings} from 'utils/storage';

export function changeEntriesBehavior(behavior: UserSettings['entriesBehavior']) {
    doDomOperation(() => {
        if ('click' === behavior) {
            const explain = document.querySelector('#explaination');
            explain.textContent = 'You can click on censored entries to see them.';

            const onClick = (e: MouseEvent) => {
                const entry = e.target as HTMLElement;
                entry.classList.remove('spoiler');
                entry.style.cursor = '';
            };

            document.querySelectorAll('.spoiler').forEach((entry: HTMLElement) => {
                entry.classList.add('spoiler-click');
                entry.addEventListener('click', onClick);
                entry.style.cursor = 'pointer';
            });
        }
    });
}

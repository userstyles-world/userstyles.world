import {doDomOperation} from 'utils/dom';
import type {UserSettings} from 'utils/storage';

export function changeEntriesBehavior(behavior: UserSettings['entriesBehavior']) {
    doDomOperation(() => {
        const explaination = document.querySelector('#explaination');
        const spoilers = document.querySelectorAll('.spoiler');
        if ('click' === behavior) {
            explaination.textContent = 'You can click on censored entries to see them.';

            spoilers.forEach((spoiler: HTMLElement) => {
                spoiler.classList.add('spoiler-click');
                spoiler.style.cursor = 'pointer';

                spoiler.addEventListener('click', (e: MouseEvent) => {
                    const entry = e.target as HTMLElement;
                    entry.classList.remove('spoiler');
                    entry.style.cursor = '';
                });
            });
        }
        if ('no-hide' === behavior) {
            explaination.remove();
            spoilers.forEach((spoiler) => spoiler.classList.remove('spoiler'));
        }
    });
}

export function updateTimestamps() {
    const list = document.querySelectorAll('time[datetime]') as NodeListOf<HTMLTimeElement>;
    list.forEach(el => el.dateTime = new Date(el.dateTime).toLocaleString());
};

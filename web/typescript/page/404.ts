// Add a simple listener to the .prev-page element to go back into history.
export function page404() {
    const prevPageElement = document.querySelector('.prev-page');
    prevPageElement && prevPageElement.addEventListener('click', () => {
        history.back();
    });
}

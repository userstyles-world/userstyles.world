export function form_sort() {
    const sortSelector = document.querySelector('.submit-form');
    sortSelector && sortSelector.addEventListener('change', () => {
        sortSelector.form.submit()
    });
}

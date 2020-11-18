const Config = {
    selector: '.card img',
    offset: '25px',
}

if ("IntersectionObserver" in window) {
    const images = document.querySelectorAll(Config.selector)
    const observer = new IntersectionObserver(entries => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const image = entry.target
                image.src = image.dataset.src
                image.classList.remove('lazy')
                observer.unobserve(image)
            }
        })
    }, { rootMargin: Config.offset })
    images.forEach(image => {
        observer.observe(image)
    })
}

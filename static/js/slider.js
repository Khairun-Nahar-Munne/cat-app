let currentIndex = 0;

const slides = document.getElementById('slides');
const totalSlides = slides.children.length;

document.getElementById('next').onclick = () => {
    currentIndex = (currentIndex + 1) % totalSlides;
    updateSlider();
};

document.getElementById('prev').onclick = () => {
    currentIndex = (currentIndex - 1 + totalSlides) % totalSlides;
    updateSlider();
};

function updateSlider() {
    slides.style.transform = `translateX(-${currentIndex * 100}%)`;
}

function vote(action, imageId) {
    fetch("/", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `action=${action}&image_id=${imageId}`
    }).then(response => response.json())
      .then(data => alert(data.status));
}

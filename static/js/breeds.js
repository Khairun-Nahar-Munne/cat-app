document.addEventListener("DOMContentLoaded", function () {
    const images = document.querySelectorAll("#breedImagesSlider img");
    let currentIndex = 0;

    const showImage = (index) => {
        images.forEach((img, i) => {
            img.style.display = i === index ? "block" : "none";
        });
    };

    const nextImage = () => {
        currentIndex = (currentIndex + 1) % images.length;
        showImage(currentIndex);
    };

    const prevImage = () => {
        currentIndex = (currentIndex - 1 + images.length) % images.length;
        showImage(currentIndex);
    };

    showImage(currentIndex);

    // Add next and prev functionality if you want
    document.getElementById("nextBtn").addEventListener("click", nextImage);
    document.getElementById("prevBtn").addEventListener("click", prevImage);
});

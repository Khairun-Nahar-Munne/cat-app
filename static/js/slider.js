
    document.addEventListener('DOMContentLoaded', function() {
        const slider = document.getElementById('breedImagesSlider');
        const dots = document.querySelectorAll('#indicatorDots button');
        let currentIndex = 0;

        function showSlide(index) {
            const slides = slider.querySelectorAll('div');
            const totalSlides = slides.length;

            // Reset all dots
            dots.forEach(dot => dot.classList.remove('bg-gray-700'));
            dots.forEach(dot => dot.classList.add('bg-gray-500'));

            // Move the slider
            slider.style.transform = `translateX(-${index * 100}%)`;

            // Mark the active dot
            dots[index].classList.add('bg-gray-700');
        }

        // Add click listeners for the dots
        dots.forEach((dot, index) => {
            dot.addEventListener('click', function() {
                currentIndex = index;
                showSlide(currentIndex);
            });
        });

        // Set up automatic slide change every 3 seconds
        setInterval(function() {
            currentIndex = (currentIndex + 1) % dots.length; // Loop back to the first image after the last one
            showSlide(currentIndex);
        }, 3000); // Change slide every 3 seconds

        // Initial call to display the first slide
        showSlide(currentIndex);
    });


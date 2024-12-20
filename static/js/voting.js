document.addEventListener("DOMContentLoaded", function () {
    const likeBtn = document.getElementById("likeBtn");
    const dislikeBtn = document.getElementById("dislikeBtn");
    const heartBtn = document.getElementById("heartBtn");
    const catImage = document.getElementById("catImage");
    const loadingImage = "/static/img/download.png"; // Adjust path as needed

    // Function to fetch a new cat image
    const fetchNewImage = async () => {
        // Show loading image while waiting for new cat image
        const originalImage = catImage.src;
        catImage.src = loadingImage;

        try {
            const response = await fetch("/voting");
            if (response.ok) {
                const data = await response.json();
                catImage.src = data.imageURL; // Update to the new image
            } else {
                console.error("Failed to fetch new image:", response.status);
                catImage.src = originalImage; // Restore original image on failure
            }
        } catch (error) {
            console.error("Error fetching new image:", error);
            catImage.src = originalImage; // Restore original image on error
        }
    };

    heartBtn.addEventListener("click", fetchNewImage); // Just fetch new image
});

document.addEventListener("DOMContentLoaded", function () {
    const likeBtn = document.getElementById("likeBtn");
    const dislikeBtn = document.getElementById("dislikeBtn");
    const heartBtn = document.getElementById("heartBtn");
    const catImage = document.getElementById("catImage");
    const loadingImage = "/static/img/download.png"; // Adjust path as needed

    const fetchNewImage = async () => {
        // Show loading image
        const originalImage = catImage.src;
        catImage.src = loadingImage;

        try {
            const response = await fetch("/voting");
            if (response.ok) {
                const parser = new DOMParser();
                const doc = parser.parseFromString(await response.text(), "text/html");
                const newImageURL = doc.querySelector("#catImage").src;
                catImage.src = newImageURL; // Update to the new image
            } else {
                console.error("Failed to fetch new image:", response.status);
                catImage.src = originalImage; // Restore original image on failure
            }
        } catch (error) {
            console.error("Error fetching new image:", error);
            catImage.src = originalImage; // Restore original image on error
        }
    };

    likeBtn.addEventListener("click", fetchNewImage);
    dislikeBtn.addEventListener("click", fetchNewImage);
    heartBtn.addEventListener("click", fetchNewImage);
});

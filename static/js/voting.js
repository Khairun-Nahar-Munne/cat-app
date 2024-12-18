document.addEventListener("DOMContentLoaded", function () {
    const likeBtn = document.getElementById("likeBtn");
    const dislikeBtn = document.getElementById("dislikeBtn");
    const heartBtn= document.getElementById("heartBtn");
    const catImage = document.getElementById("catImage");

    const fetchNewImage = async () => {
        try {
            const response = await fetch("/voting");
            if (response.ok) {
                const parser = new DOMParser();
                const doc = parser.parseFromString(await response.text(), "text/html");
                const newImageURL = doc.querySelector("#catImage").src;
                catImage.src = newImageURL;
            }
        } catch (error) {
            console.error("Error fetching new image:", error);
        }
    };

    likeBtn.addEventListener("click", fetchNewImage);
    dislikeBtn.addEventListener("click", fetchNewImage);
    heartBtn.addEventListener("click", fetchNewImage);
});

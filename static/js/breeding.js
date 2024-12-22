document.addEventListener("DOMContentLoaded", function () {
    const breedSearch = document.getElementById("breedSearch");
    const breedDropdown = document.getElementById("breedDropdown");
    const clearSearch = document.getElementById("clearSearch");
    const breedName = document.getElementById("breedName");
    const breedOrigin = document.getElementById("breedOrigin");
    const breedId = document.getElementById("breedId");
    const breedDescription = document.getElementById("breedDescription");
    const breedWikipedia = document.getElementById("breedWikipedia");
    const sliderImages = document.getElementById("sliderImages");
    const sliderDots = document.getElementById("sliderDots");
    const breedInfo = document.getElementById("breedInfo");
    const breedsSection = document.getElementById("breedsSection");

    let sliderInterval;

    // Fetch all breeds
    fetch("/api/breeds")
        .then(response => response.json())
        .then(breeds => {
            const breedList = breeds.map(breed => `
                <div class="p-2 cursor-pointer hover:bg-gray-200" data-breed-id="${breed.id}">
                    ${breed.name}
                </div>`).join("");
            breedDropdown.innerHTML = breedList;

            // Automatically fetch and display the first breed's info
            if (breeds.length > 0) {
                const firstBreed = breeds[0];
                breedSearch.value = firstBreed.name;
                fetchBreedInfo(firstBreed.id);
                clearSearch.classList.remove("hidden");
            }
        });

    // Show dropdown on input click
    breedSearch.addEventListener("click", () => {
        breedDropdown.classList.remove("hidden");
    });

    // Hide dropdown when clicking outside
    document.addEventListener("click", (e) => {
        if (!breedDropdown.contains(e.target) && e.target !== breedSearch) {
            breedDropdown.classList.add("hidden");
        }
    });

    // Select breed from dropdown
    breedDropdown.addEventListener("click", (e) => {
        const breedID = e.target.dataset.breedId;
        const breedNameSelected = e.target.textContent;

        breedSearch.value = breedNameSelected;
        clearSearch.classList.remove("hidden");
        breedDropdown.classList.add("hidden");

        fetchBreedInfo(breedID);
    });

    // Clear search
    clearSearch.addEventListener("click", () => {
        breedSearch.value = "";
        clearSearch.classList.add("hidden");
        breedSearch.placeholder = "Please select a breed...";
    });

    // Fetch breed info and images
    function fetchBreedInfo(breedID) {
        fetch(`/api/breed/${breedID}`)
            .then(response => response.json())
            .then(data => {
                const breed = data.breed;
                const images = data.images;

                breedName.textContent = breed.name;
                breedOrigin.textContent = `${breed.origin}`;
                breedId.textContent = `${breed.id}`;
                breedDescription.textContent = breed.description;
                breedWikipedia.href = breed.wikipedia_url;

                setupSlider(images);
                breedInfo.classList.remove("hidden");
            });
    }

    // Setup slider
    function setupSlider(images) {
        clearInterval(sliderInterval);

        sliderImages.innerHTML = images.map((img, idx) => `
            <img 
                src="${img.url}" 
                class="w-full h-full object-cover ${idx === 0 ? "block" : "hidden"}" 
                data-index="${idx}" 
                alt="Breed Image">
        `).join("");

        sliderDots.innerHTML = images.map((_, idx) => `
            <div 
                class="w-2 h-2 rounded-full mx-1 cursor-pointer ${idx === 0 ? "bg-blue-500" : "bg-gray-300"}" 
                data-index="${idx}">
            </div>
        `).join("");

        let currentIndex = 0;
        const imageElements = sliderImages.querySelectorAll("img");
        const dotElements = sliderDots.querySelectorAll("div");

        function changeImage(index) {
            imageElements.forEach((img, idx) => {
                img.classList.toggle("hidden", idx !== index);
            });
            dotElements.forEach((dot, idx) => {
                dot.classList.toggle("bg-blue-500", idx === index);
                dot.classList.toggle("bg-gray-300", idx !== index);
            });
            currentIndex = index;
        }

        function startSlider() {
            sliderInterval = setInterval(() => {
                currentIndex = (currentIndex + 1) % images.length;
                changeImage(currentIndex);
            }, 1000); 
        }

        startSlider();

        dotElements.forEach(dot => {
            dot.addEventListener("click", () => {
                const index = parseInt(dot.dataset.index, 10);
                changeImage(index);
            });
        });
    }
});

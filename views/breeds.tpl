<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Breeds</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/style.css">
    <!-- Tailwind CSS CDN -->
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 text-center font-sans">

    <!-- Header -->
    <header class="bg-gray-200 py-4">
        <div class="flex justify-center space-x-8 font-bold text-lg text-gray-700">
            <a href="/voting" class="cursor-pointer hover:text-gray-500">Voting</a>
            <a href="/breeds" class="cursor-pointer hover:text-gray-500">Breeds</a>
            <a href="/favorites" class="cursor-pointer hover:text-gray-500">Favorites</a>
        </div>
    </header>

    <!-- Search Field -->
    <main class="mt-8">
        <div class="flex justify-center">
            <div class="relative w-[700px]">
                <input 
                    type="text" 
                    id="breedSearch" 
                    placeholder="{{.Breed.Name}}" 
                    class="px-4 py-2 rounded-md border w-full focus:outline-none focus:ring-2 focus:ring-blue-500" 
                    onfocus="showDropdown()" 
                    oninput="filterBreeds()"> <!-- Removed readonly here -->
                
                <!-- Clear icon to reset the input -->
                <button id="clearButton" class="absolute right-2 top-2 text-gray-500" onclick="clearBreedSearch()">
                    <i class="fa-solid fa-xmark"></i>
                </button>
                
                <ul 
                    id="breedDropdown" 
                    class="absolute left-0 w-full max-h-64 overflow-y-scroll bg-white border rounded-md hidden z-10">
                    {{range .Breeds}}
                    <li 
                        class="px-4 py-2 cursor-pointer hover:bg-gray-100" 
                        onclick="selectBreed('{{.ID}}', '{{.Name}}')">
                        {{.Name}}
                    </li>
                    {{end}}
                </ul>
            </div>
        </div>

        <!-- Breed Information -->
        {{if .Breed}}
      <section class="mt-4 w-[700px] mx-auto">
    <div class="mt-1 relative">
        <div id="breedImagesSlider" class="flex justify-center items-center space-x-4 overflow-hidden">
            {{range .BreedImages}}
            <div class="w-full h-96 hidden">
                <img src="{{.URL}}" alt="Cat Image" class="w-full h-full object-cover rounded-lg">
            </div>
            {{end}}
        </div>
    </div>

    <!-- Indicator Dots - moved below the image slider -->
    <div id="indicatorDots" class="flex justify-center pb-4 space-x-2 mt-2">
        {{range .BreedImages}}
        <button class="w-3 h-3 rounded-full bg-gray-500 hover:bg-gray-700 focus:outline-none"></button>
        {{end}}
    </div>

    <div class="text-left mt-2">
        <h2 class="text-2xl font-bold text-gray-800 inline">{{.Breed.Name}}</h2>
        <p class="text-gray-600 mt-2 inline ml-2">({{.Breed.Origin}})</p>
        <p class="text-gray-600 mt-2 inline ml-2">({{.Breed.ID}})</p>
        <p class="mt-6 text-gray-700">{{.Breed.Description}}</p>
        <p class="mt-4">
            <a href="{{.Breed.Wikipedia_URL}}" class="text-blue-500 hover:underline">Learn more on Wikipedia</a>
        </p>
    </div>
</section>

        {{end}}

    </main>

    <!-- Scripts -->
<script>
    let currentSlide = 0;

    // Save the original placeholder value
    const originalPlaceholder = document.getElementById('breedSearch').getAttribute('placeholder');
    
    // Flag to prevent placeholder overwrite after selecting a breed
    let isBreedSelected = false;

    function showDropdown() {
        document.getElementById('breedDropdown').classList.remove('hidden');
    }

    function filterBreeds() {
        const searchValue = document.getElementById('breedSearch').value.toLowerCase();
        const items = document.querySelectorAll('#breedDropdown li');
        items.forEach(item => {
            if (item.innerText.toLowerCase().includes(searchValue)) {
                item.classList.remove('hidden');
            } else {
                item.classList.add('hidden');
            }
        });
    }

    function selectBreed(breedId, breedName) {
        const breedSearchInput = document.getElementById('breedSearch');
        breedSearchInput.value = breedName;  // Set the breed name

        // Only reset placeholder if user manually cleared the input
        if (!isBreedSelected) {
            breedSearchInput.setAttribute('placeholder', originalPlaceholder); // Reset placeholder to original
        }
        
        isBreedSelected = true; // Mark that a breed was selected
        document.getElementById('breedDropdown').classList.add('hidden');
        document.getElementById('clearButton').classList.add('hidden');  // Show the clear button
        window.location.href = `/breeds?breed_id=${breedId}`;
    }

    function clearBreedSearch() {
        const breedSearchInput = document.getElementById('breedSearch');
        const clearButton = document.getElementById('clearButton');
        const breedDropdown = document.getElementById('breedDropdown');
        breedSearchInput.value = ''; // Clear the breed name

        // Remove the placeholder
        breedSearchInput.setAttribute('placeholder', 'Please Select');

        // Hide the clear button
        clearButton.classList.add('hidden');
        
        // Show the dropdown
        breedDropdown.classList.remove('hidden');

    }


    function showSlide(index) {
        const slides = document.querySelectorAll("#breedImagesSlider div");
        const dots = document.querySelectorAll("#indicatorDots button");

        // Hide all images and remove active dot
        slides.forEach(slide => slide.classList.add("hidden"));
        dots.forEach(dot => dot.classList.remove("bg-gray-700"));

        // Show the selected image and activate the corresponding dot
        slides[index].classList.remove("hidden");
        dots[index].classList.add("bg-gray-700");
        currentSlide = index;
    }

    function showNextSlide() {
        const slides = document.querySelectorAll("#breedImagesSlider div");
        currentSlide = (currentSlide + 1) % slides.length;
        showSlide(currentSlide);
    }

    // Initialize the first slide
    document.addEventListener("DOMContentLoaded", function() {
        const slides = document.querySelectorAll("#breedImagesSlider div");
        const dots = document.querySelectorAll("#indicatorDots button");

        if (slides.length > 0) {
            showSlide(0); // Show the first slide
            setInterval(showNextSlide, 2000); // Change slide every 2 seconds
        }

        // Add click event to dots
        dots.forEach((dot, index) => {
            dot.addEventListener("click", function() {
                showSlide(index); // Show corresponding image when clicking on a dot
            });
        });
    });

    document.addEventListener('click', (e) => {
        const dropdown = document.getElementById('breedDropdown');
        if (!document.getElementById('breedSearch').contains(e.target) && !dropdown.contains(e.target)) {
            dropdown.classList.add('hidden');
        }
    });
</script>

</body>
</html>

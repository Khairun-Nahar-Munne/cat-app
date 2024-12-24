<!-- views/index.tpl -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat App</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-white text-center  font-sans">
    <!-- Header -->
    <header class="border-b mb-4">
    <div class="max-w-screen-md mx-auto px-4">
        <div class="flex justify-center space-x-8 py-4">
        <button id="votingTab" class="flex flex-col items-center font-medium text-orange-500 hover:text-red-600 transition-colors duration-300 focus:outline-none"> 
            <span class="text-xl"><i class="fa-solid fa-up-down"></i></span>
            <span class="text-sm">Voting</span>
        </button>
        <button id="breedsTab" class="flex flex-col items-center font-medium text-gray-400 hover:text-red-500 transition-colors duration-300 focus:outline-none">
            <span class="text-xl"><i class="fa-solid fa-magnifying-glass"></i></span>
            <span class="text-sm">Breeds</span>
        </button>
        <button id="favouriteTab" class="flex flex-col items-center font-medium text-gray-400 hover:text-red-500 transition-colors duration-300 focus:outline-none">
            <span class="text-xl"><i class="fa-regular fa-heart"></i></span>
            <span class="text-sm">Favs</span>
        </button>
        </div>
    </div>
    </header>

    <!-- Voting Section -->
    <div id="votingSection" class="md:max-w-screen-md max-w-screen-sm  mx-auto px-4 py-8 md:px-6  md:border-2 rounded-xl md:py-6">
        <div class="success-message" id="successMessage"></div>

        <main>
            <div class="bg-white rounded-3xl shadow-lg overflow-hidden">
                <div id="loadingSpinnerImg" class="flex items-center justify-center w-full h-[500px] hidden">
                    <div class="animate-spin rounded-full h-20 w-20 border-b-4 border-orange-500">
                        <i class="fa-solid fa-cat text-6xl text-orange-500"></i>
                    </div>
                </div>
                <img id="catImage" src="{{.ImageURL}}" alt="Cat Image" class="w-full h-[500px] max-w-full object-cover">
            </div>
        </main>

        <input type="hidden" id="imageId" value="{{.ImageID}}">

        <footer class="mt-2">
            <div class="flex justify-between items-center px-4">
                <button id="fetchImageButton" class="text-3xl text-gray-400 hover:text-red-500">
                    <i class="fa-regular fa-heart"></i>
                </button>
                <div class="flex space-x-6">
                    <button id="likeBtn" class="text-3xl text-gray-400 hover:text-green-500">
                        <i class="fa-regular fa-thumbs-up"></i>
                    </button>
                    <button id="dislikeBtn" class="text-3xl text-gray-400 hover:text-red-500">
                        <i class="fa-regular fa-thumbs-down"></i>
                    </button>
                </div>
            </div>
        </footer>
    </div>

    <!-- Favourite Section -->
    <div id="favouriteSection" class="md:max-w-screen-md max-w-screen-sm  mx-auto px-4 py-8 md:px-6 md:border-2 rounded-xl md:py-6 hidden">
    <div class="flex justify-between items-center mb-6">
        <div class="flex gap-2">
        <button id="gridViewBtn" class="px-4 py-2 text-orange-500 text-2xl rounded hover:text-gray-600 transition'">
            <i class="fa-solid fa-table"></i>
        </button>
        <button id="listViewBtn" class="px-4 py-2 text-gray-500 text-2xl rounded hover:text-orange-600 transition">
            <i class="fa-solid fa-list"></i>
        </button>
        </div>
    </div>
    
    <div id="loadingSpinner" class="flex items-center justify-center w-full h-[500px] hidden">
        <div class="animate-spin rounded-full h-20 w-20 border-b-4 border-orange-500">
            <i class="fa-solid fa-cat text-6xl text-orange-500"></i>
        </div>
    </div>
    
    <div id="favoriteContent" class="min-h-[200px] max-h-[450px] overflow-y-auto"></div>
    </div>

    <!-- Breeds Section -->
    <div id="breedsSection" class="md:max-w-screen-md max-w-screen-sm  mx-auto px-4 py-8 md:px-6 md:border-2 rounded-xl md:py-6 hidden">
        <div class="search-section">
            <div class="relative">
                <input 
                    type="text" 
                    id="breedSearch" 
                    class="w-full p-3 rounded-lg border border-gray-200 focus:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:ring-opacity-30 transition-colors duration-200"
                    placeholder="Please select a breed..."
                >
                <button 
                    id="clearSearch" 
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 text-xl hidden hover:text-gray-700"
                >
                    &times;
                </button>
                <div 
                    id="breedDropdown" 
                    class="hidden absolute left-0 right-0 mt-1 max-h-40 overflow-y-auto border border-gray-200 bg-white rounded-lg shadow-lg z-10"
                ></div>
            </div>
        </div>
        <div id="loadingSpinner" class="flex items-center justify-center w-full h-[500px] hidden">
            <div class="animate-spin rounded-full h-20 w-20 border-b-4 border-orange-500">
                <i class="fa-solid fa-cat text-6xl text-orange-500"></i>
            </div>
        </div>
                <div id="breedInfo" class="mt-8 hidden">
            <div id="breedImageSlider" class="relative mb-6">
                <div id="sliderImages" class="w-full h-[300px] overflow-hidden rounded-lg">
                    <!-- Images will be inserted here -->
                    <img src="" alt="breed" class="w-full h-full object-cover">
                </div>
                <div id="sliderDots" class="flex justify-center gap-2 mt-3">
                    <!-- Dots will be dynamically added here -->
                    <button class="w-2 h-2 rounded-full bg-gray-300 hover:bg-gray-400 transition-colors duration-200"></button>
                    <button class="w-2 h-2 rounded-full bg-blue-500"></button>
                </div>
            </div>

            <div class="flex items-baseline gap-2 mb-2">
                <h3 id="breedName" class="text-xl sm:text-sm font-bold text-gray-800"></h3>
                <p id="breedOrigin" class="text-gray-600 leading-relaxed before:content-['('] after:content-[')']"></p>
                <p id="breedId" class="text-gray-600 leading-relaxed italic"></p>
            </div>
            <div class="text-left">
                <p id="breedDescription" class="text-gray-600 mb-2 leading-relaxed"></p>
                <a 
                    id="breedWikipedia" 
                    href="" 
                    target="_blank" 
                    class="text-sm uppercase text-orange-500 hover:text-orange-700 transition-colors"
                >
                    WIKIPEDIA
                </a>
            </div>
        </div>
    </div>

    <script src="/static/js/main.js"></script>
    <script src="/static/js/breeding.js"></script>
</body>
</html>
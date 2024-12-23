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
<body class="bg-white text-center font-sans">
    <!-- Header -->
    <header class="border-b">
    <div class="max-w-screen-md mx-auto px-4">
        <div class="flex justify-center space-x-8 py-4">
        <button id="votingTab" class="flex flex-col items-center font-medium text-red-500 hover:text-red-600 transition-colors duration-300 focus:outline-none"> 
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
    <div id="votingSection" class="max-w-screen-md mx-auto px-4 py-8">
        <div class="success-message" id="successMessage"></div>

        <main>
            <div class="bg-white rounded-3xl shadow-lg overflow-hidden">
                <div id="loadingSpinnerImg" class="text-center py-8 w-full h-[500px] max-w-full object-cover rounded-lg hidden">
                    <div class="animate-spin rounded-full h-20 w-20 border-b-4 border-blue-500 mx-auto">
                        <i class="fa-solid fa-cat text-6xl text-orange-500"></i>
                    </div>
                </div>
                <img id="catImage" src="{{.ImageURL}}" alt="Cat Image" class="w-full h-[500px] max-w-full object-cover">
            </div>
        </main>

        <input type="hidden" id="imageId" value="{{.ImageID}}">

        <footer class="mt-4">
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
    <div id="favouriteSection" class="max-w-screen-md mx-auto px-4 mt-8 hidden">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold">My Favorite Cats</h1>
        <div class="flex gap-2">
        <button id="gridViewBtn" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition">
            Grid View
        </button>
        <button id="listViewBtn" class="px-4 py-2 text-gray-400 rounded hover:text-gray-600 transition">
            <i class="fa-solid fa-list"></i>
        </button>
        </div>
    </div>
    <div id="loadingSpinner" class="text-center py-8 hidden">
        <div class="rounded-full h-20 w-20 border-blue-500 mx-auto flex items-center justify-center">
        <i class="fa-solid fa-cat text-6xl"></i>
        </div>
    </div>
    <div id="favoriteContent" class="min-h-[200px] max-h-[600px] overflow-y-auto"></div>
    </div>

    <!-- Breeds Section -->
    <div id="breedsSection" class="max-w-screen-md mx-auto px-4 mt-8 hidden">
        <div class="search-section">
            <div class="relative">
                <input type="text" id="breedSearch" class="p-2 rounded border w-full" placeholder="Please select a breed..." readonly>
                <button id="clearSearch" class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 hidden">
                    &times;
                </button>
            </div>
            <div id="breedDropdown" class="hidden mt-2 max-h-40 overflow-y-auto border bg-white"></div>
        </div>
        <div id="breedInfo" class="mt-4 hidden">
            <div id="breedImageSlider" class="mt-4 relative">
                <div id="sliderImages" class="w-full h-[300px] overflow-hidden rounded-lg"></div>
                <div id="sliderDots" class="flex justify-center mt-2"></div>
            </div>
            <h3 id="breedName" class="text-xl font-bold"></h3>
            <p id="breedOrigin"></p>
            <p id="breedId"></p>
            <p id="breedDescription"></p>
            <a id="breedWikipedia" href="" target="_blank" class="text-blue-500">Wikipedia Link</a>
        </div>
    </div>

    <script src="/static/js/main.js"></script>
    <script src="/static/js/breeding.js"></script>
</body>
</html>
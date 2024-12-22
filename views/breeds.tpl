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
    <!-- Add reference to separate JS file -->
    <script src="/static/js/breeds.js" defer></script>
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
                    class="px-4 py-2 rounded-md border w-full focus:outline-none focus:ring-2 focus:ring-blue-500">
                    
                <button id="clearButton" class="absolute right-2 top-2 text-gray-500">
                    <i class="fa-solid fa-xmark"></i>
                </button>
                                
                <ul id="breedDropdown" class="absolute left-0 w-full max-h-64 overflow-y-scroll bg-white border rounded-md hidden z-10">
                    {{range .Breeds}}
                    <li 
                        class="px-4 py-2 cursor-pointer hover:bg-gray-100"
                        data-breed-id="{{.ID}}" 
                        data-breed-name="{{.Name}}">
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

            <!-- Indicator Dots -->
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
    <script src="static/js/cat-vot.js"></script>
</body>
</html>
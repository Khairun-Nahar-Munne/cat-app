<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Voting</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" rel="stylesheet">
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

    <!-- Main Image Section -->
    <main class="mt-8 flex justify-center">
        <div class="flex justify-center items-center w-full max-w-screen-md px-4">
            <img id="catImage" src="{{.ImageURL}}" alt="Cat Image" class="w-full h-[500px] max-w-full object-cover rounded-lg">
        </div>
    </main>

    <!-- Footer -->
    <footer class="mt-8">
        <div class="w-full max-w-screen-md mx-auto px-8">
            <div class="flex justify-between items-center">
                <!-- Heart Button (Left Side) -->
                <div>
                    <button id="heartBtn" class="text-3xl text-black hover:text-red-500">
                        <i class="fa-regular fa-heart"></i>
                    </button>
                </div>

                <!-- Like and Dislike Buttons (Right Side) -->
                <div class="flex space-x-4">
                    <button id="likeBtn" class="text-3xl px-4 py-2 text-black rounded hover:text-green-600">
                        <i class="fa-regular fa-thumbs-up"></i>
                    </button>
                    <button id="dislikeBtn" class="text-3xl px-4 py-2 text-black rounded hover:text-red-600">
                        <i class="fa-regular fa-thumbs-down"></i>
                    </button>
                </div>
            </div>
        </div>
    </footer>

    <!-- Scripts -->
    <script src="/static/js/voting.js"></script>
</body>
</html>
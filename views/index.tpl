<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Slider</title>
    <script defer src="/static/js/slider.js"></script>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100">
    <div class="container mx-auto py-8">
        <div id="slider" class="relative w-full overflow-hidden">
            <div id="slides" class="flex transition-transform duration-300 ease-in-out">
                {{range .CatImages}}
                <div class="relative">
                    <img src="{{.URL}}" alt="Cat Image" class="w-full h-64 object-cover">
                    <button onclick="vote('like', '{{.ID}}')" class="absolute bottom-2 left-2 bg-green-500 text-white px-4 py-2">Like</button>
                    <button onclick="vote('dislike', '{{.ID}}')" class="absolute bottom-2 right-2 bg-red-500 text-white px-4 py-2">Dislike</button>
                </div>
                {{end}}
            </div>
            <button id="prev" class="absolute left-0 top-1/2 transform -translate-y-1/2 bg-gray-700 text-white px-4">Prev</button>
            <button id="next" class="absolute right-0 top-1/2 transform -translate-y-1/2 bg-gray-700 text-white px-4">Next</button>
        </div>
    </div>
</body>
</html>

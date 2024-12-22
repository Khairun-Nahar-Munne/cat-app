// static/js/main.js

// Global state
let currentSlide = 0;
let breeds = [];
let isBreedSelected = false;
let originalPlaceholder;

// Initialize everything when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeTabs();
    initializeVoting();
    const favouriteSection = document.getElementById('favouriteSection');
    const favoriteContent = document.getElementById('favoriteContent');
    const loadingSpinner = document.getElementById('loadingSpinner');
    const gridViewBtn = document.getElementById('gridViewBtn');
    const listViewBtn = document.getElementById('listViewBtn');
    const favouriteTab = document.getElementById('favouriteTab');
    let currentView = 'grid';
    let hasLoadedFavorites = false;

    // Initially hide the favorite content section
    if (favouriteSection) {
        favouriteSection.classList.add('hidden');
    }

    function showLoading() {
        loadingSpinner.classList.remove('hidden');
        favoriteContent.classList.add('hidden');
    }

    function hideLoading() {
        loadingSpinner.classList.add('hidden');
        favoriteContent.classList.remove('hidden');
    }

    async function getFavorites() {
        try {
            showLoading();
            const response = await fetch('/api/favourite');
            const data = await response.json();
            
            if (data.status === 'success') {
                displayFavorites(data.favorites);
                hasLoadedFavorites = true;
            } else {
                console.error('Error fetching favorites:', data.message);
                favoriteContent.innerHTML = `
                    <div class="text-center py-8 text-red-500">
                        Error loading favorites. Please try again.
                    </div>
                `;
            }
        } catch (error) {
            console.error('Error:', error);
            favoriteContent.innerHTML = `
                <div class="text-center py-8 text-red-500">
                    Error loading favorites. Please try again.
                </div>
            `;
        } finally {
            hideLoading();
        }
    }

    async function removeFavorite(favoriteId) {
        try {
            showLoading();
            const response = await fetch(`/api/favourite/${favoriteId}`, {
                method: 'DELETE'
            });
            
            const data = await response.json();
            if (data.status === 'success') {
                hasLoadedFavorites = false; // Reset the flag so we can fetch again
                getFavorites();
            } else {
                console.error('Error removing favorite:', data.message);
                alert('Error removing favorite. Please try again.');
                hideLoading();
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Error removing favorite. Please try again.');
            hideLoading();
        }
    }

    function displayFavorites(favorites) {
        if (!favorites || favorites.length === 0) {
            favoriteContent.innerHTML = `
                <div class="text-center py-8 text-gray-500">
                    No favorites found. Add some cats to your favorites!
                </div>
            `;
            return;
        }

        const container = document.createElement('div');
        
        if (currentView === 'grid') {
            container.className = 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6';
            
            favorites.forEach(favorite => {
                const card = document.createElement('div');
                card.className = 'bg-white rounded-lg shadow-lg overflow-hidden';
                card.innerHTML = `
                    <div class="relative aspect-w-16 aspect-h-9">
                        <img src="${favorite.image.url}" alt="Cat" 
                            class="w-full h-48 object-cover">
                    </div>
                    <div class="p-4">
                        <div class="flex justify-between items-center">
                            <span class="text-sm text-gray-500">
                                Added: ${new Date(favorite.created_at).toLocaleDateString()}
                            </span>
                            <button onclick="removeFavorite('${favorite.id}')" 
                                    class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition">
                                Remove
                            </button>
                        </div>
                    </div>
                `;
                container.appendChild(card);
            });
        } else {
            container.className = 'space-y-4';
            
            favorites.forEach(favorite => {
                const listItem = document.createElement('div');
                listItem.className = 'flex items-center bg-white rounded-lg shadow-lg p-4';
                listItem.innerHTML = `
                    <img src="${favorite.image.url}" alt="Cat" 
                        class="w-24 h-24 object-cover rounded-lg">
                    <div class="ml-4 flex-grow">
                        <span class="text-sm text-gray-500">
                            Added: ${new Date(favorite.created_at).toLocaleDateString()}
                        </span>
                    </div>
                    <button onclick="removeFavorite('${favorite.id}')" 
                            class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition">
                        Remove
                    </button>
                `;
                container.appendChild(listItem);
            });
        }

        favoriteContent.innerHTML = '';
        favoriteContent.appendChild(container);
    }

    // Add click event listener to the Favorites tab
    favouriteTab.addEventListener('click', () => {
        // Show the favorites section
        // Fetch favorites if not already loaded or need refresh
   
            getFavorites();
    });

    gridViewBtn.addEventListener('click', () => {
        currentView = 'grid';
        gridViewBtn.className = 'px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition';
        listViewBtn.className = 'px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition';
        if (hasLoadedFavorites) {
            getFavorites();
        }
    });

    listViewBtn.addEventListener('click', () => {
        currentView = 'list';
        gridViewBtn.className = 'px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition';
        listViewBtn.className = 'px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition';
        if (hasLoadedFavorites) {
            getFavorites();
        }
    });

    window.removeFavorite = removeFavorite;
});
// Tab Management
function initializeTabs() {
    const votingTab = document.getElementById('votingTab');
    const breedsTab = document.getElementById('breedsTab');
    const favouriteTab = document.getElementById('favouriteTab');
    const votingSection = document.getElementById('votingSection');
    const breedsSection = document.getElementById('breedsSection');
    const favouriteSection = document.getElementById('favouriteSection');

    votingTab.addEventListener('click', () => {
        votingSection.classList.remove('hidden');
        breedsSection.classList.add('hidden');
        favouriteSection.classList.add('hidden');
        votingTab.classList.add('active');
        breedsTab.classList.remove('active');
        favouriteTab.classList.remove('active');
    });

    breedsTab.addEventListener('click', () => {
        votingSection.classList.add('hidden');
        breedsSection.classList.remove('hidden');
        favouriteSection.classList.add('hidden');
        breedsTab.classList.add('active');
        votingTab.classList.remove('active');
        favouriteTab.classList.remove('active');
    });

    favouriteTab.addEventListener('click', () => {
        votingSection.classList.add('hidden');
        breedsSection.classList.add('hidden');
        favouriteSection.classList.remove('hidden'); // Show favouriteSection
        favouriteTab.classList.add('active');
        votingTab.classList.remove('active');
        breedsTab.classList.remove('active');
    });
}


// Voting Section Functionality
function initializeVoting() {
    const likeBtn = document.getElementById('likeBtn');
    const dislikeBtn = document.getElementById('dislikeBtn');
   
    likeBtn.addEventListener('click', () => {
        const imageId = document.getElementById('imageId').value;
        submitVote(1, imageId);
    });

    dislikeBtn.addEventListener('click', () => {
        const imageId = document.getElementById('imageId').value;
        submitVote(-1, imageId);
    });
}




function showError(message) {
    alert('Error: ' + message);
}

// Event listener for the heart button
document.getElementById('fetchImageButton').addEventListener('click', function() {
    const imageId = document.getElementById('imageId').value;
    submitFavoriteAndFetch(imageId);
});

function submitVote(value, imageId) {
    const buttons = document.querySelectorAll('button');
    const spinner = document.getElementById('loadingSpinnerImg');
    const catImage = document.getElementById('catImage');
    
    // Disable buttons and show spinner, hide current image
    buttons.forEach(button => button.disabled = true);
    spinner.classList.remove('hidden');
    catImage.classList.add('hidden');
    fetch('/api/vote', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: `value=${value}&image_id=${imageId}`
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            updateImage(data.image_url, data.image_id);
        } else {
            showError(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showError('An error occurred while submitting your vote');
    })
    .finally(() => {
        buttons.forEach(button => button.disabled = false);
    });
}

function submitFavoriteAndFetch(imageId) {
    const buttons = document.querySelectorAll('button');
    const spinner = document.getElementById('loadingSpinnerImg');
    const catImage = document.getElementById('catImage');
    
    // Disable buttons and show spinner, hide current image
    buttons.forEach(button => button.disabled = true);
    spinner.classList.remove('hidden');
    catImage.classList.add('hidden');

    // First submit the favorite
    fetch('/api/favorite', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: `image_id=${imageId}`
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            // Then fetch a new image
            return fetch('/api/cat/fetch');
        } else {
            throw new Error(data.message);
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            updateImage(data.image_url, data.image_id);
        } else {
            showError(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showError('An error occurred: ' + error.message);
    })
    .finally(() => {
        // Hide spinner, show image and re-enable buttons
        spinner.classList.add('hidden');
        catImage.classList.remove('hidden');
        buttons.forEach(button => button.disabled = false);
    });
}

// Make sure your updateImage function also handles the spinner
function updateImage(imageUrl, imageId) {
    const catImage = document.getElementById('catImage');
    const spinner = document.getElementById('loadingSpinnerImg');
    
    // Create a new image object to preload the image
    const newImage = new Image();
   
    newImage.onload = function() {
        // Once the new image is loaded, update the src and show it
        catImage.src = imageUrl;
        document.getElementById('imageId').value = imageId;
        
        // Show the image and hide the spinner
        catImage.classList.remove('hidden');
        spinner.classList.add('hidden');
    };
    
    // Start loading the new image
    newImage.src = imageUrl;
}



// Breeds Section Functionality


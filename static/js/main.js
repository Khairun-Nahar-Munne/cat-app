// static/js/main.js

// Global state
let currentSlide = 0;
let breeds = [];
let isBreedSelected = false;
let originalPlaceholder;
const tabs = document.querySelectorAll('button[id$="Tab"]');

// Add click event listener to each tab

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
        container.classList.add('max-w-screen-lg', 'mx-auto', 'px-1');  // Ensures consistent max width and padding
    
        if (currentView === 'grid') {
            container.classList.add('grid', 'grid-cols-2', 'md:grid-cols-2', 'lg:grid-cols-3', 'gap-6');
            
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
            container.classList.add('space-y-4', 'w-full', 'md:w-[700px]');
            
            favorites.forEach(favorite => {
                const listItem = document.createElement('div');
                listItem.className = 'flex items-center justify-between bg-white rounded-lg shadow-lg';
                listItem.innerHTML = `
                    <img src="${favorite.image.url}" alt="Cat" 
                        class="w-28 h-24 object-cover rounded-lg">
                    <button onclick="removeFavorite('${favorite.id}')" 
                            class="px-3 ml-4 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition">
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
        gridViewBtn.className = 'px-4 py-2 text-orange-500 text-2xl rounded hover:text-gray-600 transition';
        listViewBtn.className = 'px-4 py-2 text-gray-500 text-2xl rounded hover:text-orange-600 transition';
        if (hasLoadedFavorites) {
            getFavorites();
        }
    });

    listViewBtn.addEventListener('click', () => {
        currentView = 'list';
        gridViewBtn.className = 'px-4 py-2 text-gray-500 text-2xl rounded hover:text-orange-600 transition';
        listViewBtn.className = 'px-4 py-2 text-orange-500 text-2xl rounded hover:text-gray-600 transition';
        if (hasLoadedFavorites) {
            getFavorites();
        }
    });

    window.removeFavorite = removeFavorite;
});

document.getElementById('fetchImageButton').addEventListener('click', function () {
    const imageId = document.getElementById('imageId').value;
    submitFavoriteAndFetch(imageId);
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
            updateImage(data.image_url, data.image_id, buttons, spinner, catImage);
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
tabs.forEach(tab => {
    tab.addEventListener('click', () => handleTabClick(tab));
  });
  
  function handleTabClick(clickedTab) {
    // Remove active state from all tabs
    tabs.forEach(tab => {
      tab.classList.remove('text-orange-500');
      tab.classList.add('text-gray-400');
    });
    
    // Add active state to clicked tab
    clickedTab.classList.remove('text-gray-400');
    clickedTab.classList.add('text-orange-500');
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
                updateImage(data.image_url, data.image_id, buttons, spinner, catImage);
            } else {
                showError(data.message);
                buttons.forEach(button => button.disabled = false);
                spinner.classList.add('hidden');
                catImage.classList.remove('hidden');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showError('An error occurred: ' + error.message);
            buttons.forEach(button => button.disabled = false);
            spinner.classList.add('hidden');
            catImage.classList.remove('hidden');
        });
}

function updateImage(imageUrl, imageId, buttons, spinner, catImage) {
    // Create a new image object to preload the image
    const newImage = new Image();
    newImage.onload = function () {
        // Once the new image is loaded, update the src and show it
        catImage.src = imageUrl;
        document.getElementById('imageId').value = imageId;

        // Show the image and hide the spinner
        catImage.classList.remove('hidden');
        spinner.classList.add('hidden');
        buttons.forEach(button => button.disabled = false);
    };

    // Handle image loading errors
    newImage.onerror = function () {
        showError('Failed to load the image.');
        buttons.forEach(button => button.disabled = false);
        spinner.classList.add('hidden');
        catImage.classList.remove('hidden');
    };

    // Start loading the new image
    newImage.src = imageUrl;
}

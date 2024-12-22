// State management
let currentSlide = 0;
let isBreedSelected = false;
let originalPlaceholder;

// Initialize on DOM content loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeElements();
    initializeEventListeners();
    initializeSlider();
    initializeClickOutside();
});

// Initialize DOM elements
function initializeElements() {
    const breedSearchInput = document.getElementById('breedSearch');
    
    if (breedSearchInput) {
        originalPlaceholder = breedSearchInput.getAttribute('placeholder');
    }
}

// Initialize event listeners
function initializeEventListeners() {
    const breedSearchInput = document.getElementById('breedSearch');
    const breedDropdown = document.getElementById('breedDropdown');
    const clearButton = document.getElementById('clearButton');

    if (breedSearchInput) {
        breedSearchInput.addEventListener('input', filterBreeds);
        breedSearchInput.addEventListener('focus', showDropdown);
    }

    if (breedDropdown) {
        const breedItems = breedDropdown.querySelectorAll('li');
        breedItems.forEach(item => {
            item.addEventListener('click', function() {
                const breedId = this.getAttribute('data-breed-id');
                const breedName = this.getAttribute('data-breed-name');
                selectBreed(breedId, breedName);
            });
        });
    }

    if (clearButton) {
        clearButton.addEventListener('click', clearBreedSearch);
    }
}

// Breed search functionality
function showDropdown() {
    const dropdown = document.getElementById('breedDropdown');
    if (dropdown) {
        dropdown.classList.remove('hidden');
    }
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
    if (breedSearchInput) {
        breedSearchInput.value = breedName;

        if (!isBreedSelected) {
            breedSearchInput.setAttribute('placeholder', originalPlaceholder);
        }
    }
    
    isBreedSelected = true;

    const dropdown = document.getElementById('breedDropdown');
    const clearButton = document.getElementById('clearButton');
    
    if (dropdown) {
        dropdown.classList.add('hidden');
    }
    
    if (clearButton) {
        clearButton.classList.add('hidden');
    }

    window.location.href = `/breeds?breed_id=${breedId}`;
}

function clearBreedSearch() {
    const breedSearchInput = document.getElementById('breedSearch');
    const clearButton = document.getElementById('clearButton');
    const breedDropdown = document.getElementById('breedDropdown');
    
    if (breedSearchInput) {
        breedSearchInput.value = '';
        breedSearchInput.setAttribute('placeholder', 'Please Select');
    }
    
    if (clearButton) {
        clearButton.classList.add('hidden');
    }
    
    if (breedDropdown) {
        breedDropdown.classList.remove('hidden');
    }
}

// Image slider functionality
function showSlide(index) {
    const slides = document.querySelectorAll("#breedImagesSlider div");
    const dots = document.querySelectorAll("#indicatorDots button");

    slides.forEach(slide => slide.classList.add("hidden"));
    dots.forEach(dot => dot.classList.remove("bg-gray-700"));

    if (slides[index] && dots[index]) {
        slides[index].classList.remove("hidden");
        dots[index].classList.add("bg-gray-700");
    }
    currentSlide = index;
}

function showNextSlide() {
    const slides = document.querySelectorAll("#breedImagesSlider div");
    if (slides.length > 0) {
        currentSlide = (currentSlide + 1) % slides.length;
        showSlide(currentSlide);
    }
}

function initializeSlider() {
    const slides = document.querySelectorAll("#breedImagesSlider div");
    const dots = document.querySelectorAll("#indicatorDots button");

    if (slides.length > 0) {
        showSlide(0);
        setInterval(showNextSlide, 2000);

        // Add click events to dots
        dots.forEach((dot, index) => {
            dot.addEventListener("click", () => showSlide(index));
        });
    }
}

// Click outside handler
function initializeClickOutside() {
    document.addEventListener('click', (e) => {
        const dropdown = document.getElementById('breedDropdown');
        const breedSearch = document.getElementById('breedSearch');
        
        if (dropdown && breedSearch && 
            !breedSearch.contains(e.target) && 
            !dropdown.contains(e.target)) {
            dropdown.classList.add('hidden');
        }
    });
}
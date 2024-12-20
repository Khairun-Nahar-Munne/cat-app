document.addEventListener('DOMContentLoaded', function() {
    // Get all button elements
    const heartBtn = document.getElementById('heartBtn');
    const likeBtn = document.getElementById('likeBtn');
    const dislikeBtn = document.getElementById('dislikeBtn');
    
    // Add event listeners to buttons
    heartBtn.addEventListener('click', function() {
        const imageId = document.getElementById('imageId').value;
        handleHeartClick(imageId);
    });

    likeBtn.addEventListener('click', function() {
        const imageId = document.getElementById('imageId').value;
        submitVote(1, imageId);
    });

    dislikeBtn.addEventListener('click', function() {
        const imageId = document.getElementById('imageId').value;
        submitVote(-1, imageId);
    });
});

function handleHeartClick(imageId) {
    const buttons = document.querySelectorAll('button');
    buttons.forEach(button => button.disabled = true);

    fetch('/voting', {
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
            updateImage(data.image_url, data.image_id);
        } else {
            alert('Error: ' + data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred while loading a new image');
    })
    .finally(() => {
        buttons.forEach(button => button.disabled = false);
    });
}

function submitVote(value, imageId) {
    const buttons = document.querySelectorAll('button');
    buttons.forEach(button => button.disabled = true);

    fetch('/voting', {
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
            showSuccessMessage('Vote submitted successfully!');
        } else {
            alert('Error: ' + data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred while submitting your vote');
    })
    .finally(() => {
        buttons.forEach(button => button.disabled = false);
    });
}

function updateImage(imageUrl, imageId) {
    const catImage = document.getElementById('catImage');
    const imageIdInput = document.getElementById('imageId');

    const newImage = new Image();
    newImage.onload = function() {
        catImage.src = imageUrl;
        imageIdInput.value = imageId;
    };
    newImage.src = imageUrl;
}

function showSuccessMessage(message) {
    const successMessage = document.getElementById('successMessage');
    successMessage.textContent = message;
    successMessage.style.display = 'block';

    setTimeout(() => {
        successMessage.style.display = 'none';
    }, 1000);
}
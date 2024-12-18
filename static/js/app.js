document.addEventListener('DOMContentLoaded', function() {
    var voteButtons = document.querySelectorAll('.btn-vote');
    var favButtons = document.querySelectorAll('.btn-fav');

    voteButtons.forEach(function(button) {
        button.addEventListener('click', function() {
            // Add voting logic here
            console.log('Voting for cat image');
        });
    });

    favButtons.forEach(function(button) {
        button.addEventListener('click', function() {
            // Add favorite logic here
            console.log('Adding cat image to favorites');
        });
    });
});
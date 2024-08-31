var reloadModal = document.getElementById("reloadModal");
var reloadBtn = document.getElementById("reloadBtn");
var reloadModalClose = document.getElementById("reloadModalClose");

// When the user clicks the button, open the modal 
reloadBtn.onclick = function() {
    reloadModal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
reloadModalClose.onclick = function() {
    reloadModal.style.display = "none";
}

// Close the modal if user clicks outside of it
window.onclick = function(event) {
    if (event.target == reloadModal) {
        reloadModal.style.display = "none";
    }
}

// Handle form submit
document.getElementById('reloadForm').addEventListener('submit', function(e) {
    e.preventDefault();

    var formData = {
        ip: document.getElementById('reloadIP').value,
    };

    fetch('/reload', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        reloadModal.style.display = "none";
        fetchData();
    })
});
var addTargetModal = document.getElementById("addTargetModal");
var addNewTargetBtn = document.getElementById("addNewTargetBtn");
var addTargetModalClose = document.getElementById("addTargetModalClose");

// When the user clicks the button, open the modal 
addNewTargetBtn.onclick = function() {
    addTargetModal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
addTargetModalClose.onclick = function() {
    addTargetModal.style.display = "none";
}

// Close the modal if user clicks outside of it
window.onclick = function(event) {
    if (event.target == addTargetModal) {
        addTargetModal.style.display = "none";
    }
}

// Handle form submit
document.getElementById('addTargetForm').addEventListener('submit', function(e) {
    e.preventDefault();

    var formData = {
        jobName: document.getElementById('jobName').value,
        scheme: document.getElementById('scheme').value,
        metricsPath: document.getElementById('metricsPath').value,
        scrapeInterval: document.getElementById('scrapeInterval').value,
        scrapeTimeout: document.getElementById('scrapeTimeout').value,
        staticConfigs: [{
            targets: document.getElementById('targets').value.split(',')
        }],
        basicAuth: {
            username: document.getElementById('username').value,
            password: document.getElementById('password').value
        }
    };

    fetch('/newtarget', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        addTargetModal.style.display = "none";
        fetchData();
    })
});
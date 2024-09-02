var targetModal = document.getElementById("targetModal");
var targetBtn = document.getElementById("targetBtn");
var targetModalClose = document.getElementById("targetModalClose");

// When the user clicks the button, open the modal 
targetBtn.onclick = function() {
    console.log('Opening Add New Target Modal');
    targetModal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
targetModalClose.onclick = function() {
    targetModal.style.display = "none";
}

// Close the modal if user clicks outside of it
window.onclick = function(event) {
    if (event.target == targetModal) {
        targetModal.style.display = "none";
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
        }]
    };

    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;

    if (username && password) {
        formData.basicAuth = {
            username: username,
            password: password
        };
    }

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
        targetModal.style.display = "none";
        fetchData();
    })
});
var modal = document.getElementById("targetModal");

// Get the button that opens the modal
var btn = document.getElementById("addNewTarget");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks the button, open the modal 
btn.onclick = function() {
    modal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
span.onclick = function() {
    modal.style.display = "none";
}

// Close the modal if user clicks outside of it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

// Handle form submit
document.getElementById('newTargetForm').addEventListener('submit', function(e) {
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
        modal.style.display = "none";
        fetchData();
    })
});
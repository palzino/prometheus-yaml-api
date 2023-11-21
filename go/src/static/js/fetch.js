document.getElementById('getData').addEventListener('click', fetchData);
        document.getElementById('postData').addEventListener('click', saveData);
        
        function fetchData() {
            fetch('/promconf')
            .then(response => response.json())
            .then(data => {
                console.log('Data received:', data); // Check the received data
                if (data.ScrapeConfigs && Array.isArray(data.ScrapeConfigs)) {
                    
                    displayData(data.ScrapeConfigs);
                } else {
                    console.error('ScrapeConfigs is not an array');
                }
            })
            .catch(error => console.error('Error:', error));
        }

        function displayData(scrapeConfigs) {
            const container = document.getElementById('dataContainer');
            container.innerHTML = ''; // Clear previous data
            console.log('Received data:', scrapeConfigs);
            scrapeConfigs.forEach(config => {
                console.log(`Config ${config}`);
                const card = document.createElement('div');
                card.className = 'card';
                card.innerHTML = `
                    <h3>${config.jobName}</h3>
                    <p>Scheme: ${config.scheme || 'http'}</p>
                    <p>Metrics Path: ${config.metricsPath || 'Default'}</p>
                    <p>Scrape Interval: ${config.scrapeInterval || 'Default'}</p>
                    <p>Scrape Timeout: ${config.scrapeTimeout || 'Default'}</p>
                    ${formatStaticConfigs(config.staticConfigs)}
                    ${formatBasicAuth(config.basicAuth)}
                `;
                container.appendChild(card);
            });
    }

        function formatStaticConfigs(staticConfigs) {
            if (!staticConfigs || !staticConfigs.length) return '<p>Static Configs: None</p>';
            let targets = staticConfigs.map(config => config.targets.join(', ')).join('; ');
            return `<p>Static Configs Targets: ${targets}</p>`;
        }

        function formatBasicAuth(basicAuth) {
            if (!basicAuth) return '<p>Basic Auth: N/A</p>';
            return `<p>Basic Auth: Username - ${basicAuth.username}, Password - [protected]</p>`;
        }
        
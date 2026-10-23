function fetchFiles() {
    fetch('/last-uploaded')
        .then(response => response.json())
        .then(data => {
            const list = document.getElementById('last-uploaded-files');
            list.innerHTML = '';
            (data || []).forEach(file => {
                const listItem = document.createElement('li');
                listItem.textContent = file.name;
                list.appendChild(listItem);
            });
        })
        .catch(error => console.error('Error fetching files:', error));
}

document.getElementById('uploadForm').addEventListener('submit', function() {
    document.getElementById('loading').style.display = 'block';
});

window.onload = fetchFiles;
setInterval(fetchFiles, 5000); // Auto-update every 5 seconds

function fetchFiles() {
    fetch('/last-uploaded')
        .then(response => response.json())
        .then(data => {
            const list = document.getElementById('last-uploaded-files');
            list.innerHTML = '';
            data.forEach(file => {
                const listItem = document.createElement('li');
                listItem.innerHTML = `${file.name} `;
                list.appendChild(listItem);
            });
        });

    fetch('/server-files')
        .then(response => response.json())
        .then(data => {
            const list = document.getElementById('server-files');
            list.innerHTML = '';
            data.forEach(file => {
                const listItem = document.createElement('li');
                listItem.innerHTML = `${file.name} <button onclick="deleteFile('${file.name}')">Delete</button>`;
                list.appendChild(listItem);
            });
        });
}

function deleteFile(fileName) {
    fetch(`/delete?file=${fileName}`, { method: 'DELETE' })
        .then(response => {
            if (response.ok) {
                fetchFiles();
            } else {
                alert('Failed to delete file');
            }
        });
}


document.getElementById('uploadForm').addEventListener('submit', function() {
    document.getElementById('loading').style.display = 'block';
});

window.onload = fetchFiles;
setInterval(fetchFiles, 5000); // Auto-update every 5 seconds
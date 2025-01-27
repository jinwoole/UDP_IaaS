// iso-manager.js
class ISOManager {
    static async load() {
        try {
            const response = await fetch('/api/isos');
            const isos = await response.json();
            const tbody = document.querySelector('#isoTable tbody');
            const select = document.getElementById('vmIso');
            
            tbody.innerHTML = '';
            select.innerHTML = '<option value="">None</option>';
            
            isos.forEach(iso => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${iso.name}</td>
                    <td>${Utils.formatSize(iso.size)}</td>
                `;
                tbody.appendChild(tr);
                
                const option = document.createElement('option');
                option.value = iso.name;
                option.textContent = iso.name;
                select.appendChild(option);
            });
        } catch (error) {
            console.error('Failed to load ISOs:', error);
        }
    }

    static async upload(event) {
        event.preventDefault();
        const formData = new FormData();
        formData.append('iso', document.getElementById('isoFile').files[0]);
        
        try {
            await fetch('/api/isos', {
                method: 'POST',
                body: formData
            });
            document.getElementById('uploadISOForm').reset();
            ISOManager.load();
        } catch (error) {
            console.error('Failed to upload ISO:', error);
        }
    }
}
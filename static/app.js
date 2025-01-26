// Navigation
function showSection(sectionId) {
    document.querySelectorAll('.section').forEach(section => section.classList.remove('active'));
    document.querySelectorAll('.nav-item').forEach(item => item.classList.remove('active'));
    document.getElementById(sectionId).classList.add('active');
    document.querySelector(`[href="#${sectionId}"]`).classList.add('active');
}

// VM Management
class VMManager {
    static async load() {
        try {
            const response = await fetch('/api/vms');
            const vms = await response.json();
            const tbody = document.querySelector('#vmTable tbody');
            tbody.innerHTML = '';
            
            vms.forEach(vm => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${vm.name}</td>
                    <td>${vm.cores}</td>
                    <td>${vm.memory} MB</td>
                    <td><span class="status status-${vm.state}">${vm.state}</span></td>
                    <td>
                        ${vm.state === 'running' ? 
                            `<button class="btn btn-secondary" onclick="VMManager.stop('${vm.name}')">Stop</button>` : ''}
                        <button class="btn btn-danger" onclick="VMManager.delete('${vm.name}')">Delete</button>
                    </td>
                `;
                tbody.appendChild(tr);
            });
        } catch (error) {
            console.error('Failed to load VMs:', error);
        }
    }

    static async create(event) {
        event.preventDefault();
        const data = {
            name: document.getElementById('vmName').value,
            cores: parseInt(document.getElementById('vmCores').value),
            memory: parseInt(document.getElementById('vmMemory').value),
            iso: document.getElementById('vmIso').value
        };

        try {
            await fetch('/api/vms', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            ModalManager.hideCreateVM();
            VMManager.load();
        } catch (error) {
            console.error('Failed to create VM:', error);
        }
    }

    static async stop(name) {
        try {
            await fetch(`/api/vms/${name}/stop`, { method: 'POST' });
            setTimeout(VMManager.load, 1000);
        } catch (error) {
            console.error('Failed to stop VM:', error);
        }
    }

    static async delete(name) {
        if (!confirm(`Are you sure you want to delete ${name}?`)) return;
        
        try {
            await fetch(`/api/vms/${name}`, { method: 'DELETE' });
            setTimeout(VMManager.load, 1000);
        } catch (error) {
            console.error('Failed to delete VM:', error);
        }
    }
}

// ISO Management
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
        const file = document.getElementById('isoFile').files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('iso', file);
        
        // Create and show progress modal
        const progressModal = document.createElement('div');
        progressModal.className = 'progress-modal';
        progressModal.innerHTML = `
            <h3>Uploading ${file.name}</h3>
            <div class="progress-bar">
                <div class="progress-bar-fill"></div>
            </div>
            <div class="progress-text">0%</div>
        `;
        document.body.appendChild(progressModal);
        
        try {
            const xhr = new XMLHttpRequest();
            xhr.upload.onprogress = (e) => {
                if (e.lengthComputable) {
                    const percent = Math.round((e.loaded / e.total) * 100);
                    progressModal.querySelector('.progress-bar-fill').style.width = percent + '%';
                    progressModal.querySelector('.progress-text').textContent = percent + '%';
                }
            };
            
            await new Promise((resolve, reject) => {
                xhr.onload = () => {
                    if (xhr.status === 200) {
                        resolve();
                    } else {
                        reject(new Error('Upload failed'));
                    }
                };
                xhr.onerror = () => reject(new Error('Upload failed'));
                xhr.open('POST', '/api/isos');
                xhr.send(formData);
            });

            // Show completion message
            progressModal.innerHTML = `
                <h3>Upload Complete</h3>
                <p>${file.name} has been uploaded successfully.</p>
                <button class="btn" onclick="this.parentElement.remove()">Close</button>
            `;
            
            document.getElementById('uploadISOForm').reset();
            ISOManager.load();
        } catch (error) {
            progressModal.innerHTML = `
                <h3>Upload Failed</h3>
                <p>Failed to upload ${file.name}.</p>
                <button class="btn" onclick="this.parentElement.remove()">Close</button>
            `;
            console.error('Failed to upload ISO:', error);
        }
    }
}

// Modal Management
class ModalManager {
    static showCreateVM() {
        document.getElementById('createVMModal').style.display = 'block';
    }

    static hideCreateVM() {
        document.getElementById('createVMModal').style.display = 'none';
        document.getElementById('createVMForm').reset();
    }
}

// Utilities
class Utils {
    static formatSize(bytes) {
        const units = ['B', 'KB', 'MB', 'GB'];
        let size = bytes;
        let unitIndex = 0;
        
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        
        return `${size.toFixed(1)} ${units[unitIndex]}`;
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('createVMForm').addEventListener('submit', VMManager.create);
    document.getElementById('uploadISOForm').addEventListener('submit', ISOManager.upload);
    
    VMManager.load();
    ISOManager.load();
    setInterval(VMManager.load, 5000);
});
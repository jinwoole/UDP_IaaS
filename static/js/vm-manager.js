// vm-manager.js
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
                        ${vm.state === 'running' ? `
                            <button class="btn btn-secondary" onclick="VMManager.stop('${vm.name}')">Stop</button>
                            <button class="btn" onclick="VMManager.openVNC('${vm.name}')">Console</button>
                        ` : `
                            <button class="btn" onclick="VMManager.start('${vm.name}')">Start</button>
                        `}
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

    static async start(name) {
        try {
            const response = await fetch(`/api/vms/${name}/start`, { method: 'POST' });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            setTimeout(VMManager.load, 1000);
        } catch (error) {
            console.error('Failed to start VM:', error);
            alert(`Failed to start VM: ${error.message}`);
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

    static async openVNC(name) {
        try {
            const response = await fetch(`/api/vms/${name}/vnc`);
            const { port } = await response.json();
            window.open(`http://${window.location.hostname}:${port}/vnc.html`, '_blank');
        } catch (error) {
            console.error('Failed to open VNC:', error);
        }
    }
}
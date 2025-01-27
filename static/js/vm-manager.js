// vm-manager.js
class VMManager {
    static async load() {
        try {
            const response = await fetch('/api/vms');
            const vms = await response.json();
            const tbody = document.querySelector('#vmTable tbody');
            
            const stoppingVMs = new Set();
            const stateChecks = Array.from(tbody.querySelectorAll('tr')).map(async tr => {
                const vmName = tr.getAttribute('data-vm');
                const stateSpan = tr.querySelector('.status');
                if (stateSpan && stateSpan.textContent === 'stopping') {
                    const state = await VMManager.checkVMState(vmName);
                    if (state !== 'stopped') {
                        stoppingVMs.add(vmName);
                    }
                }
            });
            
            await Promise.all(stateChecks);
            tbody.innerHTML = '';
            
            vms.forEach(vm => {
                if (stoppingVMs.has(vm.name)) {
                    vm.state = 'stopping';
                }
    
                const tr = document.createElement('tr');
                tr.setAttribute('data-vm', vm.name);
                
                let actionButtons = '';
                if (vm.state === 'stopping') {
                    actionButtons = `<span class="status status-stopping">Stopping...</span>`;
                } else if (vm.state === 'running') {
                    actionButtons = `
                        <button class="btn btn-secondary" onclick="VMManager.stop('${vm.name}')">Stop</button>
                        <button class="btn" onclick="VMManager.openVNC('${vm.name}')">Console</button>
                    `;
                } else {
                    actionButtons = `<button class="btn" onclick="VMManager.start('${vm.name}')">Start</button>`;
                }
    
                const deleteButton = vm.state !== 'deleting' ? 
                    `<button class="btn btn-danger" onclick="VMManager.delete('${vm.name}')">Delete</button>` : 
                    `<span class="status status-deleting">Deleting...</span>`;
    
                tr.innerHTML = `
                    <td>${vm.name}</td>
                    <td>${vm.cores}</td>
                    <td>${vm.memory} MB</td>
                    <td><span class="status status-${vm.state}">${vm.state}</span></td>
                    <td>${actionButtons} ${deleteButton}</td>
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
    // vm-manager.js
    static async checkStateUntilStopped(name) {
        const maxRetries = 30;
        let retries = 0;
    
        const check = async () => {
            const stateResponse = await fetch(`/api/vms/${name}/state`);
            if (!stateResponse.ok) return;
            const { state } = await stateResponse.json();
            
            if (state === 'stopped') {
                await VMManager.load();
                return true;
            }
            
            if (retries++ < maxRetries) {
                setTimeout(() => check(), 500); // 0.5초마다 체크
            }
        };
    
        return check();
    }
    
    static async stop(name) {
        try {
            // UI를 stopping으로 변경
            const tr = document.querySelector(`tr[data-vm="${name}"]`);
            tr.cells[3].innerHTML = '<span class="status status-stopping">stopping</span>';
            tr.cells[4].innerHTML = '<span class="status status-stopping">Stopping...</span>';
            
            // 중지 요청 및 상태 체크
            await fetch(`/api/vms/${name}/stop`, { method: 'POST' });
            await this.checkStateUntilStopped(name);
        } catch (error) {
            console.error('Failed to stop VM:', error);
            alert(error.message);
            await VMManager.load();
        }
    }

    static async delete(name) {
        if (!confirm(`Are you sure you want to delete ${name}?`)) return;
        
        try {
            const tr = document.querySelector(`tr[data-vm="${name}"]`);
            const stateCell = tr.cells[3];
            const actionCell = tr.cells[4];
            
            stateCell.innerHTML = '<span class="status status-deleting">deleting</span>';
            actionCell.innerHTML = '<span class="status status-deleting">Deleting...</span>';
            
            await fetch(`/api/vms/${name}`, { method: 'DELETE' });
            await new Promise(resolve => setTimeout(resolve, 5000));
            await VMManager.load();
        } catch (error) {
            console.error('Failed to delete VM:', error);
            await VMManager.load();
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
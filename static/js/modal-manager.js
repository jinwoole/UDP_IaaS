// modal-manager.js
class ModalManager {
    static showCreateVM() {
        document.getElementById('createVMModal').style.display = 'block';
    }

    static hideCreateVM() {
        document.getElementById('createVMModal').style.display = 'none';
        document.getElementById('createVMForm').reset();
    }
}
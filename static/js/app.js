// app.js
document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('createVMForm').addEventListener('submit', VMManager.create);
    document.getElementById('uploadISOForm').addEventListener('submit', ISOManager.upload);
    
    VMManager.load();
    ISOManager.load();
    setInterval(VMManager.load, 5000);
});
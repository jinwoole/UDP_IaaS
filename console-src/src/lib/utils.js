export function formatSize(bytes) {
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    
    return `${size.toFixed(1)} ${units[unitIndex]}`;
  }
  
  export function formatState(state) {
    const baseClasses = 'px-2 py-1 rounded text-sm';
    switch (state) {
      case 'running':
        return `${baseClasses} bg-green-100 text-green-800`;
      case 'stopped':
        return `${baseClasses} bg-red-100 text-red-800`;
      case 'paused':
        return `${baseClasses} bg-yellow-100 text-yellow-800`;
      case 'stopping':
      case 'starting':
        return `${baseClasses} bg-blue-100 text-blue-800 animate-pulse`;
      case 'deleting':
        return `${baseClasses} bg-gray-100 text-gray-800 animate-pulse`;
      default:
        return `${baseClasses} bg-gray-100 text-gray-800`;
    }
  }
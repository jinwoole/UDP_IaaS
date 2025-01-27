<script>
    import { onMount } from 'svelte';
    import { formatSize } from '$lib/utils';
  
    let isos = [];
    let selectedFile = null;
  
    async function loadISOs() {
      try {
        const res = await fetch('/api/isos');
        isos = await res.json();
      } catch (error) {
        console.error('Failed to load ISOs:', error);
      }
    }
  
    async function handleUpload(event) {
      event.preventDefault();
      if (!selectedFile) return;
  
      const formData = new FormData();
      formData.append('iso', selectedFile);
  
      try {
        await fetch('/api/isos', {
          method: 'POST',
          body: formData
        });
        selectedFile = null;
        event.target.reset();
        await loadISOs();
      } catch (error) {
        console.error('Failed to upload ISO:', error);
      }
    }
  
    onMount(loadISOs);
  </script>
  
  <div class="min-h-screen bg-gradient-to-r from-gray-200 to-gray-400 flex flex-col items-center justify-center px-4 py-8">
    <div class="w-full max-w-3xl bg-white rounded-3xl shadow-2xl p-8">
      <h1 class="text-2xl font-semibold text-gray-800 mb-6 text-center">ISO Management</h1>
      <form class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-4 mb-6" on:submit={handleUpload}>
        <input
          class="block w-full sm:w-auto text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:bg-blue-600 file:text-white hover:file:bg-blue-700 cursor-pointer"
          type="file" accept=".iso" required
          on:change={(e) => selectedFile = e.target.files[0]}
        />
        <button
          type="submit"
          class="bg-blue-600 text-white px-6 py-2 rounded-full hover:bg-blue-800 transition-colors focus:outline-none"
        >
          Upload
        </button>
      </form>
      <div class="overflow-hidden rounded-xl">
        <table class="min-w-full border-collapse">
          <thead class="bg-gray-100">
            <tr>
              <th class="py-3 px-4 text-left text-gray-600">Name</th>
              <th class="py-3 px-4 text-left text-gray-600">Size</th>
            </tr>
          </thead>
          <tbody>
            {#each isos as iso (iso.name)}
              <tr class="border-b last:border-none hover:bg-gray-50">
                <td class="py-3 px-4 text-gray-700">{iso.name}</td>
                <td class="py-3 px-4 text-gray-700">{formatSize(iso.size)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
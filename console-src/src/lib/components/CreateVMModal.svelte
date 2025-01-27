<script>
    import { createEventDispatcher, onMount } from 'svelte';
    
    const dispatch = createEventDispatcher();
    
    let isos = [];
    let formData = {
      name: '',
      cores: 1,
      memory: 1024,
      iso: ''
    };
  
    async function loadISOs() {
      try {
        const res = await fetch('/api/isos');
        isos = await res.json();
      } catch (error) {
        console.error('Failed to load ISOs:', error);
      }
    }
  
    async function handleSubmit() {
      try {
        await fetch('/api/vms', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(formData)
        });
        dispatch('created');
        dispatch('close');
      } catch (error) {
        console.error('Failed to create VM:', error);
      }
    }
  
    onMount(loadISOs);
  </script>
  
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
    <div class="bg-white rounded-lg p-6 max-w-md w-full">
      <h2 class="text-xl font-bold mb-4">Create Virtual Machine</h2>
      
      <form on:submit|preventDefault={handleSubmit}>
        <div class="mb-4">
          <label class="block mb-1">Name</label>
          <input
            type="text"
            class="w-full border rounded px-2 py-1"
            bind:value={formData.name}
            required
          />
        </div>
  
        <div class="mb-4">
          <label class="block mb-1">CPU Cores</label>
          <input
            type="number"
            class="w-full border rounded px-2 py-1"
            bind:value={formData.cores}
            min="1"
            required
          />
        </div>
  
        <div class="mb-4">
          <label class="block mb-1">Memory (MB)</label>
          <input
            type="number"
            class="w-full border rounded px-2 py-1"
            bind:value={formData.memory}
            min="512"
            required
          />
        </div>
  
        <div class="mb-4">
          <label class="block mb-1">ISO Image</label>
          <select
            class="w-full border rounded px-2 py-1"
            bind:value={formData.iso}
          >
            <option value="">None</option>
            {#each isos as iso}
              <option value={iso.name}>{iso.name}</option>
            {/each}
          </select>
        </div>
  
        <div class="flex justify-end space-x-2">
          <button
            type="button"
            class="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700"
            on:click={() => dispatch('close')}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Create
          </button>
        </div>
      </form>
    </div>
  </div>
  
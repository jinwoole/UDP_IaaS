<script>
    import { onMount, onDestroy } from 'svelte';
    import CreateVMModal from '$lib/components/CreateVMModal.svelte';
    import { formatState } from '$lib/utils';
  
    let vms = [];
    let showModal = false;
    let interval;
  
    async function loadVMs() {
      try {
        const res = await fetch('/api/vms');
        vms = await res.json();
      } catch (error) {
        console.error('Failed to load VMs:', error);
      }
    }
  
    async function handleStart(name) {
      try {
        await fetch(`/api/vms/${name}/start`, { method: 'POST' });
        await loadVMs();
      } catch (error) {
        console.error('Failed to start VM:', error);
      }
    }
  
    async function handleStop(name) {
      try {
        await fetch(`/api/vms/${name}/stop`, { method: 'POST' });
        await loadVMs();
      } catch (error) {
        console.error('Failed to stop VM:', error);
      }
    }
  
    async function handleDelete(name) {
      if (!confirm(`Are you sure you want to delete ${name}?`)) return;
      
      try {
        await fetch(`/api/vms/${name}`, { method: 'DELETE' });
        await loadVMs();
      } catch (error) {
        console.error('Failed to delete VM:', error);
      }
    }
  
    async function handleVNC(name) {
      try {
        const res = await fetch(`/api/vms/${name}/vnc`);
        const { port } = await res.json();
        window.open(`http://${window.location.hostname}:${port}/vnc.html`, '_blank');
      } catch (error) {
        console.error('Failed to open VNC:', error);
      }
    }
  
    onMount(() => {
      loadVMs();
      interval = setInterval(loadVMs, 5000);
    });
  
    onDestroy(() => {
      clearInterval(interval);
    });
  </script>
  
  <div class="flex items-center justify-center min-h-screen">
    <h1 class="text-2xl font-semibold text-gray-800">
      Welcome to UDP IaaS
    </h1>
  </div>
  
  <div class="bg-white rounded-lg shadow p-6">
    <div class="mb-4">
      <button 
        class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        on:click={() => showModal = true}
      >
        Create VM
      </button>
    </div>
  
    <table class="w-full">
      <thead>
        <tr class="border-b">
          <th class="text-left py-2">Name</th>
          <th class="text-left py-2">Cores</th>
          <th class="text-left py-2">Memory</th>
          <th class="text-left py-2">State</th>
          <th class="text-left py-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each vms as vm (vm.name)}
          <tr class="border-b">
            <td class="py-2">{vm.name}</td>
            <td class="py-2">{vm.cores}</td>
            <td class="py-2">{vm.memory} MB</td>
            <td class="py-2">
              <span class={formatState(vm.state)}>
                {vm.state}
              </span>
            </td>
            <td class="py-2 space-x-2">
              {#if vm.state === 'running'}
                <button
                  class="bg-gray-600 text-white px-3 py-1 rounded hover:bg-gray-700"
                  on:click={() => handleStop(vm.name)}
                >
                  Stop
                </button>
                <button
                  class="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
                  on:click={() => handleVNC(vm.name)}
                >
                  Console
                </button>
              {:else if vm.state === 'stopped'}
                <button
                  class="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
                  on:click={() => handleStart(vm.name)}
                >
                  Start
                </button>
              {/if}
              <button
                class="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700"
                on:click={() => handleDelete(vm.name)}
              >
                Delete
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
  
  {#if showModal}
    <CreateVMModal 
      on:close={() => showModal = false}
      on:created={loadVMs}
    />
  {/if}
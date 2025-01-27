<script>
import '../app.css';
import { page } from '$app/stores';
export const prerender = true;

</script>

<div class="flex min-h-screen bg-gradient-to-tr from-gray-100 to-white">
  <div class="relative flex flex-col w-64 backdrop-blur bg-white/70 text-gray-700 p-4 rounded-r-xl shadow-xl">
    <div class="mb-2 p-4">
		<a href="/" class="text-xl font-semibold">UDP IaaS</a>
    </div>
    <nav class="flex flex-col gap-2 p-2 text-base text-gray-700">
      <a href="/vm" class={`${$page.url.pathname === '/vm' ? 'bg-gray-200' : ''} flex items-center w-full p-3 rounded-lg hover:bg-gray-100`}>
        <span class="mr-4">🔗</span>Virtual Machines
      </a>
      <a href="/isos" class={`${$page.url.pathname === '/isos' ? 'bg-gray-200' : ''} flex items-center w-full p-3 rounded-lg hover:bg-gray-100`}>
        <span class="mr-4">📄</span>ISO Images
      </a>
    </nav>
  </div>

  <div class="flex-1">
    <header class="backdrop-blur bg-white/70 shadow rounded-b-xl">
    </header>
    <main class="max-w-7xl mx-auto px-4 py-8">
      <slot />
    </main>
  </div>
</div>
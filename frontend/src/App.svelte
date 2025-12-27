<script>
  import { onMount } from 'svelte';

  let message = '';
  let items = [];
  let newItem = '';
  let loading = false;

  const API_URL = 'http://localhost:8080/api';

  onMount(async () => {
    fetchHello();
    fetchItems();
  });

  async function fetchHello() {
    try {
      const response = await fetch(`${API_URL}/hello`);
      const data = await response.json();
      message = data.message;
    } catch (error) {
      console.error('Error:', error);
    }
  }

  async function fetchItems() {
    loading = true;
    try {
      const response = await fetch(`${API_URL}/items`);
      if (!response.ok) throw new Error('Failed to fetch');
      const data = await response.json();
      items = data.items || [];
    } catch (error) {
      console.error('Error fetching items:', error);
    } finally {
      loading = false;
    }
  }

  async function deleteItem(index) {
    try {
      const response = await fetch(`${API_URL}/items`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ index }),
      });

      if (response.ok) {
        await fetchItems();
      }
    } catch (error) {
      console.error('Error:', error);
    }
  }

  async function addItem() {
    if (!newItem.trim()) return;

    try {
      const response = await fetch(`${API_URL}/items`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ text: newItem }),
      });

      if (response.ok) {
        await fetchItems();
        newItem = '';
      }
    } catch (error) {
      console.error('Error:', error);
    }
  }
</script>

<main class="min-h-screen bg-gray-50 py-12 px-4">
  <div class="max-w-2xl mx-auto">
    <h1 class="text-4xl font-bold text-center text-gray-900 mb-8">
      Svelte + Go App
    </h1>

    <div class="bg-white rounded-xl shadow-md p-6 mb-6">
      <h2 class="text-2xl font-semibold text-gray-800 mb-2">{message}</h2>
    </div>

    <div class="bg-white rounded-xl shadow-md p-6">
      <h3 class="text-xl font-semibold text-gray-800 mb-4">Items</h3>

      {#if items.length === 0}
        <p class="text-gray-500 text-center py-4">No items yet</p>
      {:else}
        <ul class="space-y-2 mb-4">
          {#each items as item, index}
            <li class="bg-gray-100 rounded-lg px-4 py-3 text-gray-700 flex justify-between items-center">
              <span>{item.text}</span>
              <button
                onclick={() => deleteItem(index)}
                class="text-red-500 hover:text-red-700 font-medium px-3 py-1 rounded hover:bg-red-50 transition-colors"
              >
                Delete
              </button>
            </li>
          {/each}
        </ul>
      {/if}

      <div class="flex gap-3 mt-4">
        <input
          type="text"
          bind:value={newItem}
          placeholder="Enter new item"
          class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          onkeypress={(e) => e.key === 'Enter' && addItem()}
        />
        <button
          onclick={addItem}
          class="px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
        >
          Add Item
        </button>
      </div>
    </div>
  </div>
</main>

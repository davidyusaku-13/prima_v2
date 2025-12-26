<script>
  import { onMount } from 'svelte';

  let message = '';
  let items = [];
  let newItem = '';
  let loading = false;

  const API_URL = 'http://localhost:8080/api';

  onMount(async () => {
    await fetchHello();
    await fetchItems();
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
      const data = await response.json();
      items = data.items;
    } catch (error) {
      console.error('Error:', error);
    } finally {
      loading = false;
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

<main>
  <h1>Svelte + Go App</h1>
  
  <div class="card">
    <h2>{message}</h2>
  </div>

  <div class="card">
    <h3>Items</h3>
    {#if loading}
      <p>Loading...</p>
    {:else}
      <ul>
        {#each items as item}
          <li>{item}</li>
        {/each}
      </ul>
    {/if}

    <div class="input-group">
      <input 
        type="text" 
        bind:value={newItem} 
        placeholder="Enter new item"
        on:keypress={(e) => e.key === 'Enter' && addItem()}
      />
      <button on:click={addItem}>Add Item</button>
    </div>
  </div>
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }

  .card {
    padding: 1.5rem;
    margin: 1rem 0;
    border-radius: 8px;
    background: #f5f5f5;
  }

  .input-group {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  input {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  button {
    padding: 0.5rem 1rem;
    background: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  button:hover {
    background: #45a049;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  li {
    padding: 0.5rem;
    margin: 0.5rem 0;
    background: white;
    border-radius: 4px;
  }
</style>
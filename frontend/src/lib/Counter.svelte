<script>
  let count = $state(0)
  const increment = () => {
    count += 1
  }
</script>

<button
  onclick={increment}
  class="px-4 py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors"
>
  count is {count}
</button>

<script>
  import Header from '$lib/components/Header.svelte';
  import TaskForm from '$lib/components/TaskForm.svelte';
  import TaskList from '$lib/components/TaskList.svelte';
  import ErrorMessage from '$lib/components/ErrorMessage.svelte';
  import { onMount } from 'svelte';
  import { user } from '$lib/stores.js';

  onMount(() => {
      const token = localStorage.getItem("jwt");
      if (!token) {
          window.location.href = "/login";
          return;
      }

      // Get username from JWT token
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));

      user.set(JSON.parse(jsonPayload).username);
  });
</script>

<ErrorMessage />
<Header />
<TaskForm />
<TaskList />

<style>
  :global(body) {
      font-family: Arial, sans-serif;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
  }
</style>
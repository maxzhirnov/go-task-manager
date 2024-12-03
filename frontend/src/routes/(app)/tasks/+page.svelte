<script>
  import Header from '$lib/components/Header.svelte';
  import TaskForm from '$lib/components/TaskForm.svelte';
  import TaskList from '$lib/components/TaskList.svelte';
  import ErrorMessage from '$lib/components/ErrorMessage.svelte';
  import { onMount } from 'svelte';
  import { user } from '$lib/stores.js';

</script>

<svelte:head>
  <title>Dashboard | ActionHub</title>
</svelte:head>

<TaskForm />
<TaskList />


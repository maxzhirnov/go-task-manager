<script>
  import { onMount } from 'svelte';
  import { user } from '$lib/stores.js';

  import ErrorMessage from '$lib/components/Common/ErrorMessage.svelte';
  import Header from '$lib/components/Header/Header.svelte';
  import TaskForm from '$lib/components/TaskForm/TaskForm2.svelte';
  import TaskList from '$lib/components/TaskList/TaskList.svelte';
</script>

<svelte:head>
  <title>Dashboard | ActionHub</title>
</svelte:head>

<TaskForm />
<TaskList />


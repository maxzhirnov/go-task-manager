<script>
    import { onMount } from 'svelte';
    import { tasks } from '../stores.js';
    import { api } from '../api.js';
    import TaskItem from './TaskItem.svelte';

    onMount(async () => {
        const fetchedTasks = await api.fetchTasks();
        tasks.set(fetchedTasks);
    });
</script>

<ul class="task-list">
    {#each $tasks as task (task.id)}
        <TaskItem {task} />
    {/each}
</ul>

<style>
    .task-list {
        list-style-type: none;
        padding: 0;
    }
</style>
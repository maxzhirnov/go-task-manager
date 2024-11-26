<script>
    import { dndzone } from 'svelte-dnd-action';
    import { tasks } from '../stores.js';
    import TaskItem from './TaskItem.svelte';
    import { api } from '$lib/api.js';
    import { showError } from '$lib/stores.js';
    import { onMount } from 'svelte';

    let dragDisabled = false;

    // Load tasks when component mounts
    onMount(async () => {
        try {
            const fetchedTasks = await api.fetchTasks();
            tasks.set(fetchedTasks);
        } catch (error) {
            showError('Failed to load tasks');
        }
    });

    function handleDndConsider(e) {
        tasks.set(e.detail.items);
    }

    async function handleDndFinalize(e) {
        const newTasks = e.detail.items;
        // console.log('DnD Finalize event:', e.detail);
        
        // Create positions map
        const positions = {};
        newTasks.forEach((task, index) => {
            positions[task.id] = index;
        });
        // console.log('Positions to update:', positions);

        try {
            await api.updateTaskPositions(positions);
            // Refresh tasks after successful update
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            // console.log('Error updating task positions:', error);
            showError('Failed to update task positions');
            // Refresh tasks to restore original order
            const originalTasks = await api.fetchTasks();
            tasks.set(originalTasks);
        }
    }
</script>

<section
    use:dndzone={{items: $tasks, dragDisabled}}
    on:consider={handleDndConsider}
    on:finalize={handleDndFinalize}
>
    {#each $tasks as task (task.id)}
        <TaskItem {task} />
    {/each}
</section>
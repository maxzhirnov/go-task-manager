<script>
    import { dragHandleZone } from 'svelte-dnd-action';
    import { tasks } from '../stores.js';
    import TaskItem from './TaskItem.svelte';
    import { api } from '$lib/api.js';
    import { showError } from '$lib/stores.js';
    import { onMount } from 'svelte';
    import {flip} from "svelte/animate";

    const flipDurationMs = 300;
    const dragDisabled = false;

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
            // console.log(`Task ${task.id} moved to position ${index}`);
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
    use:dragHandleZone={{
        items: $tasks,
        flipDurationMs,
    }}
    on:consider={handleDndConsider}
    on:finalize={handleDndFinalize}
>
    {#each $tasks as task (task.id)}
        <div class="task-wrapper" animate:flip="{{ duration: flipDurationMs }}">
            <TaskItem {task} />
        </div>
    {/each}
</section>

<style>
    .task-wrapper {
        width: 100%;
    }
</style>
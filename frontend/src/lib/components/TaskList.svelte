<!-- TaskList.svelte -->
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

    // Add status filter
    let selectedStatus = 'active'; // Default to active tasks
    const statusOptions = [
        { value: 'active', label: 'Active' },  // "Active" means all except completed
        { value: 'all', label: 'All' },
        { value: 'pending', label: 'Pending' },
        { value: 'in_progress', label: 'In Progress' },
        { value: 'completed', label: 'Completed' }
    ];

    // Modified filtered tasks computation
    $: filteredTasks = 
        selectedStatus === 'active' ? $tasks.filter(task => task.status !== 'completed') :
        selectedStatus === 'all' ? $tasks :
        $tasks.filter(task => task.status === selectedStatus);

    onMount(async () => {
        try {
            const fetchedTasks = await api.fetchTasks();
            tasks.set(fetchedTasks);
        } catch (error) {
            showError('Failed to load tasks');
        }
    });

    function handleDndConsider(e) {
        // Only update positions within the filtered view
        const updatedTasks = [...$tasks];
        const filteredIndices = $tasks.map((task, index) => 
            selectedStatus === 'all' || task.status === selectedStatus ? index : -1
        ).filter(index => index !== -1);

        e.detail.items.forEach((task, newFilteredIndex) => {
            updatedTasks[filteredIndices[newFilteredIndex]] = task;
        });

        tasks.set(updatedTasks);
    }

    async function handleDndFinalize(e) {
        const newTasks = [...$tasks];
        const positions = {};
        
        newTasks.forEach((task, index) => {
            positions[task.id] = index;
        });

        try {
            await api.updateTaskPositions(positions);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            showError('Failed to update task positions');
            const originalTasks = await api.fetchTasks();
            tasks.set(originalTasks);
        }
    }
</script>

<div class="filter-container">
    <select bind:value={selectedStatus}>
        {#each statusOptions as option}
            <option value={option.value}>{option.label}</option>
        {/each}
    </select>
</div>

<section
    use:dragHandleZone={{
        items: filteredTasks,
        flipDurationMs,
    }}
    on:consider={handleDndConsider}
    on:finalize={handleDndFinalize}
>
    {#each filteredTasks as task (task.id)}
        <div class="task-wrapper" animate:flip="{{ duration: flipDurationMs }}">
            <TaskItem {task} />
        </div>
    {/each}
</section>

<style>
    .task-wrapper {
        width: 100%;
    }

    .filter-container {
        margin-bottom: 1rem;
    }

    select {
        padding: 8px;
        border-radius: 4px;
        border: 1px solid #ddd;
        font-size: 14px;
        width: 200px;
    }

    select:focus {
        outline: none;
        border-color: #2196F3;
    }
</style>
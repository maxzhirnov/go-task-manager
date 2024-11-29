<!-- TaskList.svelte -->
<script>
    import { dragHandleZone } from 'svelte-dnd-action';
    import { tasks } from '../stores.js';
    import TaskItem from './TaskItem.svelte';
    import { api } from '$lib/api.js';
    import { showError } from '$lib/stores.js';
    import { onMount } from 'svelte';
    import {flip} from "svelte/animate";
    import { mdiFilter } from '@mdi/js';

    const flipDurationMs = 300;
    const dragDisabled = false;

    // Add status filter
    let selectedStatus = 'active'; // Default to active tasks
    const statusOptions = [
        { value: 'active', label: 'Active' },  // "Active" means all except completed
        { value: 'pending', label: 'Pending' },
        { value: 'in_progress', label: 'In Progress' },
        { value: 'all', label: 'All' },
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
        const updatedTasks = [...$tasks];
        const filteredIndices = $tasks.map((task, index) => {
            if (selectedStatus === 'all') return index;
            if (selectedStatus === 'active') return task.status !== 'completed' ? index : -1;
            return task.status === selectedStatus ? index : -1;
        }).filter(index => index !== -1);

        e.detail.items.forEach((task, newFilteredIndex) => {
            updatedTasks[filteredIndices[newFilteredIndex]] = task;
        });

        tasks.set(updatedTasks);
    }

    async function handleDndFinalize(e) {
        const positions = {};
        
        // Get only the visible tasks based on current filter
        const visibleTasks = selectedStatus === 'active' 
            ? $tasks.filter(task => task.status !== 'completed')
            : selectedStatus === 'all' 
                ? $tasks 
                : $tasks.filter(task => task.status === selectedStatus);

        // Update positions only for visible tasks
        visibleTasks.forEach((task, index) => {
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
<div class="task-list-container">
    <div class="filter-container">
        <svg viewBox="0 0 24 24" width="24" height="24">
            <path fill="currentColor" d={mdiFilter} />
        </svg>
        <select bind:value={selectedStatus}>
            {#each statusOptions as option}
                <option value={option.value}>
                    {option.label}
                </option>
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
</div>
<style>
    .task-list-container {
        margin-top: 2.5rem;
    }

    .task-wrapper {
        width: 100%;
    }

    .filter-container {
        display: flex;
        margin-bottom: 1rem;
        align-items: center;
        gap: 8px;
    }

    .filter-container select {
        width: 150px;
        appearance: none;
        padding: 8px 32px 8px 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: #fff url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23666%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.4-12.8z%22%2F%3E%3C%2Fsvg%3E") no-repeat right 12px top 50%;
        background-size: 12px auto;
        font-size: 14px;
        cursor: pointer;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    .filter-container select:hover {
        border-color: #aaa;
    }

    .filter-container select:focus {
        outline: none;
        border-color: #2196F3;
        box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
    }

    .filter-container select option {
        padding: 8px;
        background-color: white;
    }


    @media screen and (max-width: 600px) {
        .filter-container select {
            width: 100%;
        }
    }
</style>
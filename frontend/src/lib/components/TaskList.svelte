<!-- TaskList.svelte -->
<script>
    import { dragHandleZone } from 'svelte-dnd-action';

    import {flip} from "svelte/animate";
    import { onMount } from 'svelte';

    import { browser } from '$app/environment';
    import { tasks } from '$lib/stores.js';
    import { api } from '$lib/api.js';
    import { showError } from '$lib/stores.js';

    import LoadingSpinner from './LoadingSpinner.svelte';
    import TaskItem from './TaskItem.svelte';
    import FilterButtons from './FilterButtons.svelte';

    let isLoading = true; //
    const flipDurationMs = 300;
    const dragDisabled = false;

    let selectedStatus = browser ? localStorage.getItem('selectedStatus') || 'active' : 'active';

    $: if (browser && selectedStatus) {
        localStorage.setItem('selectedStatus', selectedStatus);
    }

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
            isLoading = true;
            const fetchedTasks = await api.fetchTasks();
            tasks.set(fetchedTasks);
        } catch (error) {
            showError('Failed to load tasks');
        }  finally {
            isLoading = false;
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
    <FilterButtons bind:selectedStatus {statusOptions} />
    
    <!-- Task list -->
    {#if isLoading}
        <LoadingSpinner/>
    {:else}
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
    {/if}
    </div>

<style>
    .task-list-container {
        margin-top: 2.5rem;
    }

    .task-wrapper {
        width: 100%;
    }
</style>
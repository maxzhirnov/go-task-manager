<script>
    import { onMount } from 'svelte';
    import { api } from '$lib/api';
    import { showError } from '$lib/stores';

    export let statistics = null;

    async function loadStatistics() {
        try {
            statistics = await api.getUserStatistics();
        } catch (error) {
            showError("Failed to load statistics");
        }
    }

    onMount(loadStatistics); // Add this line
</script>

<div class="statistics-container">
    {#if !statistics}
        <p>Loading statistics...</p>
    {:else}
        <div class="stats-grid">
            <div class="stat-card">
                <h3>Total Tasks</h3>
                <p>{statistics.total_tasks}</p>
            </div>
            <div class="stat-card">
                <h3>Completed</h3>
                <p>{statistics.completed_tasks}</p>
            </div>
            <div class="stat-card">
                <h3>In Progress</h3>
                <p>{statistics.in_progress_tasks}</p>
            </div>
            <div class="stat-card">
                <h3>Pending</h3>
                <p>{statistics.pending_tasks}</p>
            </div>
            <div class="stat-card">
                <h3>Created Today</h3>
                <p>{statistics.tasks_created_today}</p>
            </div>
        </div>
    {/if}
</div>

<style>
    .statistics-container {
        padding: 1rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        margin-top: 1rem;
    }

    .stat-card {
        background: #f5f5f5;
        padding: 1rem;
        border-radius: 8px;
        text-align: center;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .stat-card h3 {
        margin: 0 0 0.5rem 0;
        color: #666;
    }

    .stat-card p {
        font-size: 2rem;
        margin: 0;
        font-weight: bold;
    }
</style>
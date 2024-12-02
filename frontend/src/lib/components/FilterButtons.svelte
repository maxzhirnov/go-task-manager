<!-- FilterButtons.svelte -->
<script>
    import { mdiFilter } from '@mdi/js';

    export let selectedStatus;
    export let statusOptions;
</script>

<div class="filter-container">
    <svg viewBox="0 0 24 24" width="24" height="24" class="filter-icon">
        <path fill="currentColor" d={mdiFilter} />
    </svg>
    
    <!-- Desktop buttons -->
    <div class="filter-buttons">
        {#each statusOptions as option}
            <button 
                class="filter-btn {selectedStatus === option.value ? 'active' : ''}"
                on:click={() => selectedStatus = option.value}
            >
                {option.label}
            </button>
        {/each}
    </div>
    
    <!-- Mobile dropdown -->
    <select bind:value={selectedStatus} class="filter-select">
        {#each statusOptions as option}
            <option value={option.value}>
                {option.label}
            </option>
        {/each}
    </select>
</div>

<style>
    .filter-container {
        display: flex;
        margin-bottom: 1rem;
        align-items: center;
        gap: 8px;
    }

    .filter-buttons {
        display: none;
        gap: 8px;
    }

    .filter-btn {
        padding: 8px 16px;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: white;
        cursor: pointer;
        transition: all 0.2s;
    }

    .filter-btn:hover {
        border-color: #2196F3;
    }

    .filter-btn.active {
        background: #2196F3;
        color: white;
        border-color: #2196F3;
    }

    .filter-select {
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

    @media screen and (min-width: 601px) {
        .filter-buttons {
            display: flex;
        }
        
        .filter-select {
            display: none;
        }
    }

    @media screen and (max-width: 600px) {
        .filter-buttons {
            display: none;
        }
        
        .filter-select {
            display: block;
            width: 100%;
        }
    }
</style>
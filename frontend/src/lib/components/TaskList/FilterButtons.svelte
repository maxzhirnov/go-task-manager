<script>
    import FilterIcon from './FilterIcon.svelte';

    export let selectedStatus;
    export let statusOptions;
</script>

<div class="filter-container">
    <div class="terminal-box">
        <div class="terminal-header">
            <span class="terminal-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </span>
            <span class="terminal-title">FILTER.sys</span>
        </div>
        
        <!-- Desktop buttons -->
        <div class="filter-buttons">
            <FilterIcon/>
            
            {#each statusOptions as option}
                <button 
                    class="filter-btn {selectedStatus === option.value ? 'active' : ''}"
                    on:click={() => selectedStatus = option.value}
                >
                    <span class="status-indicator"></span>
                    <span class="btn-text">{option.label}</span>
                </button>
            {/each}
        </div>
        
        <!-- Mobile select -->
        <div class="select-wrapper">
            <div class="icon"><FilterIcon/></div>
            <select bind:value={selectedStatus} class="filter-select">
                {#each statusOptions as option}
                    <option value={option.value}>
                        {option.label}
                    </option>
                {/each}
            </select>
        </div>
    </div>
</div>

<style>
    .icon {
        margin-right: .5rem;
    }
    .filter-container {
        margin-bottom: 1rem;
        font-family: "JetBrains Mono", monospace;
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
    }

    .terminal-header {
        background: #2d3436;
        padding: 0.3rem 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
    }

    .filter-buttons {
        display: none;
        gap: 0.5rem;
        padding: 0.5rem;
        align-items: center;
    }

    .filter-icon {
        color: #00b894;
        margin-right: 0.5rem;
    }

    .filter-btn {
        background: transparent;
        border: 1px solid #0984e3;
        color: #0984e3;
        padding: 0.4rem 0.8rem;
        border-radius: 3px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
        transition: all 0.3s ease;
    }

    .filter-btn:hover {
        background: rgba(9, 132, 227, 0.1);
        box-shadow: 0 0 8px rgba(9, 132, 227, 0.2);
    }

    .filter-btn.active {
        background: rgba(0, 184, 148, 0.1);
        border-color: #00b894;
        color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .status-indicator {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: currentColor;
    }

    .select-wrapper {
        display: none;
        padding: 0.5rem;
        position: relative;
    }

    .filter-select {
        width: 100%;
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #fff;
        padding: 0.4rem 2rem;
        border-radius: 3px;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
        appearance: none;
        cursor: pointer;
    }

    .filter-select:focus {
        outline: none;
        border-color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    @media screen and (min-width: 601px) {
        .filter-buttons {
            display: flex;
        }
        
        .select-wrapper {
            display: none;
        }
    }

    @media screen and (max-width: 600px) {
        .filter-buttons {
            display: none;
        }
        
        .select-wrapper {
            display: flex;
            align-items: center;
        }
    }
</style>
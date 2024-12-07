<script>
    import FilterIcon from './FilterIcon.svelte';

    export let selectedStatus;
    export let statusOptions;
</script>

<div class="filter-container">
    <div class="filter-group">
        {#each statusOptions as option}
            <button 
                class="cyber-filter {selectedStatus === option.value ? 'active' : ''}"
                on:click={() => selectedStatus = option.value}
            >
                <div class="filter-content">
                    <div class="hex-border">
                        <svg viewBox="0 0 24 24" class="hex-svg">
                            <polygon 
                                points="12,2 22,7 22,17 12,22 2,17 2,7" 
                                class="hex-shape"
                            />
                        </svg>
                        <div class="status-core"></div>
                    </div>
                    <span class="filter-text">{option.label}</span>
                </div>
                <div class="filter-glow"></div>
            </button>
        {/each}
    </div>

    <!-- Mobile select -->
    <div class="mobile-filter">
        <FilterIcon />
        <select bind:value={selectedStatus} class="filter-select">
            {#each statusOptions as option}
                <option value={option.value}>{option.label}</option>
            {/each}
        </select>
    </div>
</div>

<style>
    .filter-container {
        margin-bottom: 1rem;
        font-family: "JetBrains Mono", monospace;
    }

    .filter-group {
        display: none;
        gap: 1rem;
        align-items: center;
        background: #1c1c1c;
        padding: 0.8rem;
        border-radius: 4px;
        border: 1px solid rgba(9, 132, 227, 0.2);
    }

    .cyber-filter {
        position: relative;
        background: transparent;
        border: none;
        padding: 0.5rem 1rem;
        cursor: pointer;
        color: #636e72;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .filter-content {
        position: relative;
        z-index: 2;
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .hex-border {
        position: relative;
        width: 24px;
        height: 24px;
    }

    .hex-svg {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }

    .hex-shape {
        fill: none;
        stroke: currentColor;
        stroke-width: 1;
        transition: all 0.3s ease;
    }

    .status-core {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 6px;
        height: 6px;
        background: currentColor;
        clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
        transition: all 0.3s ease;
    }

    .filter-text {
        font-size: 0.8rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .filter-glow {
        position: absolute;
        inset: 0;
        background: radial-gradient(
            circle at var(--x, 50%) var(--y, 50%),
            rgba(0, 184, 148, 0.1),
            transparent 50%
        );
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    /* Active state */
    .cyber-filter.active {
        color: #00b894;
    }

    .cyber-filter.active .hex-shape {
        stroke: #00b894;
        stroke-width: 2;
    }

    .cyber-filter.active .status-core {
        background: #00b894;
        box-shadow: 0 0 10px #00b894;
    }

    /* Hover effects */
    .cyber-filter:hover {
        color: #0984e3;
    }

    .cyber-filter:hover .filter-glow {
        opacity: 1;
    }

    .cyber-filter:hover .status-core {
        transform: translate(-50%, -50%) scale(1.2);
    }

    /* Mobile styles */
    .mobile-filter {
        display: none;
        padding: 0.5rem;
        gap: 0.5rem;
        background: #1c1c1c;
        border: 1px solid rgba(9, 132, 227, 0.2);
        border-radius: 4px;
    }

    .filter-select {
        width: 100%;
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #fff;
        padding: 0.5rem;
        border-radius: 3px;
        font-family: inherit;
        font-size: 0.8rem;
        appearance: none;
    }

    /* Media queries */
    @media screen and (min-width: 601px) {
        .filter-group { display: flex; }
        .mobile-filter { display: none; }
    }

    @media screen and (max-width: 600px) {
        .filter-group { display: none; }
        .mobile-filter { display: flex; }
    }

    /* Animation for status core */
    @keyframes pulse {
        0%, 100% { transform: translate(-50%, -50%) scale(1); }
        50% { transform: translate(-50%, -50%) scale(1.2); }
    }

    .cyber-filter.active .status-core {
        animation: pulse 2s infinite;
    }
</style>
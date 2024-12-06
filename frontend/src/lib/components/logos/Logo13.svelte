<!-- src/lib/components/Logo.svelte -->
<script>
    import { goto } from '$app/navigation';
    
    export let size = '32';
    export let clickable = true;

    $: width = size;
    $: height = size;
</script>

<div class="logo-container">
    {#if clickable}
        <a href="/tasks" class="logo-link" on:click|preventDefault={() => goto('/tasks')}>
            <div class="logo">
                <div class="hex-container">
                    <svg {width} {height} viewBox="0 0 24 24" class="hex-svg">
                        <!-- Hexagon background -->
                        <path 
                            class="hex-bg" 
                            d="M12 2L20 7V17L12 22L4 17V7L12 2Z" 
                        />
                        <!-- Terminal line -->
                        <path 
                            class="terminal-line" 
                            d="M8 12L10 14L14 10" 
                            stroke-linecap="round" 
                            stroke-linejoin="round"
                        />
                        <!-- Dots -->
                        <circle class="dot" cx="7" cy="8" r="0.5" />
                        <circle class="dot" cx="17" cy="8" r="0.5" />
                        <circle class="dot" cx="7" cy="16" r="0.5" />
                        <circle class="dot" cx="17" cy="16" r="0.5" />
                    </svg>
                    <div class="cursor">█</div>
                </div>
                <div class="text">
                    <span class="action">Action</span>
                    <span class="hub">Hub</span>
                    <span class="status">[1.0]</span>
                </div>
            </div>
        </a>
    {:else}
        <div class="logo">
            <!-- Same content -->
        </div>
    {/if}
</div>

<style>
    .logo-container {
        display: flex;
        align-items: center;
    }

    .logo {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", monospace;
        position: relative;
    }

    .hex-container {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .hex-svg {
        filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
    }

    .hex-bg {
        fill: #1c1c1c;
        stroke: #2980b9;
        stroke-width: 1.5;
    }

    .terminal-line {
        stroke: #2ecc71;
        stroke-width: 1.5;
        opacity: 0.8;
    }

    .dot {
        fill: #3498db;
        opacity: 0.8;
    }

    .cursor {
        position: absolute;
        color: #2ecc71;
        font-size: 0.8rem;
        animation: blink 1s steps(1) infinite;
    }

    .text {
        display: flex;
        align-items: baseline;
        gap: 0.2rem;
    }

    .action {
        color: #2c3e50;
        font-weight: bold;
        font-size: 1rem;
    }

    .hub {
        color: #2980b9;
        font-weight: bold;
        font-size: 1rem;
    }

    .status {
        color: #7f8c8d;
        font-size: 0.6rem;
        margin-left: 0.2rem;
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    /* Hover effects */
    .logo:hover .hex-bg {
        stroke: #3498db;
        filter: drop-shadow(0 0 2px #3498db);
    }

    .logo:hover .terminal-line {
        stroke: #27ae60;
        animation: pulse 2s infinite;
    }

    .logo:hover .dot {
        animation: glow 1.5s infinite alternate;
    }

    @keyframes pulse {
        0%, 100% { opacity: 0.8; }
        50% { opacity: 1; }
    }

    @keyframes glow {
        from { opacity: 0.8; }
        to { opacity: 1; filter: drop-shadow(0 0 2px #3498db); }
    }

    /* Scan line effect */
    .hex-container::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            transparent 50%,
            rgba(46, 204, 113, 0.1) 50%
        );
        background-size: 100% 4px;
        pointer-events: none;
        opacity: 0.5;
    }

    /* Optional: Add matrix rain effect */
    .hex-container::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(180deg,
            transparent 0%,
            rgba(46, 204, 113, 0.1) 50%,
            transparent 100%
        );
        background-size: 100% 20px;
        animation: matrix 1s linear infinite;
        pointer-events: none;
        opacity: 0.1;
    }

    @keyframes matrix {
        0% { background-position: 0 0; }
        100% { background-position: 0 20px; }
    }
</style>
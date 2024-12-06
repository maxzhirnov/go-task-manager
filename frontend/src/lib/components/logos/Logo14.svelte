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
                <div class="chip-container">
                    <svg {width} {height} viewBox="0 0 24 24" class="chip-svg">
                        <!-- Main circuit board -->
                        <rect 
                            class="board" 
                            x="4" y="4" 
                            width="16" height="16" 
                            rx="2"
                        />
                        
                        <!-- Circuit paths -->
                        <path 
                            class="circuit-lines"
                            d="M4 12h4M16 12h4M12 4v4M12 16v4"
                            stroke-linecap="round"
                        />
                        
                        <!-- CPU core -->
                        <rect 
                            class="cpu-core"
                            x="8" y="8"
                            width="8" height="8"
                            rx="1"
                        />
                        
                        <!-- Connection points -->
                        <circle class="connection" cx="8" cy="8" r="1" />
                        <circle class="connection" cx="16" cy="8" r="1" />
                        <circle class="connection" cx="8" cy="16" r="1" />
                        <circle class="connection" cx="16" cy="16" r="1" />
                        
                        <!-- Data pulse -->
                        <circle class="data-pulse" cx="12" cy="12" r="2" />
                    </svg>
                    <div class="power-indicator"></div>
                </div>
                <div class="text">
                    <span class="action">Action</span>
                    <span class="separator">/</span>
                    <span class="hub">Hub</span>
                    <span class="version">v0.0.1</span>
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
    }

    .chip-container {
        position: relative;
        filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
    }

    .board {
        fill: #2d3436;
        stroke: #0984e3;
        stroke-width: 1;
    }

    .circuit-lines {
        stroke: #00b894;
        stroke-width: 1;
        opacity: 0.8;
    }

    .cpu-core {
        fill: #2d3436;
        stroke: #00b894;
        stroke-width: 1;
    }

    .connection {
        fill: #0984e3;
        filter: drop-shadow(0 0 2px rgba(9, 132, 227, 0.5));
    }

    .data-pulse {
        fill: #00b894;
        opacity: 0.5;
        animation: pulse 2s ease-in-out infinite;
    }

    .power-indicator {
        position: absolute;
        top: 25%;
        left: 25%;
        width: 50%;
        height: 50%;
        background: radial-gradient(
            circle at center,
            rgba(0, 184, 148, 0.2) 0%,
            transparent 70%
        );
        pointer-events: none;
        animation: glow 3s ease-in-out infinite;
    }

    .text {
        display: flex;
        align-items: baseline;
        gap: 0.2rem;
    }

    .action {
        color: #2d3436;
        font-weight: bold;
        font-size: 1rem;
    }

    .separator {
        color: #00b894;
        font-weight: normal;
        font-size: 1rem;
    }

    .hub {
        color: #0984e3;
        font-weight: bold;
        font-size: 1rem;
    }

    .version {
        color: #636e72;
        font-size: 0.6rem;
        margin-left: 0.3rem;
    }

    @keyframes pulse {
        0%, 100% { 
            opacity: 0.5;
            transform: scale(1);
        }
        50% { 
            opacity: 0.8;
            transform: scale(1.2);
        }
    }

    @keyframes glow {
        0%, 100% {
            opacity: 0.6;
            transform: scale(1);
        }
        50% {
            opacity: 0.8;
            transform: scale(1.1);
        }
    }

    /* Hover effects */
    .logo:hover .board {
        stroke: #00b894;
    }

    .logo:hover .circuit-lines {
        animation: circuit 2s linear infinite;
    }

    .logo:hover .connection {
        animation: blink 1.5s infinite alternate;
    }

    @keyframes circuit {
        0% {
            stroke-dasharray: 0 12;
            stroke-dashoffset: 0;
        }
        100% {
            stroke-dasharray: 12 12;
            stroke-dashoffset: 24;
        }
    }

    @keyframes blink {
        from { opacity: 0.8; }
        to { opacity: 1; filter: drop-shadow(0 0 3px #0984e3); }
    }

    /* Circuit pattern overlay */
    .chip-container::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.1) 1px,
                transparent 1px
            );
        background-size: 4px 4px;
        pointer-events: none;
        opacity: 0.5;
    }
</style>
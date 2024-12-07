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
                    <!-- Main hexagon -->
                    <svg class="hex-frame" viewBox="0 0 100 100">
                        <polygon 
                            points="50,5 95,27.5 95,72.5 50,95 5,72.5 5,27.5" 
                            class="hex-outline"
                        />
                        <!-- Inner circuit lines -->
                        <path d="M20,50 H40" class="circuit-line" />
                        <path d="M60,50 H80" class="circuit-line" />
                        <path d="M50,20 V40" class="circuit-line" />
                        <path d="M50,60 V80" class="circuit-line" />
                        <!-- Connection nodes -->
                        <circle cx="50" cy="50" r="15" class="core-circle" />
                        <circle cx="50" cy="50" r="3" class="core-dot" />
                        <circle cx="20" cy="50" r="2" class="node" />
                        <circle cx="80" cy="50" r="2" class="node" />
                        <circle cx="50" cy="20" r="2" class="node" />
                        <circle cx="50" cy="80" r="2" class="node" />
                    </svg>
                    <!-- Core text -->
                    <div class="core-text">A::H</div>
                </div>
                <div class="text-container">
                    <div class="logo-text">
                        <span class="action">ACTION</span>
                        <span class="separator">/</span>
                        <span class="hub">HUB</span>
                    </div>
                    <div class="status-line">[SYS:ONLINE][v0.0.1]</div>
                </div>
            </div>
        </a>
    {:else}
        <div class="logo">
            <!-- Same content without click handler -->
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
        gap: 1rem;
        font-family: "JetBrains Mono", monospace;
    }

    .hex-container {
        position: relative;
        width: 40px;
        height: 40px;
    }

    .hex-frame {
        width: 100%;
        height: 100%;
    }

    .hex-outline {
        fill: none;
        stroke: #0984e3;
        stroke-width: 1;
    }

    .circuit-line {
        stroke: #00b894;
        stroke-width: 0.5;
        opacity: 0.6;
    }

    .core-circle {
        fill: none;
        stroke: #0984e3;
        stroke-width: 0.5;
        opacity: 0.3;
    }

    .core-dot {
        fill: #00b894;
        filter: drop-shadow(0 0 2px #00b894);
    }

    .node {
        fill: #0984e3;
        opacity: 0.8;
    }

    .core-text {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: #00b894;
        font-size: 0.6rem;
        font-weight: bold;
        white-space: nowrap;
    }

    .text-container {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .logo-text {
        font-size: 1rem;
        font-weight: bold;
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .action { color: #0984e3; }
    .separator { color: #00b894; }
    .hub { color: #0984e3; }

    .status-line {
        font-size: 0.5rem;
        color: #00b894;
        opacity: 0.8;
    }

    /* Animations */
    .core-dot {
        animation: pulse 2s infinite;
    }

    .circuit-line {
        stroke-dasharray: 20;
        animation: circuit 4s linear infinite;
    }

    .node {
        animation: nodePulse 3s infinite alternate;
    }

    @keyframes pulse {
        0%, 100% { 
            transform: scale(1);
            opacity: 1;
        }
        50% { 
            transform: scale(1.2);
            opacity: 0.8;
        }
    }

    @keyframes circuit {
        to {
            stroke-dashoffset: 40;
        }
    }

    @keyframes nodePulse {
        to {
            opacity: 0.4;
        }
    }

    /* Hover effects */
    .logo:hover .hex-outline {
        filter: drop-shadow(0 0 2px #0984e3);
    }

    .logo:hover .core-dot {
        filter: drop-shadow(0 0 4px #00b894);
    }

    .logo:hover .status-line {
        animation: glitch 0.3s ease infinite;
    }

    @keyframes glitch {
        0%, 100% { transform: translate(0); }
        20% { transform: translate(-1px, 1px); }
        40% { transform: translate(-1px, -1px); }
        60% { transform: translate(1px, 1px); }
    }
</style>
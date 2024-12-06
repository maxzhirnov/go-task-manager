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
                <div class="screen-border">
                    <div class="screen">
                        <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <!-- Pixel art elements -->
                            <rect x="4" y="8" width="2" height="2" class="pixel" />
                            <rect x="6" y="6" width="2" height="2" class="pixel" />
                            <rect x="8" y="4" width="2" height="2" class="pixel" />
                            <rect x="10" y="6" width="2" height="2" class="pixel" />
                            <rect x="12" y="8" width="2" height="2" class="pixel" />
                            <rect x="14" y="10" width="2" height="2" class="pixel" />
                            <rect x="16" y="12" width="2" height="2" class="pixel" />
                            
                            <!-- Loading bar -->
                            <rect x="4" y="16" width="16" height="2" class="loading-bar" />
                        </svg>
                        <div class="scanlines"></div>
                        <div class="screen-glare"></div>
                    </div>
                </div>
                <div class="text-container">
                    <span class="logo-text">
                        <span class="action blink">>_</span>
                        <span class="action">ACTION</span>
                        <span class="hub">[HUB]</span>
                    </span>
                    <span class="subtext">v1.0 © 1988</span>
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
        gap: 1rem;
        font-family: "Press Start 2P", "Courier New", monospace;
    }

    .screen-border {
        background: #2C3E50;
        padding: 4px;
        border-radius: 4px;
        box-shadow: 
            inset -2px -2px 0 #1a252f,
            inset 2px 2px 0 #34495e,
            0 2px 4px rgba(0,0,0,0.1);
    }

    .screen {
        background: #1a252f;
        position: relative;
        overflow: hidden;
        padding: 4px;
        border: 1px solid #34495e;
    }

    .pixel {
        fill: #3498db;
        animation: pixelate 0.5s steps(2) infinite;
    }

    .loading-bar {
        fill: #3498db;
        animation: loading 2s steps(16) infinite;
    }

    .scanlines {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            transparent 50%,
            rgba(52, 152, 219, 0.1) 50%
        );
        background-size: 100% 4px;
        animation: scroll 8s linear infinite;
        pointer-events: none;
    }

    .screen-glare {
        position: absolute;
        top: -50%;
        left: -50%;
        width: 200%;
        height: 200%;
        background: radial-gradient(
            circle at center,
            rgba(255, 255, 255, 0.1) 0%,
            transparent 70%
        );
        pointer-events: none;
    }

    .text-container {
        display: flex;
        flex-direction: column;
    }

    .logo-text {
        font-size: 1.2rem;
        color: #2C3E50;
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .action {
        color: #2C3E50;
        text-shadow: 1px 1px 0 rgba(44, 62, 80, 0.1);
    }

    .blink {
        animation: blink 1s steps(2) infinite;
        color: #3498db;
    }

    .hub {
        color: #e74c3c;
        text-shadow: 1px 1px 0 rgba(231, 76, 60, 0.1);
    }

    .subtext {
        font-size: 0.6rem;
        color: #7f8c8d;
        margin-top: 0.25rem;
    }

    @keyframes pixelate {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.8; }
    }

    @keyframes loading {
        from { width: 0; }
        to { width: 16px; }
    }

    @keyframes scroll {
        from { transform: translateY(0); }
        to { transform: translateY(4px); }
    }

    @keyframes blink {
        0%, 100% { opacity: 1; }
        50% { opacity: 0; }
    }

    /* Hover effects */
    .logo:hover .screen-border {
        box-shadow: 
            inset -2px -2px 0 #1a252f,
            inset 2px 2px 0 #34495e,
            0 4px 8px rgba(0,0,0,0.15);
    }

    .logo:hover .pixel,
    .logo:hover .loading-bar {
        fill: #2980b9;
    }

    .logo:hover .action {
        color: #234257;
    }

    .logo:hover .hub {
        color: #c0392b;
    }

    /* Static noise effect with reduced opacity */
    .screen::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyBAMAAADsEZWCAAAAGFBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgE1xQAAAABh0Uk5TAAECAwQFBgcICQoLDA0ODxAREhMUFRYX0UzQnAAAAJlJREFUGJVjYGBgEDIxMWE4wYAGmBhEhCRYTzCgAQkWoXnzpmAIYDEvmjdvPoYAi+k85ni7oQmwh8+bNzESQ4A93HTeDBQ+ewZUPruVA1SAw4FFgcUCQ4DDgUWBxRxDgMOBRYHFHEOAw4FFgcUCQ4DDgUWBxRxDgMOBRYHFAkOAw4FFgcUcQ4DDgUWBxQJDgMOBRYHFHEMAAI4UFEFh7bKvAAAAAElFTkSuQmCC");
        opacity: 0.05;
        pointer-events: none;
    }
</style>
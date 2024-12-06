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
        background: #888;
        padding: 4px;
        border-radius: 2px;
        box-shadow: 
            inset -2px -2px 0 #444,
            inset 2px 2px 0 #ccc;
    }

    .screen {
        background: #000;
        position: relative;
        overflow: hidden;
        padding: 4px;
    }

    .pixel {
        fill: #00FF00;
        animation: pixelate 0.5s steps(2) infinite;
    }

    .loading-bar {
        fill: #00FF00;
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
            rgba(0, 0, 0, 0.5) 50%
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
        color: #00FF00;
        text-shadow: 0 0 5px rgba(0, 255, 0, 0.5);
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .action {
        color: #00FF00;
    }

    .blink {
        animation: blink 1s steps(2) infinite;
    }

    .hub {
        color: #FFFF00;
    }

    .subtext {
        font-size: 0.6rem;
        color: #888;
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

    /* CRT power off effect on hover */
    .logo:hover .screen {
        animation: poweroff 0.15s linear;
    }

    @keyframes poweroff {
        0% { transform: scale(1, 1); }
        50% { transform: scale(1, 0.005); }
        100% { transform: scale(0, 0); }
    }

    /* Static noise effect */
    .screen::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyBAMAAADsEZWCAAAAGFBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgE1xQAAAABh0Uk5TAAECAwQFBgcICQoLDA0ODxAREhMUFRYX0UzQnAAAAJlJREFUGJVjYGBgEDIxMWE4wYAGmBhEhCRYTzCgAQkWoXnzpmAIYDEvmjdvPoYAi+k85ni7oQmwh8+bNzESQ4A93HTeDBQ+ewZUPruVA1SAw4FFgcUCQ4DDgUWBxRxDgMOBRYHFHEOAw4FFgcUCQ4DDgUWBxRxDgMOBRYHFAkOAw4FFgcUcQ4DDgUWBxQJDgMOBRYHFHEMAAI4UFEFh7bKvAAAAAElFTkSuQmCC");
        opacity: 0.1;
        pointer-events: none;
    }
</style>
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
                <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <!-- Grid background -->
                    <path class="grid" d="M0 6h24M0 12h24M0 18h24" stroke="#2C2C2C" stroke-width="0.5" />
                    <path class="grid" d="M6 0v24M12 0v24M18 0v24" stroke="#2C2C2C" stroke-width="0.5" />
                    
                    <!-- Main "screen" rectangle -->
                    <rect x="4" y="4" width="16" height="16" class="screen" />
                    
                    <!-- Pixel art "A" -->
                    <rect x="8" y="8" width="2" height="2" class="pixel" />
                    <rect x="10" y="8" width="2" height="2" class="pixel" />
                    <rect x="12" y="8" width="2" height="2" class="pixel" />
                    <rect x="8" y="10" width="2" height="2" class="pixel" />
                    <rect x="12" y="10" width="2" height="2" class="pixel" />
                    <rect x="8" y="12" width="2" height="2" class="pixel" />
                    <rect x="10" y="12" width="2" height="2" class="pixel" />
                    <rect x="12" y="12" width="2" height="2" class="pixel" />
                    <rect x="8" y="14" width="2" height="2" class="pixel" />
                    <rect x="12" y="14" width="2" height="2" class="pixel" />
                    
                    <!-- Decorative scanline -->
                    <rect x="4" y="11" width="16" height="0.5" class="scanline" />
                    
                    <!-- Corner accents -->
                    <path d="M2 2h4v4" class="corner" />
                    <path d="M22 2h-4v4" class="corner" />
                    <path d="M2 22h4v-4" class="corner" />
                    <path d="M22 22h-4v-4" class="corner" />
                </svg>
                <span class="logo-text">Action<span class="hub">HUB</span></span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <!-- Same SVG content as above -->
        </div>
    {/if}
</div>

<style>
    .logo-container {
        display: flex;
        align-items: center;
    }

    .logo-link {
        text-decoration: none;
        color: inherit;
    }

    .logo {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .grid {
        opacity: 0.2;
    }

    .screen {
        fill: #000;
        stroke: #00FF95;
        stroke-width: 1;
        filter: drop-shadow(0 0 2px #00FF95);
    }

    .pixel {
        fill: #00FF95;
        animation: flicker 4s infinite;
    }

    .scanline {
        fill: #00FF95;
        opacity: 0.5;
        animation: scan 2s linear infinite;
    }

    .corner {
        stroke: #00FF95;
        stroke-width: 1;
        fill: none;
    }

    .logo-text {
        font-family: "VT323", "Courier New", monospace;
        font-size: 1.75rem;
        color: #00FF95;
        text-shadow: 0 0 10px rgba(0, 255, 149, 0.5);
        letter-spacing: 1px;
    }

    .hub {
        background: #00FF95;
        color: #000;
        padding: 0 4px;
        margin-left: 2px;
    }

    /* Hover effects */
    .logo:hover .screen {
        stroke: #00FFA3;
        filter: drop-shadow(0 0 4px #00FF95);
    }

    .logo:hover .pixel {
        animation: pixelate 0.5s steps(2) infinite;
    }

    @keyframes scan {
        0% {
            transform: translateY(-7px);
            opacity: 0;
        }
        50% {
            opacity: 0.5;
        }
        100% {
            transform: translateY(7px);
            opacity: 0;
        }
    }

    @keyframes flicker {
        0%, 100% { opacity: 1; }
        96%, 98% { opacity: 0.8; }
        97% { opacity: 0.5; }
    }

    @keyframes pixelate {
        0%, 100% { transform: translateX(0); }
        50% { transform: translateX(1px); }
    }

    /* CRT screen effect */
    .logo::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            rgba(18, 16, 16, 0) 50%,
            rgba(0, 0, 0, 0.25) 50%
        );
        background-size: 100% 4px;
        z-index: 2;
        opacity: 0.1;
        pointer-events: none;
    }

    /* Retro screen glow */
    .logo:hover .logo-text {
        animation: textGlow 2s infinite;
    }

    @keyframes textGlow {
        0%, 100% { text-shadow: 0 0 10px rgba(0, 255, 149, 0.5); }
        50% { text-shadow: 0 0 20px rgba(0, 255, 149, 0.8); }
    }
</style>
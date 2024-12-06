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
                <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <!-- Sunset gradient background -->
                    <rect x="0" y="0" width="24" height="24" class="sunset-bg"/>
                    
                    <!-- Palm tree silhouette -->
                    <path 
                        d="M12 20V14M10 14c0-2 2-3 2-3s2 1 2 3M8 12c0-2 4-2 4-2s4 0 4 2" 
                        class="palm-tree"
                    />
                    <path 
                        d="M7 10c1-1 5-1 5-1s4 0 5 1" 
                        class="palm-leaves"
                    />
                    
                    <!-- Grid lines -->
                    <path 
                        d="M0 16h24M0 18h24M0 20h24" 
                        class="grid-lines"
                    />
                    
                    <!-- Neon frame -->
                    <rect 
                        x="2" y="2" 
                        width="20" height="20" 
                        class="neon-frame"
                    />
                </svg>
                <span class="logo-text">
                    <span class="action">Action</span>
                    <span class="hub">HUB</span>
                </span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <!-- Same SVG content -->
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
        position: relative;
    }

    .sunset-bg {
        fill: url(#sunset);
    }

    .palm-tree {
        stroke: #000;
        stroke-width: 1.5;
        filter: drop-shadow(0 0 2px rgba(0, 0, 0, 0.5));
    }

    .palm-leaves {
        stroke: #000;
        stroke-width: 1.5;
        filter: drop-shadow(0 0 2px rgba(0, 0, 0, 0.5));
    }

    .grid-lines {
        stroke: rgba(255, 255, 255, 0.3);
        stroke-width: 0.5;
    }

    .neon-frame {
        fill: none;
        stroke: #FF00FF;
        stroke-width: 0.5;
        filter: drop-shadow(0 0 4px #FF00FF);
    }

    .logo-text {
        font-family: "Brush Script MT", cursive;
        font-size: 1.75rem;
        letter-spacing: 1px;
        filter: drop-shadow(0 0 2px rgba(255, 255, 255, 0.5));
    }

    .action {
        background: linear-gradient(to bottom, #FF69B4, #FF00FF);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        font-weight: bold;
    }

    .hub {
        color: #1ed3d3;
        font-weight: bold;
        text-shadow: 0 0 10px #1baaaa;
        margin-left: 4px;
    }

    /* Hover animations */
    .logo:hover .neon-frame {
        animation: neonPulse 2s infinite;
    }

    .logo:hover .palm-tree,
    .logo:hover .palm-leaves {
        animation: sway 4s ease-in-out infinite;
    }

    @keyframes neonPulse {
        0%, 100% {
            filter: drop-shadow(0 0 4px #FF00FF);
        }
        50% {
            filter: drop-shadow(0 0 8px #FF00FF);
        }
    }

    @keyframes sway {
        0%, 100% {
            transform: rotate(0deg);
        }
        50% {
            transform: rotate(5deg);
        }
    }

    /* Add sunset gradient definition */
    .logo svg {
        --sunset-start: #FF6B6B;
        --sunset-middle: #FF69B4;
        --sunset-end: #4B0082;
        background: linear-gradient(
            to bottom,
            var(--sunset-start),
            var(--sunset-middle) 50%,
            var(--sunset-end)
        );
    }

    /* Retro scan lines effect */
    .logo::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: repeating-linear-gradient(
            0deg,
            rgba(0, 0, 0, 0.1),
            rgba(0, 0, 0, 0.1) 1px,
            transparent 1px,
            transparent 2px
        );
        pointer-events: none;
    }

    /* Chrome effect for text */
    .logo:hover .action {
        animation: chromeEffect 3s infinite;
    }

    @keyframes chromeEffect {
        0%, 100% {
            filter: brightness(100%);
        }
        50% {
            filter: brightness(150%) saturate(150%);
        }
    }
</style>
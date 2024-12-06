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
                    <!-- Main hexagon -->
                    <path 
                        d="M12 2L20.5 7V17L12 22L3.5 17V7L12 2Z" 
                        stroke="#333" 
                        stroke-width="2" 
                        class="hexagon"
                    />
                    <!-- Inner details -->
                    <path 
                        d="M12 6L16.5 8.5V13.5L12 16L7.5 13.5V8.5L12 6Z" 
                        stroke="#666" 
                        stroke-width="1.5" 
                        class="inner-hex"
                    />
                    <!-- Accent dot -->
                    <circle 
                        cx="12" 
                        cy="11" 
                        r="2" 
                        class="accent-dot"
                    />
                    <!-- Pulse rings -->
                    <circle 
                        cx="12" 
                        cy="11" 
                        r="4" 
                        class="pulse-ring"
                    />
                </svg>
                <span class="logo-text">ActionHub</span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path 
                    d="M12 2L20.5 7V17L12 22L3.5 17V7L12 2Z" 
                    stroke="#333" 
                    stroke-width="2" 
                    class="hexagon"
                />
                <path 
                    d="M12 6L16.5 8.5V13.5L12 16L7.5 13.5V8.5L12 6Z" 
                    stroke="#666" 
                    stroke-width="1.5" 
                    class="inner-hex"
                />
                <circle 
                    cx="12" 
                    cy="11" 
                    r="2" 
                    class="accent-dot"
                />
                <circle 
                    cx="12" 
                    cy="11" 
                    r="4" 
                    class="pulse-ring"
                />
            </svg>
            <span class="logo-text">ActionHub</span>
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
        gap: 0.5rem;
    }

    .logo svg {
        transition: all 0.3s ease;
    }

    .hexagon {
        fill: none;
        stroke: #333;
        transition: all 0.3s ease;
    }

    .inner-hex {
        fill: none;
        stroke: #666;
        transition: all 0.3s ease;
    }

    .accent-dot {
        fill: #FF4081;
        filter: drop-shadow(0 0 2px #FF4081);
    }

    .pulse-ring {
        fill: none;
        stroke: #FF4081;
        stroke-width: 0.5;
        stroke-opacity: 0;
    }

    .logo:hover svg {
        transform: scale(1.1);
    }

    .logo:hover .hexagon {
        stroke: #444;
        filter: drop-shadow(0 0 2px #333);
    }

    .logo:hover .inner-hex {
        stroke: #888;
        transform-origin: center;
        animation: rotate 4s linear infinite;
    }

    .logo:hover .accent-dot {
        animation: pulse 2s ease-in-out infinite;
    }

    .logo:hover .pulse-ring {
        animation: ring-pulse 2s ease-out infinite;
    }

    .logo-text {
        font-size: 1.5rem;
        font-weight: bold;
        background: linear-gradient(135deg, #333 0%, #666 50%, #333 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: 0.5px;
        text-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    @keyframes rotate {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    @keyframes pulse {
        0% {
            transform: scale(1);
            fill: #FF4081;
        }
        50% {
            transform: scale(1.2);
            fill: #FF80AB;
        }
        100% {
            transform: scale(1);
            fill: #FF4081;
        }
    }

    @keyframes ring-pulse {
        0% {
            stroke-opacity: 0.5;
            stroke-width: 0.5;
            transform: scale(0.8);
        }
        100% {
            stroke-opacity: 0;
            stroke-width: 2;
            transform: scale(1.5);
        }
    }
</style>
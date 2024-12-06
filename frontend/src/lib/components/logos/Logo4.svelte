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
                    <!-- Background circle -->
                    <circle cx="12" cy="12" r="11" class="bg-circle" />
                    
                    <!-- Geometric shapes -->
                    <rect x="6" y="6" width="6" height="6" class="square" />
                    <circle cx="15" cy="15" r="3" class="small-circle" />
                    <path d="M18 6L6 18" class="diagonal-line" />
                    
                    <!-- Decorative elements -->
                    <path d="M3 12h4M17 12h4" class="side-lines" />
                    <circle cx="12" cy="12" r="1.5" class="center-dot" />
                </svg>
                <span class="logo-text">ActionHub</span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="11" class="bg-circle" />
                <rect x="6" y="6" width="6" height="6" class="square" />
                <circle cx="15" cy="15" r="3" class="small-circle" />
                <path d="M18 6L6 18" class="diagonal-line" />
                <path d="M3 12h4M17 12h4" class="side-lines" />
                <circle cx="12" cy="12" r="1.5" class="center-dot" />
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
        gap: 0.75rem;
    }

    .bg-circle {
        fill: #F5F5F5;
        stroke: #E53935;
        stroke-width: 1.5;
    }

    .square {
        fill: #1E88E5;
        transform-origin: center;
        transition: transform 0.3s ease;
    }

    .small-circle {
        fill: #FDD835;
        transform-origin: center;
        transition: transform 0.3s ease;
    }

    .diagonal-line {
        stroke: #000;
        stroke-width: 1.5;
        stroke-linecap: round;
    }

    .side-lines {
        stroke: #000;
        stroke-width: 1.5;
        stroke-linecap: round;
        opacity: 0.7;
    }

    .center-dot {
        fill: #000;
    }

    .logo-text {
        font-family: 'Futura', 'Arial', sans-serif;
        font-size: 1.5rem;
        font-weight: bold;
        letter-spacing: 1px;
        color: #000;
        text-transform: uppercase;
    }

    /* Hover animations */
    .logo:hover .square {
        transform: rotate(45deg);
    }

    .logo:hover .small-circle {
        transform: scale(1.2);
    }

    .logo:hover .bg-circle {
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0% {
            stroke-width: 1.5;
        }
        50% {
            stroke-width: 2.5;
        }
        100% {
            stroke-width: 1.5;
        }
    }

    /* Bauhaus-inspired color variations on hover */
    .logo:hover .square {
        fill: #1565C0;
    }

    .logo:hover .small-circle {
        fill: #FBC02D;
    }

    .logo:hover .center-dot {
        animation: blink 2s infinite;
    }

    @keyframes blink {
        0%, 100% {
            transform: scale(1);
        }
        50% {
            transform: scale(1.5);
        }
    }

    /* Responsive adjustments */
    @media (max-width: 768px) {
        .logo-text {
            font-size: 1.25rem;
        }
    }
</style>
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
                    <!-- Arrow/Smile combination -->
                    <path 
                        class="smile-arrow"
                        d="M4 12C4 12 7 15 12 15C17 15 20 12 20 12" 
                        stroke-linecap="round"
                    />
                    <!-- Arrow head -->
                    <path 
                        class="arrow-head"
                        d="M16 9L20 12L16 15" 
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
                <span class="logo-text">
                    action<span class="hub">hub</span>
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
        gap: 0.5rem;
    }

    .smile-arrow {
        stroke: #FF9900;
        stroke-width: 2;
        transition: all 0.3s ease;
    }

    .arrow-head {
        stroke: #FF9900;
        stroke-width: 2;
        transition: all 0.3s ease;
    }

    .logo-text {
        font-family: Arial, sans-serif;
        font-size: 1.5rem;
        font-weight: bold;
        color: #232F3E;
        letter-spacing: -0.5px;
    }

    .hub {
        color: #FF9900;
    }

    /* Hover animations */
    .logo:hover .smile-arrow {
        stroke-width: 2.5;
        filter: drop-shadow(0 0 2px rgba(255, 153, 0, 0.5));
    }

    .logo:hover .arrow-head {
        stroke-width: 2.5;
        filter: drop-shadow(0 0 2px rgba(255, 153, 0, 0.5));
        animation: pointRight 1s ease-in-out infinite;
    }

    @keyframes pointRight {
        0%, 100% {
            transform: translateX(0);
        }
        50% {
            transform: translateX(2px);
        }
    }

    /* Mobile optimization */
    @media (max-width: 768px) {
        .logo-text {
            font-size: 1.25rem;
        }
    }
</style>
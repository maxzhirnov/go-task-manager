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
                    <!-- Outer circle -->
                    <circle cx="12" cy="12" r="10" stroke="#4CAF50" stroke-width="2" class="circle"/>
                    <!-- Arrow -->
                    <path 
                        d="M8 12L11 15L16 9" 
                        stroke="#4CAF50" 
                        stroke-width="2" 
                        stroke-linecap="round" 
                        stroke-linejoin="round"
                        class="checkmark"
                    />
                    <!-- Action lines -->
                    <path 
                        d="M12 2V4M12 20V22M2 12H4M20 12h2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" 
                        stroke="#45a049" 
                        stroke-width="2" 
                        stroke-linecap="round"
                        class="rays"
                    />
                </svg>
                <span class="logo-text">ActionHub</span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="10" stroke="#4CAF50" stroke-width="2" class="circle"/>
                <path 
                    d="M8 12L11 15L16 9" 
                    stroke="#4CAF50" 
                    stroke-width="2" 
                    stroke-linecap="round" 
                    stroke-linejoin="round"
                    class="checkmark"
                />
                <path 
                    d="M12 2V4M12 20V22M2 12H4M20 12h2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" 
                    stroke="#45a049" 
                    stroke-width="2" 
                    stroke-linecap="round"
                    class="rays"
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
        position: relative;
    }

    .logo svg {
        transition: all 0.3s ease;
    }

    .logo:hover svg {
        transform: scale(1.1);
    }

    .logo:hover .circle {
        stroke: #45a049;
        filter: drop-shadow(0 0 2px #4CAF50);
    }

    .logo:hover .checkmark {
        stroke: #4CAF50;
        stroke-dasharray: 20;
        animation: dash 1.5s ease-in-out infinite;
    }

    .logo:hover .rays {
        animation: spin 4s linear infinite;
    }

    .logo-text {
        font-size: 1.5rem;
        font-weight: bold;
        background: linear-gradient(135deg, #4CAF50 0%, #45a049 50%, #357a38 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: 0.5px;
        text-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    @keyframes dash {
        0% {
            stroke-dashoffset: 20;
        }
        50% {
            stroke-dashoffset: 0;
        }
        100% {
            stroke-dashoffset: 20;
        }
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }
</style>
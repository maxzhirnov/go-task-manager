<script>
    import { goto } from '$app/navigation';
    
    // Optional props for customization
    export let size = '32'; // SVG size
    export let clickable = true; // Whether logo should navigate on click

    // Set width and height from size
    $: width = size;
    $: height = size;
</script>

<div class="logo-container">
    {#if clickable}
        <a href="/tasks" class="logo-link" on:click|preventDefault={() => goto('/tasks')}>
            <div class="logo">
                <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M2 17L12 22L22 17" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M2 12L12 17L22 12" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span class="logo-text">ActionHub</span>
            </div>
        </a>
    {:else}
        <div class="logo">
            <svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M2 17L12 22L22 17" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M2 12L12 17L22 12" stroke="#4CAF50" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
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
        transition: transform 0.3s ease;
    }

    .logo:hover svg {
        transform: scale(1.1);
    }

    .logo-text {
        font-size: 1.5rem;
        font-weight: bold;
        background: linear-gradient(135deg, #4CAF50, #45a049);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: 0.5px;
    }
</style>
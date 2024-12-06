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
                    <!-- DNA Helix background -->
                    <path class="dna-strand" d="M6 4c4 3 8 3 12 0" />
                    <path class="dna-strand" d="M6 8c4 3 8 3 12 0" />
                    <path class="dna-strand" d="M6 12c4 3 8 3 12 0" />
                    <path class="dna-strand" d="M6 16c4 3 8 3 12 0" />
                    <path class="dna-strand" d="M6 20c4 3 8 3 12 0" />
                    
                    <!-- Circuit paths -->
                    <path class="circuit" d="M2 12h4m12 0h4M12 2v4m0 12v4" />
                    
                    <!-- Central hexagon -->
                    <path class="hex-outer" d="M12 4L20 8L20 16L12 20L4 16L4 8Z" />
                    <path class="hex-inner" d="M12 8L16 10L16 14L12 16L8 14L8 10Z" />
                    
                    <!-- Animated particles -->
                    <circle class="particle" cx="12" cy="12" r="1" />
                    <circle class="particle" cx="14" cy="14" r="0.5" />
                    <circle class="particle" cx="10" cy="10" r="0.5" />
                    
                    <!-- Energy rings -->
                    <circle class="energy-ring" cx="12" cy="12" r="3" />
                    <circle class="energy-ring" cx="12" cy="12" r="5" />
                </svg>
                <span class="logo-text">
                    <span class="action">Action</span>
                    <span class="hub">Hub</span>
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

    .logo {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        position: relative;
    }

    .dna-strand {
        stroke: #6366F1;
        stroke-width: 0.5;
        opacity: 0.3;
        filter: blur(0.5px);
    }

    .circuit {
        stroke: #818CF8;
        stroke-width: 0.5;
        stroke-dasharray: 4;
        animation: dash 20s linear infinite;
    }

    .hex-outer {
        fill: none;
        stroke: #4F46E5;
        stroke-width: 1.5;
        filter: drop-shadow(0 0 3px #4F46E5);
    }

    .hex-inner {
        fill: none;
        stroke: #818CF8;
        stroke-width: 1;
        animation: rotate 10s linear infinite;
    }

    .particle {
        fill: #F472B6;
        filter: drop-shadow(0 0 2px #F472B6);
        animation: float 3s ease-in-out infinite;
    }

    .energy-ring {
        fill: none;
        stroke: #A78BFA;
        stroke-width: 0.5;
        opacity: 0.5;
        animation: pulse 4s ease-in-out infinite;
    }

    .logo-text {
        font-family: 'Inter', sans-serif;
        font-size: 1.5rem;
        font-weight: bold;
    }

    .action {
        background: linear-gradient(135deg, #4F46E5, #818CF8);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        position: relative;
    }

    .hub {
        color: #F472B6;
        position: relative;
        margin-left: 2px;
    }

    .hub::after {
        content: '';
        position: absolute;
        bottom: -2px;
        left: 0;
        width: 100%;
        height: 2px;
        background: #F472B6;
        transform: scaleX(0);
        transition: transform 0.3s ease;
    }

    .logo:hover .hub::after {
        transform: scaleX(1);
    }

    @keyframes dash {
        to {
            stroke-dashoffset: 100;
        }
    }

    @keyframes rotate {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    @keyframes float {
        0%, 100% {
            transform: translate(0, 0);
        }
        50% {
            transform: translate(2px, -2px);
        }
    }

    @keyframes pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 0.5;
        }
        50% {
            transform: scale(1.2);
            opacity: 0.2;
        }
    }

    /* Glitch effect on hover */
    .logo:hover .action {
        animation: glitch 0.5s ease-in-out infinite;
    }

    @keyframes glitch {
        0%, 100% {
            transform: translate(0);
        }
        20% {
            transform: translate(-1px, 1px);
        }
        40% {
            transform: translate(-1px, -1px);
        }
        60% {
            transform: translate(1px, 1px);
        }
        80% {
            transform: translate(1px, -1px);
        }
    }

    /* Quantum field effect */
    .logo::before {
        content: '';
        position: absolute;
        top: -5px;
        left: -5px;
        right: -5px;
        bottom: -5px;
        background: radial-gradient(
            circle at 50% 50%,
            rgba(99, 102, 241, 0.1),
            transparent 70%
        );
        opacity: 0;
        transition: opacity 0.3s ease;
        pointer-events: none;
    }

    .logo:hover::before {
        opacity: 1;
    }
</style>
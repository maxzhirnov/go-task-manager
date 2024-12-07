<!-- LoadingSpinner.svelte -->
<script>
    export let size = '40';
    export let message = 'LOADING';
</script>

<div class="loading-container">
    <div class="spinner-wrapper">
        <!-- Hexagonal Frame -->
        <svg class="spinner-hex" viewBox="0 0 100 100" width={size} height={size}>
            <polygon 
                class="hex-frame"
                points="50 3, 93.5 25, 93.5 75, 50 97, 6.5 75, 6.5 25"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            />
        </svg>

        <!-- Circular Spinner -->
        <svg class="spinner-circle" viewBox="0 0 100 100" width={size} height={size}>
            <circle
                class="progress-ring"
                cx="50"
                cy="50"
                r="45"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-dasharray="283"
                stroke-linecap="round"
            />
        </svg>

        <!-- Inner Elements -->
        <svg class="spinner-inner" viewBox="0 0 100 100" width={size} height={size}>
            <!-- Inner segments -->
            {#each Array(6) as _, i}
                <line
                    class="segment"
                    x1="50"
                    y1="15"
                    x2="50"
                    y2="25"
                    stroke="currentColor"
                    stroke-width="2"
                    transform="rotate({i * 60} 50 50)"
                />
            {/each}
        </svg>

        <!-- Center Dot -->
        <div class="center-dot"></div>
    </div>
    
    <!-- Loading Text -->
    <div class="loading-text">
        <span class="text">{message}</span>
        <span class="dots">
            <span class="dot">.</span>
            <span class="dot">.</span>
            <span class="dot">.</span>
        </span>
    </div>
</div>

<style>
    .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
        min-height: 200px;
        color: #0984e3;
    }

    .spinner-wrapper {
        position: relative;
        width: var(--size, 40px);
        height: var(--size, 40px);
    }

    .spinner-hex,
    .spinner-circle,
    .spinner-inner {
        position: absolute;
        top: 0;
        left: 0;
    }

    .hex-frame {
        opacity: 0.5;
    }

    .progress-ring {
        stroke-dashoffset: 283;
        animation: rotate 2s linear infinite;
    }

    .spinner-inner {
        animation: rotate-reverse 4s linear infinite;
    }

    .segment {
        opacity: 0.7;
    }

    .center-dot {
        position: absolute;
        top: 50%;
        left: 50%;
        width: 4px;
        height: 4px;
        background: currentColor;
        border-radius: 50%;
        transform: translate(-50%, -50%);
        animation: pulse 1.5s ease-in-out infinite;
    }

    .loading-text {
        font-family: 'JetBrains Mono', monospace;
        font-size: 0.8rem;
        display: flex;
        align-items: center;
        gap: 0.2rem;
        color: #0984e3;
    }

    .dots {
        display: flex;
    }

    .dot {
        animation: blink 1.4s infinite;
        animation-fill-mode: both;
    }

    .dot:nth-child(2) {
        animation-delay: 0.2s;
    }

    .dot:nth-child(3) {
        animation-delay: 0.4s;
    }

    @keyframes rotate {
        from {
            transform: rotate(0deg);
            stroke-dashoffset: 283;
        }
        to {
            transform: rotate(360deg);
            stroke-dashoffset: 0;
        }
    }

    @keyframes rotate-reverse {
        from { transform: rotate(360deg); }
        to { transform: rotate(0deg); }
    }

    @keyframes pulse {
        0%, 100% {
            transform: translate(-50%, -50%) scale(1);
            opacity: 1;
        }
        50% {
            transform: translate(-50%, -50%) scale(1.5);
            opacity: 0.5;
        }
    }

    @keyframes blink {
        0%, 100% { opacity: 0; }
        50% { opacity: 1; }
    }

    /* Glowing effect on container */
    .loading-container::after {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        width: calc(var(--size, 40px) * 1.5);
        height: calc(var(--size, 40px) * 1.5);
        background: radial-gradient(
            circle,
            rgba(9, 132, 227, 0.1) 0%,
            transparent 70%
        );
        transform: translate(-50%, -50%);
        animation: glow 2s ease-in-out infinite;
    }

    @keyframes glow {
        0%, 100% { opacity: 0.5; }
        50% { opacity: 0.8; }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .progress-ring,
        .spinner-inner,
        .center-dot,
        .dot {
            animation: none;
        }

        .loading-container::after {
            display: none;
        }
    }
</style>
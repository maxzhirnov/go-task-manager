<script>
    import Time from "svelte-time";
    import { dayjs } from "svelte-time";
    
    export let timestamp;
    export let daysToReRelative = 7;
    const daysInMs = daysToReRelative * 24 * 60 * 60 * 1000;
    $: isWithinDays = (Date.now() - new Date(timestamp).getTime()) < daysInMs;
    $: isCurrentYear = new Date(timestamp).getFullYear() === new Date().getFullYear();
    $: format = isCurrentYear ? "MMM D" : "MMM D, YYYY";
</script>

<div class="timestamp-container">
    <span class="timestamp-prefix">[TIME::</span>
    <span class="timestamp-value">
        {#if isWithinDays}
            <Time relative live timestamp={timestamp} />
        {:else}
            <Time timestamp={timestamp} {format} />
        {/if}
    </span>
    <span class="timestamp-suffix">]</span>
</div>

<style>
    .timestamp-container {
        font-family: "JetBrains Mono", monospace;
        font-size: 0.75rem;
        display: inline-flex;
        align-items: center;
        gap: 0.25rem;
        padding: 0.2rem 0.4rem;
        border-radius: 2px;
        background: rgba(9, 132, 227, 0.1);
        color: #0984e3;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;
    }

    .timestamp-container::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 1px;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(9, 132, 227, 0.2),
            transparent
        );
        animation: scan 2s linear infinite;
    }

    .timestamp-prefix {
        color: #00b894;
        font-weight: 500;
        letter-spacing: 0.05em;
    }

    .timestamp-value {
        color: #0984e3;
    }

    .timestamp-suffix {
        color: #00b894;
    }

    /* Hover effect */
    .timestamp-container:hover {
        background: rgba(9, 132, 227, 0.15);
        box-shadow: 0 0 8px rgba(9, 132, 227, 0.2);
    }

    @keyframes scan {
        from { transform: translateX(-100%); }
        to { transform: translateX(100%); }
    }

    /* Dark mode */
    @media (prefers-color-scheme: dark) {
        .timestamp-container {
            background: rgba(9, 132, 227, 0.15);
        }

        .timestamp-value {
            color: #74b9ff;
        }
    }

    /* High contrast mode */
    @media (prefers-contrast: high) {
        .timestamp-container {
            background: transparent;
            border: 1px solid currentColor;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .timestamp-container::before {
            animation: none;
        }
    }
</style>
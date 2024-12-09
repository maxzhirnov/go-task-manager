<script>
    import Time from "svelte-time";
    import { dayjs } from "svelte-time";
    
    export let timestamp;
    export let daysToReRelative = 7;
    const daysInMs = daysToReRelative * 24 * 60 * 60 * 1000;
    
    $: isWithinDays = (Date.now() - new Date(timestamp).getTime()) < daysInMs;
    $: isCurrentYear = new Date(timestamp).getFullYear() === new Date().getFullYear();
    $: format = isCurrentYear ? "MMM D, HH:mm" : "MMM D, YYYY HH:mm";
    
    let isHovered = false;
    
    function handleMouseEnter() {
        isHovered = true;
    }
    
    function handleMouseLeave() {
        isHovered = false;
    }
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div 
    class="timestamp-container"
    on:mouseenter={handleMouseEnter}
    on:mouseleave={handleMouseLeave}
>
    <span class="timestamp-prefix">[T::</span>
    <span class="timestamp-value">
        {#if isHovered}
            <Time timestamp={timestamp} format="YYYY-MM-DD HH:mm:ss" />
        {:else if isWithinDays}
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
        padding: 0.15rem 0.3rem;
        border-radius: 3px;
        /* background: #2d3436; */
        color: #74b9ff;
        transition: all 0.2s ease;
        /* border: 1px solid #0984e3; */
        cursor: crosshair;
    }

    .timestamp-prefix {
        color: #00b894;
        font-weight: 500;
        font-size: .6rem;
        opacity: 0.8;
    }

    .timestamp-value {
        color: inherit;
        min-width: 4rem;
    }

    .timestamp-suffix {
        color: #00b894;
        opacity: 0.8;
    }

    .timestamp-container:hover {
        background: #2d3436;
        border-color: #00b894;
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .timestamp-container {
            transition: none;
        }
        
        .timestamp-container:hover {
            transform: none;
        }
    }
</style>
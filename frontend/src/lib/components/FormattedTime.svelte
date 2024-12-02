<script>
    import Time from "svelte-time";
    import { dayjs } from "svelte-time";
    
    export let timestamp;
    export let daysToReRelative = 7;
    const daysInMs = daysToReRelative * 24 * 60 * 60 * 1000;
    $: isWithinDays = (Date.now() - new Date(timestamp).getTime()) < daysInMs;
    
    // Check if the year is current
    $: isCurrentYear = new Date(timestamp).getFullYear() === new Date().getFullYear();
    $: format = isCurrentYear ? "MMM D" : "MMM D, YYYY";
    </script>
    
    <div class="timestamp">
        {#if isWithinDays}
            <Time relative live timestamp={timestamp} />
        {:else}
            <Time timestamp={timestamp} {format} />
        {/if}
    </div>
    

    <style>
        .timestamp {
            font-size: 0.8rem;
            color: #666;
            font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
            display: inline-block;
            border-radius: 4px;
        }

        /* For dark backgrounds */
        @media (prefers-color-scheme: dark) {
            .timestamp {
                color: #999;
            }
            
        }
        </style>
<script>
    import Time from "svelte-time";
    import { dayjs } from "svelte-time";
    
    export let timestamp;
    const sevenDaysInMs = 3 * 24 * 60 * 60 * 1000;
    $: isWithinSevenDays = (Date.now() - new Date(timestamp).getTime()) < sevenDaysInMs;
    
    // Check if the year is current
    $: isCurrentYear = new Date(timestamp).getFullYear() === new Date().getFullYear();
    $: format = isCurrentYear ? "MMM D" : "MMM D, YYYY";
    </script>
    
    <div class="timestamp">
        {#if isWithinSevenDays}
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
            
            .timestamp:hover {
            }
        }
        </style>
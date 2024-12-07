<!-- StatsIcon.svelte -->
<script>
    export let type = 'total';
    export let size = '24';
    
    // Derive width and height from size
    $: width = size;
    $: height = size;
</script>

<svg viewBox="0 0 24 24" {width} {height} class="stats-icon {type}">
    {#if type === 'total'}
        <!-- Total Tasks Icon -->
        <g class="icon-group">
            <rect x="4" y="4" width="16" height="16" class="frame" fill="none" stroke="currentColor" stroke-width="1.5" rx="2"/>
            <path d="M8 8h8M8 12h8M8 16h8" class="lines" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <circle cx="4" cy="4" r="1" class="node"/>
            <circle cx="20" cy="4" r="1" class="node"/>
            <circle cx="4" cy="20" r="1" class="node"/>
            <circle cx="20" cy="20" r="1" class="node"/>
        </g>

    {:else if type === 'completed'}
        <!-- Completed Icon -->
        <g class="icon-group">
            <circle cx="12" cy="12" r="8" class="frame" fill="none" stroke="currentColor" stroke-width="1.5"/>
            <path d="M8 12l3 3 5-5" class="check" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 4v2M12 18v2M4 12h2M18 12h2" class="rays" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </g>

    {:else if type === 'in_progress'}
        <!-- In Progress Icon -->
        <g class="icon-group">
            <path d="M12 4v4M16 6l-2 3.5M8 6l2 3.5" class="rays" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <circle cx="12" cy="12" r="6" class="frame" fill="none" stroke="currentColor" stroke-width="1.5"/>
            <path class="pulse-ring" d="M12 12h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" transform="rotate(45, 12, 12)"/>
        </g>

    {:else if type === 'pending'}
        <!-- Pending Icon -->
        <g class="icon-group">
            <path d="M12 4L19 8V16L12 20L5 16V8L12 4Z" class="frame" fill="none" stroke="currentColor" stroke-width="1.5"/>
            <circle cx="12" cy="12" r="3" class="core" fill="none" stroke="currentColor" stroke-width="1.5"/>
            <path d="M12 9v3l2 2" class="hand" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </g>

    {:else if type === 'today'}
        <!-- Today Icon -->
        <g class="icon-group">
            <rect x="4" y="6" width="16" height="14" class="frame" fill="none" stroke="currentColor" stroke-width="1.5" rx="2"/>
            <path d="M8 4v4M16 4v4" class="pins" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <rect x="8" y="10" width="8" height="6" class="highlight" fill="none" stroke="currentColor" stroke-width="1.5" rx="1"/>
        </g>
    {/if}
</svg>

<style>
    .stats-icon {
        transition: all 0.3s ease;
    }

    .frame {
        transition: all 0.3s ease;
    }

    .node {
        fill: currentColor;
        transition: all 0.3s ease;
    }

    .lines, .check, .rays, .pins, .hand, .highlight {
        transition: all 0.3s ease;
    }

    /* Hover Animations */
    .stats-icon:hover .frame {
        stroke-dasharray: 60;
        animation: frameDash 2s linear infinite;
    }

    .stats-icon:hover .node {
        animation: pulse 1.5s infinite alternate;
    }

    .stats-icon:hover .lines,
    .stats-icon:hover .check,
    .stats-icon:hover .rays,
    .stats-icon:hover .pins,
    .stats-icon:hover .hand,
    .stats-icon:hover .highlight {
        filter: drop-shadow(0 0 2px currentColor);
    }

    .stats-icon.in_progress .pulse-ring {
        transform-origin: center;
        animation: rotate 2s linear infinite;
    }

    @keyframes frameDash {
        to {
            stroke-dashoffset: -120;
        }
    }

    @keyframes pulse {
        to {
            transform: scale(1.2);
            opacity: 0.7;
        }
    }

    @keyframes rotate {
        to {
            transform: rotate(405deg);
        }
    }
</style>
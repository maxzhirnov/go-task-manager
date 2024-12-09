<!-- DragHandle.svelte -->
<script>
    import { dragHandle } from 'svelte-dnd-action';
</script>

<div class="drag-handle" use:dragHandle>
    <svg viewBox="0 0 24 24" width="24" height="24">
        <!-- Circuit-style grip dots -->
        {#each Array(3) as _, i}
            <g class="grip-line">
                <circle 
                    cx="8" 
                    cy={6 + i * 6} 
                    r="1.5" 
                    class="grip-dot"
                />
                <line 
                    x1="10" 
                    y1={6 + i * 6} 
                    x2="16" 
                    y2={6 + i * 6} 
                    class="grip-line"
                />
                <circle 
                    cx="16" 
                    cy={6 + i * 6} 
                    r="1.5" 
                    class="grip-dot"
                />
            </g>
        {/each}
    </svg>
</div>

<style>
    .drag-handle {
        cursor: grab;
        padding: 0.25rem;
        border-radius: 2px;
        transition: all 0.3s ease;
        color: #0984e3;
    }

    .drag-handle:hover {
        background: rgba(9, 132, 227, 0.1);
    }

    .drag-handle:active {
        cursor: grabbing;
        background: rgba(9, 132, 227, 0.2);
    }

    svg {
        display: block;
    }

    .grip-dot {
        fill: currentColor;
        transition: all 0.3s ease;
    }

    .grip-line {
        stroke: currentColor;
        stroke-width: 1;
        opacity: 0.5;
        transition: all 0.3s ease;
    }

    /* Hover animations */
    .drag-handle:hover .grip-dot {
        animation: pulse 1.5s infinite;
    }

    .drag-handle:hover .grip-line {
        opacity: 0.8;
        stroke-dasharray: 6;
        animation: dash 1s linear infinite;
    }

    @keyframes pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 1;
        }
        50% {
            transform: scale(1);
            opacity: 0.8;
        }
    }

    @keyframes dash {
        to {
            stroke-dashoffset: -12;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .drag-handle:hover .grip-dot,
        .drag-handle:hover .grip-line {
            animation: none;
        }
    }
</style>
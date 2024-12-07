<!-- DragHandleAlt.svelte -->
<script>
    import { dragHandle } from 'svelte-dnd-action';
</script>

<div class="drag-handle" use:dragHandle>
    <svg viewBox="0 0 24 24" width="24" height="24">
        <!-- Circuit board pattern -->
        <path 
            class="circuit-path"
            d="M6 4v16M18 4v16"
            stroke="currentColor"
            stroke-width="1.5"
            fill="none"
        />
        {#each Array(3) as _, i}
            <g class="connection-point">
                <circle 
                    cx="6" 
                    cy={8 + i * 4} 
                    r="2" 
                    class="node"
                />
                <circle 
                    cx="18" 
                    cy={8 + i * 4} 
                    r="2" 
                    class="node"
                />
                <line 
                    x1="8" 
                    y1={8 + i * 4} 
                    x2="16" 
                    y2={8 + i * 4} 
                    class="connection-line"
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
        box-shadow: 0 0 8px rgba(9, 132, 227, 0.2);
    }

    .drag-handle:active {
        cursor: grabbing;
        background: rgba(9, 132, 227, 0.2);
    }

    svg {
        display: block;
    }

    .circuit-path {
        opacity: 0.3;
    }

    .node {
        fill: currentColor;
        stroke: none;
        transition: all 0.3s ease;
    }

    .connection-line {
        stroke: currentColor;
        stroke-width: 1;
        opacity: 0.5;
        stroke-dasharray: 4;
        transition: all 0.3s ease;
    }

    /* Hover animations */
    .drag-handle:hover .node {
        animation: nodeGlow 1.5s infinite alternate;
    }

    .drag-handle:hover .connection-line {
        animation: lineDash 1s linear infinite;
    }

    @keyframes nodeGlow {
        from {
            filter: drop-shadow(0 0 0 currentColor);
        }
        to {
            filter: drop-shadow(0 0 2px currentColor);
        }
    }

    @keyframes lineDash {
        to {
            stroke-dashoffset: -8;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .drag-handle:hover .node,
        .drag-handle:hover .connection-line {
            animation: none;
        }
    }
</style>
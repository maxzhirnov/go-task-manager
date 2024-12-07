<svg viewBox="0 0 24 24" width="20" height="20" class="filter-icon">
    <!-- Main filter shape -->
    <path
        class="filter-path"
        d="M4 6h16M6 12h12M8 18h8"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
        stroke-linejoin="round"
    />
    
    <!-- Connection nodes -->
    <circle cx="4" cy="6" r="1" class="node" />
    <circle cx="20" cy="6" r="1" class="node" />
    <circle cx="6" cy="12" r="1" class="node" />
    <circle cx="18" cy="12" r="1" class="node" />
    <circle cx="8" cy="18" r="1" class="node" />
    <circle cx="16" cy="18" r="1" class="node" />
</svg>

<style>
    .filter-icon {
        transition: all 0.3s ease;
    }

    .filter-path {
        transition: all 0.3s ease;
    }

    .node {
        fill: currentColor;
        transition: all 0.3s ease;
    }

    /* Hover animations */
    .filter-icon:hover .filter-path {
        stroke-dasharray: 24;
        animation: dash 1.5s linear infinite;
    }

    .filter-icon:hover .node {
        animation: pulse 1.5s infinite alternate;
    }

    @keyframes dash {
        to {
            stroke-dashoffset: -48;
        }
    }

    @keyframes pulse {
        from {
            transform: scale(1);
            opacity: 1;
        }
        to {
            transform: scale(1.5);
            opacity: 0.5;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .filter-icon:hover .filter-path,
        .filter-icon:hover .node {
            animation: none;
        }
    }
</style>
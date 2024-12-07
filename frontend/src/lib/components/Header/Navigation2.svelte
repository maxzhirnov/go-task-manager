<script>
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';

    const navItems = [
        { path: '/tasks', label: 'TASKS', hex: '0x01' },
        { path: '/statistics', label: 'STATS', hex: '0x02' }
    ];
</script>

<nav class="cyber-nav">
    <div class="nav-grid">
        {#each navItems as item}
            <a 
                class="nav-node {$page.url.pathname === item.path ? 'active' : ''}"
                href={item.path}
                on:click|preventDefault={() => goto(item.path)}
            >
                <div class="node-content">
                    <!-- Hexagonal Frame -->
                    <svg class="hex-frame" viewBox="0 0 100 100">
                        <polygon 
                            points="50,3 93.5,25 93.5,75 50,97 6.5,75 6.5,25" 
                            class="hex-shape"
                        />
                        <!-- Circuit lines -->
                        <path d="M25,50 H45" class="circuit-line" />
                        <path d="M55,50 H75" class="circuit-line" />
                        <!-- Connection points -->
                        <circle cx="25" cy="50" r="2" class="connection-point" />
                        <circle cx="75" cy="50" r="2" class="connection-point" />
                    </svg>

                    <!-- Content -->
                    <div class="node-label">
                        <span class="label-text">{item.label}</span>
                        <span class="hex-code">{item.hex}</span>
                    </div>

                    <!-- Status indicator -->
                    <div class="status-indicator"></div>
                </div>
            </a>
        {/each}
    </div>

    <!-- Circuit decoration -->
    <div class="circuit-lines">
        <div class="line"></div>
        <div class="line"></div>
    </div>
</nav>

<style>
    .cyber-nav {
        position: relative;
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
        overflow: hidden;
    }

    .nav-grid {
        display: flex;
        gap: 2rem;
        position: relative;
        z-index: 2;
    }

    .nav-node {
        text-decoration: none;
        color: #0984e3;
        position: relative;
        transition: all 0.3s ease;
    }

    .node-content {
        position: relative;
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.5rem 1rem;
    }

    .hex-frame {
        width: 24px;
        height: 24px;
    }

    .hex-shape {
        fill: none;
        stroke: currentColor;
        stroke-width: 1;
        transition: all 0.3s ease;
    }

    .circuit-line {
        stroke: currentColor;
        stroke-width: 0.5;
        opacity: 0.5;
        stroke-dasharray: 10;
        stroke-dashoffset: 20;
        transition: all 0.3s ease;
    }

    .connection-point {
        fill: currentColor;
        transition: all 0.3s ease;
    }

    .node-label {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
    }

    .label-text {
        font-size: 0.9rem;
        font-weight: 500;
    }

    .hex-code {
        font-size: 0.6rem;
        opacity: 0.7;
    }

    .status-indicator {
        position: absolute;
        top: 50%;
        left: 0.5rem;
        width: 4px;
        height: 4px;
        background: currentColor;
        transform: translateY(-50%);
        clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
        opacity: 0;
        transition: all 0.3s ease;
    }

    /* Active state */
    .nav-node.active {
        color: #00b894;
    }

    .nav-node.active .hex-shape {
        stroke-width: 2;
        filter: drop-shadow(0 0 3px currentColor);
    }

    .nav-node.active .circuit-line {
        opacity: 1;
        stroke-dashoffset: 0;
    }

    .nav-node.active .connection-point {
        animation: pulse 2s infinite;
    }

    .nav-node.active .status-indicator {
        opacity: 1;
    }

    /* Hover effects */
    .nav-node:hover {
        color: #00b894;
    }

    .nav-node:hover .circuit-line {
        opacity: 1;
    }

    .nav-node:hover .connection-point {
        transform: scale(1.5);
    }

    /* Circuit decoration */
    .circuit-lines {
        position: absolute;
        inset: 0;
        pointer-events: none;
        opacity: 0.1;
    }

    .line {
        position: absolute;
        background: linear-gradient(
            90deg,
            transparent,
            #00b894,
            transparent
        );
        height: 1px;
    }

    .line:nth-child(1) {
        top: 30%;
        left: 0;
        right: 0;
    }

    .line:nth-child(2) {
        bottom: 30%;
        left: 0;
        right: 0;
    }

    @keyframes pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 1;
        }
        50% {
            transform: scale(1.5);
            opacity: 0.5;
        }
    }

    @media screen and (max-width: 768px) {
        .nav-grid {
            justify-content: center;
        }
    }
</style>
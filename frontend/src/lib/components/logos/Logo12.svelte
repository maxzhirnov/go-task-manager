<!-- src/lib/components/Logo.svelte -->
<script>
    import { goto } from '$app/navigation';
    
    export let size = '32';
    export let clickable = true;

    $: width = size;
    $: height = size;
</script>

<div class="logo-container">
    {#if clickable}
        <a href="/tasks" class="logo-link" on:click|preventDefault={() => goto('/tasks')}>
            <div class="logo">
                <div class="terminal-border">
                    <div class="terminal-header">
                        <div class="terminal-buttons">
                            <span class="btn close"></span>
                            <span class="btn minimize"></span>
                            <span class="btn maximize"></span>
                        </div>
                    </div>
                    <div class="terminal-screen">
                        <div class="ascii-art">
                            [A::H]
                            >_sys
                        </div>
                        <div class="command-line">
                            <span class="cursor">█</span>
                        </div>
                        <div class="matrix-rain"></div>
                    </div>
                </div>
                <div class="text-container">
                    <div class="logo-text">
                        <span class="action">ACTION</span>
                        <span class="separator">::</span>
                        <span class="hub">HUB</span>
                    </div>
                    <div class="status-line">[CORE=ACTIVE]</div>
                </div>
            </div>
        </a>
    {:else}
        <div class="logo">
            <a href="/tasks" class="logo-link" on:click|preventDefault={() => console.log('ACTIONHUB')}>
                <div class="logo">
                    <div class="terminal-border">
                        <div class="terminal-header">
                            <div class="terminal-buttons">
                                <span class="btn close"></span>
                                <span class="btn minimize"></span>
                                <span class="btn maximize"></span>
                            </div>
                        </div>
                        <div class="terminal-screen">
                            <div class="ascii-art">
                                [A::H]
                                >_sys
                            </div>
                            <div class="command-line">
                                <span class="cursor">█</span>
                            </div>
                            <div class="matrix-rain"></div>
                        </div>
                    </div>
                    <div class="text-container">
                        <div class="logo-text">
                            <span class="action">ACTION</span>
                            <span class="separator">::</span>
                            <span class="hub">HUB</span>
                        </div>
                        <div class="status-line">[CORE=ACTIVE]</div>
                    </div>
                </div>
            </a>
        </div>
    {/if}
</div>

<style>
    .logo-container {
        display: flex;
        align-items: center;
    }

    .logo {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", "Fira Code", monospace;
    }

    .terminal-border {
        background: #2b2b2b;
        border-radius: 4px;
        box-shadow: 0 2px 8px rgba(0,0,0,0.15);
        overflow: hidden;
        width: fit-content;
    }

    .terminal-header {
        background: #1c1c1c;
        padding: 4px 6px;
        display: flex;
        align-items: center;
        border-bottom: 1px solid #3c3c3c;
    }

    .terminal-buttons {
        display: flex;
        gap: 4px;
    }

    .btn {
        width: 6px;
        height: 6px;
        border-radius: 50%;
    }

    .close { background: #ff5f56; }
    .minimize { background: #ffbd2e; }
    .maximize { background: #27c93f; }

    .terminal-screen {
        padding: 4px 6px;
        background: #1c1c1c;
        position: relative;
        min-width: 60px;
        height: 32px;
    }

    .ascii-art {
        color: #00ff00;
        font-size: 0.6rem;
        white-space: pre;
        opacity: 0.8;
        line-height: 1;
    }

    .command-line {
        color: #fff;
        font-size: 0.7rem;
        display: flex;
        gap: 4px;
    }

    .cursor {
        animation: blink 1s steps(1) infinite;
    }

    .matrix-rain {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        pointer-events: none;
        opacity: 0.1;
        background: linear-gradient(180deg,
            transparent 0%,
            rgba(0, 255, 0, 0.2) 50%,
            transparent 100%
        );
        background-size: 100% 20px;
        animation: matrix 1s linear infinite;
    }

    .text-container {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .logo-text {
        font-size: 1rem;
        font-weight: bold;
        display: flex;
        align-items: center;
        gap: 2px;
    }

    .action {
        color: #2c3e50;
    }

    .separator {
        color: #7f8c8d;
        font-weight: normal;
    }

    .hub {
        color: #2980b9;
    }

    .status-line {
        font-size: 0.5rem;
        color: #2ecc71;
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    @keyframes matrix {
        0% { background-position: 0 0; }
        100% { background-position: 0 20px; }
    }

    /* Hover effects */
    .logo:hover .terminal-screen {
        background: #1a1a1a;
    }

    .logo:hover .matrix-rain {
        opacity: 0.15;
    }

    .logo:hover .ascii-art {
        opacity: 1;
    }

    .logo:hover .status-line {
        animation: glitch 0.3s ease infinite;
    }

    @keyframes glitch {
        0%, 100% { transform: translate(0); }
        20% { transform: translate(-1px, 1px); }
        40% { transform: translate(-1px, -1px); }
        60% { transform: translate(1px, 1px); }
    }

    /* Scan line effect */
    .terminal-screen::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            transparent 50%,
            rgba(0, 0, 0, 0.1) 50%
        );
        background-size: 100% 4px;
        pointer-events: none;
    }
</style>
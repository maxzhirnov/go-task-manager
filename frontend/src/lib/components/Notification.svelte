<script>
    import { successMessage } from '$lib/stores';
</script>

{#if $successMessage}
    <div class="success-terminal" role="alert">
        <div class="success-content">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">SYSTEM_SUCCESS.log</span>
            </div>
            <div class="success-body">
                <div class="status-line">
                    <span class="prompt">>_</span>
                    <span class="status-code">[SUCCESS_0x{Math.floor(Math.random() * 1000).toString(16).toUpperCase().padStart(3, '0')}]</span>
                </div>
                <div class="success-text">{$successMessage}</div>
                <div class="progress-bar">
                    <div class="progress-fill"></div>
                </div>
            </div>
        </div>
        <div class="matrix-effect"></div>
        <div class="scan-line"></div>
    </div>
{/if}

<style>
    .success-terminal {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 1000;
        max-width: 400px;
        background: #1c1c1c;
        border: 1px solid #00b894;
        border-radius: 4px;
        overflow: hidden;
        font-family: 'JetBrains Mono', monospace;
        animation: slideIn 0.3s ease-out;
        box-shadow: 0 0 15px rgba(0, 184, 148, 0.2);
    }

    .success-content {
        position: relative;
        z-index: 2;
    }

    .terminal-header {
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(0, 184, 148, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .success-body {
        padding: 1rem;
    }

    .status-line {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
        animation: blink 1s steps(1) infinite;
    }

    .status-code {
        color: #00b894;
        font-size: 0.8rem;
        font-weight: bold;
    }

    .success-text {
        color: #fff;
        font-size: 0.9rem;
        padding-left: 1.5rem;
        line-height: 1.4;
        margin-bottom: 1rem;
    }

    .progress-bar {
        height: 2px;
        background: rgba(0, 184, 148, 0.2);
        border-radius: 1px;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        background: #00b894;
        animation: progress 3s linear forwards;
    }

    .matrix-effect {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            rgba(0, 184, 148, 0.05) 50%,
            rgba(0, 184, 148, 0) 50%
        );
        background-size: 100% 4px;
        pointer-events: none;
        opacity: 0.5;
    }

    .scan-line {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: rgba(0, 184, 148, 0.2);
        animation: scan 2s linear infinite;
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    @keyframes progress {
        0% { width: 0; }
        100% { width: 100%; }
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    @keyframes scan {
        from { transform: translateY(-100%); }
        to { transform: translateY(100%); }
    }

    /* Auto-dismiss animation */
    .success-terminal {
        animation: 
            slideIn 0.3s ease-out,
            fadeOut 3s forwards;
    }

    @keyframes fadeOut {
        0% { opacity: 1; }
        70% { opacity: 1; }
        100% { 
            opacity: 0;
            visibility: hidden;
        }
    }

    /* Mobile responsiveness */
    @media (max-width: 480px) {
        .success-terminal {
            top: 10px;
            right: 10px;
            left: 10px;
            max-width: none;
        }

        .success-body {
            padding: 0.8rem;
        }

        .success-text {
            font-size: 0.8rem;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .success-terminal {
            animation: slideIn 0.3s ease-out;
        }

        .matrix-effect,
        .scan-line {
            display: none;
        }

        .prompt {
            animation: none;
        }

        .progress-fill {
            animation: none;
            width: 100%;
        }
    }

    /* High contrast mode */
    @media (prefers-contrast: high) {
        .success-terminal {
            border-width: 2px;
        }

        .success-text {
            color: #ffffff;
        }
    }
</style>
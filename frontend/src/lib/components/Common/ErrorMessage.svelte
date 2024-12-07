<script>
    import { errorMessage } from '$lib/stores.js';
</script>

{#if $errorMessage}
    <div class="error-terminal" role="alert">
        <div class="error-content">
            <div class="error-header">
                <span class="error-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="error-title">SYSTEM_ERROR.log</span>
            </div>
            <div class="error-body">
                <div class="error-line">
                    <span class="prompt">>_</span>
                    <span class="error-code">[ERROR_0x{Math.floor(Math.random() * 1000).toString(16).toUpperCase().padStart(3, '0')}]</span>
                </div>
                <div class="error-text">{$errorMessage}</div>
            </div>
        </div>
        <div class="glitch-effect"></div>
        <div class="scan-line"></div>
    </div>
{/if}

<style>
    .error-terminal {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 1000;
        max-width: 400px;
        background: #1c1c1c;
        border: 1px solid #ff6b6b;
        border-radius: 4px;
        overflow: hidden;
        font-family: 'JetBrains Mono', monospace;
        animation: slideIn 0.3s ease-out, fadeOut 3s forwards;
        box-shadow: 0 0 15px rgba(255, 107, 107, 0.2);
    }

    .error-content {
        position: relative;
        z-index: 2;
    }

    .error-header {
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(255, 107, 107, 0.2);
    }

    .error-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .error-title {
        color: #ff6b6b;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .error-body {
        padding: 1rem;
    }

    .error-line {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .prompt {
        color: #ff6b6b;
        font-size: 0.9rem;
        animation: blink 1s steps(1) infinite;
    }

    .error-code {
        color: #ff6b6b;
        font-size: 0.8rem;
        font-weight: bold;
    }

    .error-text {
        color: #fff;
        font-size: 0.9rem;
        padding-left: 1.5rem;
        line-height: 1.4;
    }

    .glitch-effect {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(255, 107, 107, 0.05);
        pointer-events: none;
        animation: glitch 1s infinite;
    }

    .scan-line {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: rgba(255, 107, 107, 0.2);
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

    @keyframes fadeOut {
        0% { opacity: 1; }
        70% { opacity: 1; }
        100% { 
            opacity: 0;
            visibility: hidden;
        }
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    @keyframes glitch {
        0% { opacity: 0; }
        1% { opacity: 0.1; transform: translateX(-2px); }
        2% { opacity: 0; transform: translateX(2px); }
        100% { opacity: 0; }
    }

    @keyframes scan {
        from { transform: translateY(-100%); }
        to { transform: translateY(100%); }
    }

    /* Mobile responsiveness */
    @media (max-width: 480px) {
        .error-terminal {
            top: 10px;
            right: 10px;
            left: 10px;
            max-width: none;
        }

        .error-body {
            padding: 0.8rem;
        }

        .error-text {
            font-size: 0.8rem;
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .error-terminal {
            animation: slideIn 0.3s ease-out;
        }

        .glitch-effect,
        .scan-line {
            display: none;
        }

        .prompt {
            animation: none;
        }
    }
</style>
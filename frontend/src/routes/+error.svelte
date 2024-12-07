<script>
    import { page } from '$app/stores';

    $: errorCode = $page.status;
    $: errorTitle = $page.error?.message || 'SYSTEM_ERROR';
    $: errorMessage = {
        404: "ERROR::REQUESTED_RESOURCE_NOT_FOUND",
        500: "ERROR::INTERNAL_SYSTEM_FAILURE",
        403: "ERROR::ACCESS_DENIED",
    }[errorCode] || 'ERROR::UNEXPECTED_SYSTEM_FAILURE';
</script>

<div class="error-container">
    <div class="terminal-box">
        <div class="terminal-header">
            <span class="terminal-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </span>
            <span class="terminal-title">SYSTEM_ERROR.log</span>
        </div>
        
        <div class="error-content">
            <div class="error-status">
                <div class="status-code">{errorCode}</div>
                <div class="status-lines">
                    <span class="status-line">> ERROR_DETECTED</span>
                    <span class="status-line">> SYSTEM_COMPROMISED</span>
                    <span class="status-line blink">> AWAITING_USER_ACTION</span>
                </div>
            </div>

            <div class="error-details">
                <div class="error-title">[ERROR_TYPE]::{errorTitle}</div>
                <div class="error-message">{errorMessage}</div>
            </div>

            <div class="matrix-bg"></div>

            <div class="action-buttons">
                <button class="cyber-button primary" on:click={() => window.location.href = "/"}>
                    <span class="btn-icon">⌂</span>
                    <span class="btn-text">RETURN_TO_HOME</span>
                </button>
                <button class="cyber-button secondary" on:click={() => history.back()}>
                    <span class="btn-icon">←</span>
                    <span class="btn-text">PREVIOUS_NODE</span>
                </button>
            </div>
        </div>
    </div>
</div>

<style>
    .error-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 2rem;
        font-family: 'JetBrains Mono', monospace;
        background-color: #0f1215;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                #161b22 1px,
                transparent 1px
            );
        background-size: 24px 24px;
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #ff6b6b;
        border-radius: 4px;
        width: 100%;
        max-width: 600px;
        overflow: hidden;
        position: relative;
        animation: glitch 0.3s infinite;
    }

    .terminal-header {
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(255, 107, 107, 0.2);
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
        color: #ff6b6b;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .error-content {
        padding: 2rem;
        position: relative;
    }

    .error-status {
        display: flex;
        align-items: center;
        gap: 2rem;
        margin-bottom: 2rem;
    }

    .status-code {
        font-size: 4rem;
        font-weight: bold;
        color: #ff6b6b;
        text-shadow: 0 0 10px rgba(255, 107, 107, 0.5);
        animation: flicker 2s infinite;
    }

    .status-lines {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .status-line {
        color: #ff6b6b;
        font-size: 0.8rem;
    }

    .blink {
        animation: blink 1s steps(1) infinite;
    }

    .error-details {
        margin-bottom: 2rem;
        padding: 1rem;
        background: rgba(255, 107, 107, 0.1);
        border-left: 3px solid #ff6b6b;
    }

    .error-title {
        color: #ff6b6b;
        font-size: 1.2rem;
        margin-bottom: 0.5rem;
    }

    .error-message {
        color: #fff;
        font-size: 0.9rem;
    }

    .action-buttons {
        display: flex;
        gap: 1rem;
    }

    .cyber-button {
        flex: 1;
        background: transparent;
        border: 1px solid currentColor;
        padding: 0.8rem;
        font-family: inherit;
        font-size: 0.8rem;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        cursor: pointer;
        position: relative;
        overflow: hidden;
        transition: all 0.3s ease;
    }

    .cyber-button.primary {
        color: #00b894;
        border-color: #00b894;
    }

    .cyber-button.secondary {
        color: #0984e3;
        border-color: #0984e3;
    }

    .cyber-button:hover {
        background: rgba(255, 255, 255, 0.1);
        box-shadow: 0 0 10px currentColor;
    }

    .matrix-bg {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            rgba(255, 107, 107, 0.05) 50%,
            rgba(255, 107, 107, 0) 50%
        );
        background-size: 100% 4px;
        pointer-events: none;
        z-index: -1;
    }

    @keyframes glitch {
        0% { transform: translate(0); }
        20% { transform: translate(-2px, 2px); }
        40% { transform: translate(-2px, -2px); }
        60% { transform: translate(2px, 2px); }
        80% { transform: translate(2px, -2px); }
        100% { transform: translate(0); }
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    @keyframes flicker {
        0%, 100% { opacity: 1; }
        92% { opacity: 1; }
        93% { opacity: 0.3; }
        94% { opacity: 1; }
        95% { opacity: 0.5; }
        96% { opacity: 1; }
    }

    /* Mobile responsiveness */
    @media (max-width: 480px) {
        .error-container {
            padding: 1rem;
        }

        .error-status {
            flex-direction: column;
            gap: 1rem;
            text-align: center;
        }

        .action-buttons {
            flex-direction: column;
        }
    }
</style>
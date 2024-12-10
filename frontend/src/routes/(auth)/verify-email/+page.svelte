
<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { Analytics } from '$lib/analytics';

    let status = 'verifying';
    let message = 'INITIALIZING VERIFICATION SEQUENCE...';
    let progress = 0;

    onMount(async () => {
        const token = $page.url.searchParams.get('token');
        if (!token) {
            status = 'error';
            message = 'ERROR: VERIFICATION TOKEN NOT FOUND';
            return;
        }

        // Simulate progress for better UX
        const progressInterval = setInterval(() => {
            if (progress < 90) progress += 10;
        }, 200);

        try {
            const response = await fetch(`/api/verify-email?token=${token}`);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'VERIFICATION SEQUENCE FAILED');
            }

            progress = 100;
            status = 'success';
            message = 'EMAIL VERIFICATION COMPLETE. ACCESS GRANTED.';
            Analytics.track('Email Verification Completed');
        } catch (error) {
            status = 'error';
            message = `ERROR: ${error.message}`;
        } finally {
            clearInterval(progressInterval);
        }
    });
</script>

<div class="container">
    <div class="terminal-box">
        <div class="terminal-header">
            <span class="terminal-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </span>
            <span class="terminal-title">EMAIL_VERIFICATION.sys</span>
        </div>
        
        <div class="terminal-content">
            <div class="system-status">
                <span class="status-line">VERIFICATION PROTOCOL v2.0</span>
                <span class="status-line">SECURE CHANNEL: ACTIVE</span>
                <span class="status-line blink">>_ PROCESSING</span>
            </div>

            {#if status === 'verifying'}
                <div class="verification-progress">
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: {progress}%"></div>
                    </div>
                    <div class="progress-status">
                        <span class="percentage">{progress}%</span>
                        <span class="status-text">{message}</span>
                    </div>
                    <div class="loading-matrix">
                        {#each Array(5) as _, i}
                            <div class="matrix-line">
                                {#each Array(8) as _, j}
                                    <span class="matrix-char">{Math.random() > 0.5 ? '1' : '0'}</span>
                                {/each}
                            </div>
                        {/each}
                    </div>
                </div>
            {:else if status === 'success'}
                <div class="success-terminal">
                    <div class="success-icon">✓</div>
                    <div class="message-box success">
                        <span class="message-prefix">[SUCCESS]</span>
                        <span class="message-text">{message}</span>
                    </div>
                    <button 
                        class="cyber-button"
                        on:click={() => goto('/login')}
                    >
                        <span class="btn-icon">⚡</span>
                        ACCESS_LOGIN_PORTAL
                    </button>
                </div>
            {:else}
                <div class="error-terminal">
                    <div class="error-icon">⚠</div>
                    <div class="message-box error">
                        <span class="message-prefix">[ERROR]</span>
                        <span class="message-text">{message}</span>
                    </div>
                    <button 
                        class="cyber-button secondary"
                        on:click={() => window.location.reload()}
                    >
                        <span class="btn-icon">↻</span>
                        RETRY_VERIFICATION
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .container {
        color: #fff;
        max-width: 450px;
        margin: 50px auto;
        padding: 1rem;
        font-family: "JetBrains Mono", monospace;
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
    }

    .terminal-box::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.05) 1px,
                transparent 1px
            );
        background-size: 10px 10px;
        pointer-events: none;
    }

    .terminal-header {
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
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

    .terminal-content {
        padding: 1.5rem;
        position: relative;
    }

    .system-status {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
        margin-bottom: 2rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
        padding-bottom: 1rem;
    }

    .status-line {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .blink {
        animation: blink 1s steps(1) infinite;
    }

    /* Terminal input styles */
    .input-group {
        margin-bottom: 1.5rem;
    }

    .input-label {
        color: #00b894;
        font-size: 0.7rem;
        margin-bottom: 0.5rem;
        letter-spacing: 0.1em;
    }

    .input-wrapper {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
        padding: 0 0.5rem;
        transition: all 0.3s ease;
    }

    .input-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
        user-select: none;
    }

    /* Terminal animations */
    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    /* Terminal decorative elements */
    .terminal-content::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(9, 132, 227, 0.2),
            transparent
        );
    }

    /* Terminal glitch effect */
    .terminal-glitch {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 184, 148, 0.05);
        opacity: 0;
        pointer-events: none;
        animation: glitch 2s infinite;
    }

    @keyframes glitch {
        0% { opacity: 0; transform: translateX(0); }
        1% { opacity: 0.1; transform: translateX(-2px); }
        2% { opacity: 0; transform: translateX(2px); }
        100% { opacity: 0; transform: translateX(0); }
    }

    /* Terminal scanline effect */
    .terminal-scanline {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 2px;
        background: rgba(0, 184, 148, 0.1);
        opacity: 0.75;
        animation: scan 6s linear infinite;
    }

    @keyframes scan {
        0% { transform: translateY(-100%); }
        100% { transform: translateY(100vh); }
    }

    /* Terminal response messages */
    .terminal-response {
        color: #00b894;
        font-size: 0.8rem;
        padding: 0.5rem;
        background: rgba(0, 184, 148, 0.05);
        border-left: 2px solid #00b894;
        margin: 1rem 0;
    }

    .terminal-error {
        color: #ff6b6b;
        font-size: 0.8rem;
        padding: 0.5rem;
        background: rgba(255, 107, 107, 0.05);
        border-left: 2px solid #ff6b6b;
        margin: 1rem 0;
    }

    /* Terminal scrollbar */
    .terminal-content::-webkit-scrollbar {
        width: 8px;
    }

    .terminal-content::-webkit-scrollbar-track {
        background: #2d3436;
    }

    .terminal-content::-webkit-scrollbar-thumb {
        background: #0984e3;
        border-radius: 4px;
    }

    .terminal-content::-webkit-scrollbar-thumb:hover {
        background: #00b894;
    }

    /* Terminal selection */
    .terminal-content ::selection {
        background: rgba(0, 184, 148, 0.3);
        color: #fff;
    }

    /* Terminal focus styles */
    .terminal-box:focus-within {
        box-shadow: 0 0 0 2px rgba(0, 184, 148, 0.2);
    }

    /* Terminal responsive adjustments */
    @media (max-width: 480px) {
        .terminal-content {
            padding: 1rem;
        }

        .status-line {
            font-size: 0.6rem;
        }

        .input-wrapper {
            padding: 0 0.3rem;
        }

        .prompt {
            font-size: 0.8rem;
        }
    }

    .verification-progress {
        padding: 1rem;
    }

    .progress-bar {
        height: 4px;
        background: #2d3436;
        border-radius: 2px;
        margin: 1rem 0;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        background: #00b894;
        transition: width 0.3s ease;
        position: relative;
    }

    .progress-fill::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.2),
            transparent
        );
        animation: progress-glow 1s linear infinite;
    }

    .progress-status {
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: #00b894;
        font-size: 0.8rem;
    }

    .loading-matrix {
        margin-top: 2rem;
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
        opacity: 0.5;
    }

    .matrix-line {
        display: flex;
        gap: 0.5rem;
        justify-content: center;
    }

    .matrix-char {
        color: #00b894;
        font-size: 0.8rem;
        animation: matrix-flicker 0.5s infinite;
    }

    .message-box {
        padding: 1rem;
        margin: 1rem 0;
        border-radius: 3px;
        display: flex;
        gap: 0.5rem;
    }

    .message-box.success {
        background: rgba(0, 184, 148, 0.1);
        border: 1px solid #00b894;
    }

    .message-box.error {
        background: rgba(255, 107, 107, 0.1);
        border: 1px solid #ff6b6b;
    }

    .cyber-button {
        /* ... (previous button styles) ... */
    }

    @keyframes progress-glow {
        0% { transform: translateX(-100%); }
        100% { transform: translateX(100%); }
    }

    @keyframes matrix-flicker {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }

    .success-terminal, .error-terminal {
        text-align: center;
        padding: 2rem 1rem;
    }

    .success-icon, .error-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
        animation: icon-pulse 2s infinite;
    }

    .success-icon {
        color: #00b894;
    }

    .error-icon {
        color: #ff6b6b;
    }

    .message-prefix {
        font-weight: bold;
    }

    .success .message-prefix {
        color: #00b894;
    }

    .error .message-prefix {
        color: #ff6b6b;
    }

    .message-text {
        color: #fff;
    }

    .cyber-button {
        width: 100%;
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.8rem;
        border-radius: 3px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        font-family: inherit;
        font-size: 0.8rem;
        transition: all 0.3s ease;
        margin-top: 1.5rem;
        text-transform: uppercase;
        letter-spacing: 1px;
        position: relative;
        overflow: hidden;
    }

    .cyber-button::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(0, 184, 148, 0.2),
            transparent
        );
        transition: 0.5s;
    }

    .cyber-button:hover::before {
        left: 100%;
    }

    .cyber-button:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    .cyber-button.secondary {
        border-color: #0984e3;
        color: #0984e3;
    }

    .cyber-button.secondary:hover {
        background: rgba(9, 132, 227, 0.1);
        box-shadow: 0 0 8px rgba(9, 132, 227, 0.3);
    }

    .terminal-content {
        position: relative;
        z-index: 1;
    }

    .terminal-content::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: repeating-linear-gradient(
            0deg,
            rgba(0, 0, 0, 0.15) 0px,
            rgba(0, 0, 0, 0.15) 1px,
            transparent 1px,
            transparent 2px
        );
        pointer-events: none;
        z-index: -1;
    }

    .status-line {
        position: relative;
        padding-left: 1.2rem;
    }

    .status-line::before {
        content: '>';
        position: absolute;
        left: 0;
        color: #00b894;
    }

    @keyframes icon-pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 1;
        }
        50% {
            transform: scale(1.1);
            opacity: 0.8;
        }
    }

    @keyframes scan-line {
        0% {
            transform: translateY(-100%);
        }
        100% {
            transform: translateY(100%);
        }
    }

    /* Scan line effect */
    .terminal-box::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 5px;
        background: linear-gradient(
            to bottom,
            transparent,
            rgba(0, 184, 148, 0.2),
            transparent
        );
        animation: scan-line 3s linear infinite;
        pointer-events: none;
    }

    /* Responsive styles */
    @media (max-width: 480px) {
        .container {
            margin: 20px auto;
            padding: 0.5rem;
        }

        .terminal-content {
            padding: 1rem;
        }

        .success-icon, .error-icon {
            font-size: 2rem;
        }

        .message-box {
            padding: 0.8rem;
            font-size: 0.9rem;
        }

        .loading-matrix {
            display: none; /* Hide matrix on mobile for better performance */
        }

        .cyber-button {
            padding: 0.6rem;
            font-size: 0.7rem;
        }
    }

    /* Dark mode support */
    @media (prefers-color-scheme: dark) {
        body {
            background-color: #0f1215;
        }
    }

    /* High contrast mode support */
    @media (prefers-contrast: high) {
        .terminal-box {
            border-width: 2px;
        }

        .message-text {
            color: #ffffff;
        }
    }

    /* Reduced motion support */
    @media (prefers-reduced-motion: reduce) {
        .cyber-button::before,
        .terminal-box::after,
        .progress-fill::after {
            animation: none;
        }

        .cyber-button:hover::before {
            display: none;
        }
    }
</style>
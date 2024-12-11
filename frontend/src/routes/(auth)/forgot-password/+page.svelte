<script>
    import { showError, showSuccess } from '$lib/stores';
    import { api } from '$lib/api';

    let email = '';
    let loading = false;
    let submitted = false;

    async function handleSubmit() {
        if (!email.trim()) {
            showError('VALIDATION_ERROR: Email field required');
            return;
        }

        loading = true;

        try {
            const response = await api.requestPasswordReset({ email });
            submitted = true;
            console.log('Password reset request response:', response);
            showSuccess('RESET_LINK_TRANSMITTED: Check your inbox');
        } catch (error) {
            console.error('Password reset request failed:', {
                error: error,
                message: error.message,
                status: error.status,
                details: error.details
            });
            
            if (error.status === 404) {
                showError('ERROR_404: Email not found in database');
            } else if (error.status === 429) {
                showError('ERROR_429: Too many reset attempts detected');
            } else {
                showError('ERROR: Reset link transmission failed');
            }
        } finally {
            loading = false;
        }
    }
</script>

<div class="container">
    <div class="terminal-box">
        <div class="terminal-header">
            <span class="terminal-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </span>
            <span class="terminal-title">PASSWORD_RESET.exe</span>
        </div>
        
        <div class="terminal-content">
            {#if !submitted}
                <div class="system-status">
                    <span class="status-line">INITIALIZING PASSWORD RESET PROTOCOL</span>
                    <span class="status-line">SECURE CHANNEL: ESTABLISHED</span>
                    <span class="status-line blink">>_ AWAITING EMAIL VERIFICATION</span>
                </div>

                <form on:submit|preventDefault={handleSubmit}>
                    <div class="input-group">
                        <div class="input-label">[REGISTERED_EMAIL]</div>
                        <div class="input-wrapper">
                            <span class="prompt">>_</span>
                            <input 
                                type="email" 
                                bind:value={email} 
                                placeholder="Enter registered email"
                                required
                                disabled={loading}
                            >
                        </div>
                    </div>

                    <button type="submit" class="terminal-button" disabled={loading}>
                        <span class="btn-icon">⚡</span>
                        <span class="btn-text">
                            {loading ? 'TRANSMITTING...' : 'INITIATE_RESET_SEQUENCE'}
                        </span>
                    </button>
                </form>
            {:else}
                <div class="success-terminal">
                    <div class="success-icon">✓</div>
                    <div class="success-content">
                        <span class="status-line">RESET LINK TRANSMITTED TO:</span>
                        <span class="email-highlight">{email}</span>
                        <span class="status-line">STATUS: AWAITING USER CONFIRMATION</span>
                        <span class="status-line">ACTION: CHECK INBOX FOR FURTHER INSTRUCTIONS</span>
                    </div>
                </div>
            {/if}

            <div class="system-footer">
                <a href="/login" class="system-link">
                    <span class="btn-icon">←</span>
                    RETURN_TO_LOGIN
                </a>
            </div>
        </div>
    </div>
</div>

<style>
    .container {
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
    }

    .system-status {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
        margin-bottom: 2rem;
    }

    .status-line {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .blink {
        animation: blink 1s steps(1) infinite;
    }

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
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
    }

    input {
        width: 100%;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem 0.5rem;
        font-family: inherit;
        font-size: 0.9rem;
    }

    input:focus {
        outline: none;
    }

    .input-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .terminal-button {
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
    }

    .terminal-button:hover:not(:disabled) {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    .terminal-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .success-terminal {
        text-align: center;
        padding: 1rem;
    }

    .success-icon {
        color: #00b894;
        font-size: 2rem;
        margin-bottom: 1rem;
    }

    .success-content {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .email-highlight {
        color: #0984e3;
        font-size: 0.9rem;
        padding: 0.5rem;
        background: rgba(9, 132, 227, 0.1);
        border-radius: 3px;
        margin: 0.5rem 0;
    }

    .system-footer {
        margin-top: 2rem;
        padding-top: 1rem;
        border-top: 1px solid rgba(9, 132, 227, 0.2);
        text-align: center;
    }

    .system-link {
        color: #0984e3;
        text-decoration: none;
        font-size: 0.8rem;
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        transition: all 0.3s ease;
    }

    .system-link:hover {
        color: #00b894;
        text-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    @keyframes blink {
        0%, 50% { opacity: 1; }
        51%, 100% { opacity: 0; }
    }

    @media (max-width: 480px) {
        .container {
            margin: 20px auto;
        }

        input {
            font-size: 16px;
        }
    }
</style>
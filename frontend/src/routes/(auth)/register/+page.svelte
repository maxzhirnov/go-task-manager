<script>
    let email = '';
    let password = '';
    let errorMessage = '';
    let successMessage = '';
    let isRegistered = false;
    let loading = false;

    async function handleSubmit(e) {
        loading = true;
        try {
            const response = await fetch("/api/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Registration failed");
            }

            isRegistered = true;
            successMessage = "INITIALIZATION COMPLETE! VERIFY YOUR ACCESS NODE VIA EMAIL.";
            errorMessage = '';

        } catch (error) {
            errorMessage = error.message;
            successMessage = '';
        } finally {
            loading = false;
        }
    }

    async function handleResendVerification() {
        loading = true;
        try {
            const response = await fetch("/api/resend-verification", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Failed to resend verification email");
            }

            successMessage = "VERIFICATION SIGNAL TRANSMITTED. CHECK YOUR INBOX.";
            errorMessage = '';
        } catch (error) {
            errorMessage = error.message;
            successMessage = '';
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
            <span class="terminal-title">NEW_USER_REGISTRATION.exe</span>
        </div>
        
        <div class="terminal-content">
            <div class="system-status">
                <span class="status-line">INITIALIZING REGISTRATION PROTOCOL...</span>
                <span class="status-line">SECURE CONNECTION: ESTABLISHED</span>
                <span class="status-line blink">>_ AWAITING USER INPUT</span>
            </div>

            {#if !isRegistered}
                <form on:submit|preventDefault={handleSubmit}>
                    <div class="input-group">
                        <div class="input-label">[USER_EMAIL]</div>
                        <div class="input-wrapper">
                            <span class="prompt">>_</span>
                            <input 
                                type="email" 
                                bind:value={email} 
                                placeholder="Enter access email"
                                required
                                disabled={loading}
                            >
                        </div>
                    </div>

                    <div class="input-group">
                        <div class="input-label">[ACCESS_KEY]</div>
                        <div class="input-wrapper">
                            <span class="prompt">>_</span>
                            <input 
                                type="password" 
                                bind:value={password} 
                                placeholder="Enter access key"
                                required
                                minlength="6"
                                disabled={loading}
                            >
                        </div>
                        <div class="input-hint">Minimum 6 characters required</div>
                    </div>

                    <button type="submit" class="terminal-button" disabled={loading}>
                        <span class="btn-icon">⚡</span>
                        <span class="btn-text">
                            {loading ? 'INITIALIZING...' : 'INITIALIZE_REGISTRATION'}
                        </span>
                    </button>
                </form>
            {:else}
                <div class="success-terminal">
                    <div class="success-icon">✓</div>
                    <p class="success-text">{successMessage}</p>
                    <button 
                        class="terminal-button secondary"
                        on:click={handleResendVerification}
                        disabled={loading}
                    >
                        <span class="btn-icon">↻</span>
                        <span class="btn-text">
                            {loading ? 'TRANSMITTING...' : 'RETRANSMIT_VERIFICATION'}
                        </span>
                    </button>
                </div>
            {/if}

            {#if errorMessage}
                <div class="error-container">
                    <span class="error-prefix">[ERROR]</span>
                    <span class="error-message">{errorMessage}</span>
                </div>
            {/if}

            <div class="system-footer">
                <span class="footer-text">EXISTING_USER?</span>
                <a href="/login" class="system-link">ACCESS_LOGIN_PORTAL</a>
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

    .input-hint {
        color: #636e72;
        font-size: 0.7rem;
        margin-top: 0.3rem;
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

    .terminal-button.secondary {
        border-color: #0984e3;
        color: #0984e3;
    }

    .error-container {
        margin-top: 1rem;
        padding: 0.8rem;
        background: rgba(231, 76, 60, 0.1);
        border: 1px solid #e74c3c;
        border-radius: 3px;
        display: flex;
        gap: 0.5rem;
        font-size: 0.8rem;
    }

    .error-prefix {
        color: #e74c3c;
    }

    .error-message {
        color: #fff;
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

    .success-text {
        color: #00b894;
        font-size: 0.9rem;
        margin-bottom: 1.5rem;
        line-height: 1.6;
    }

    .system-footer {
        margin-top: 2rem;
        padding-top: 1rem;
        border-top: 1px solid rgba(9, 132, 227, 0.2);
        text-align: center;
        font-size: 0.7rem;
    }

    .footer-text {
        color: #636e72;
        margin-right: 0.5rem;
    }

    .system-link {
        color: #0984e3;
        text-decoration: none;
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
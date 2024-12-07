<!-- <script>
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { showError, showSuccess } from '$lib/stores';
    import { api } from '$lib/api';

    let newPassword = '';
    let confirmPassword = '';
    let loading = false;

    // Get token from URL
    const token = $page.url.searchParams.get('token');

    async function handleSubmit() {
        if (newPassword !== confirmPassword) {
            showError("Passwords don't match");
            return;
        }

        loading = true;
        try {
            await api.resetPassword({
                token,
                new_password: newPassword
            });
            showSuccess('Password has been reset successfully');
            
            // Check for refresh token and redirect accordingly
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                goto('/tasks');
            } else {
                goto('/login');
            }
        } catch (error) {
            showError(error.message);
        } finally {
            loading = false;
        }
    }
</script>

<div class="container">
    <h1>Reset Password</h1>
    <form on:submit|preventDefault={handleSubmit}>
        <input
            type="password"
            bind:value={newPassword}
            placeholder="New password"
            required
        />
        <input
            type="password"
            bind:value={confirmPassword}
            placeholder="Confirm new password"
            required
        />
        <button type="submit" disabled={loading}>
            {loading ? 'Resetting...' : 'Reset Password'}
        </button>
    </form>
</div>


<style>
    .container {
        font-family: Arial, sans-serif;
        max-width: 400px;
        margin: 50px auto;
        padding: 2rem;
        text-align: center;
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    h1 {
        color: #333;
        margin-bottom: 1.5rem;
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .form-group {
        position: relative;
        text-align: left;
    }

    input {
        width: 100%;
        padding: 0.75rem;
        font-size: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        box-sizing: border-box;
    }

    input:disabled {
        background: #f5f5f5;
        cursor: not-allowed;
    }

    button {
        margin-top: 1rem;
        padding: 0.75rem;
        font-size: 1rem;
        background-color: #4CAF50;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    button:hover:not(:disabled) {
        background-color: #45a049;
    }

    button:disabled {
        background-color: #cccccc;
        cursor: not-allowed;
    }

</style> -->

<script>
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { showError, showSuccess } from '$lib/stores';
    import { api } from '$lib/api';

    let newPassword = '';
    let confirmPassword = '';
    let loading = false;

    const token = $page.url.searchParams.get('token');

    async function handleSubmit() {
        if (newPassword !== confirmPassword) {
            showError("ACCESS_ERROR: Password sequences do not match");
            return;
        }

        loading = true;
        try {
            await api.resetPassword({
                token,
                new_password: newPassword
            });
            showSuccess('ACCESS_UPDATE: Password reset successful');
            
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                goto('/tasks');
            } else {
                goto('/login');
            }
        } catch (error) {
            showError(`SYSTEM_ERROR: ${error.message}`);
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
            <div class="system-status">
                <span class="status-line">INITIALIZING PASSWORD UPDATE PROTOCOL</span>
                <span class="status-line">TOKEN VERIFICATION: COMPLETE</span>
                <span class="status-line blink">>_ AWAITING NEW CREDENTIALS</span>
            </div>

            <form on:submit|preventDefault={handleSubmit}>
                <div class="input-group">
                    <div class="input-label">[NEW_ACCESS_KEY]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input
                            type="password"
                            bind:value={newPassword}
                            placeholder="Enter new access key"
                            required
                            disabled={loading}
                        />
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-label">[VERIFY_ACCESS_KEY]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input
                            type="password"
                            bind:value={confirmPassword}
                            placeholder="Confirm new access key"
                            required
                            disabled={loading}
                        />
                    </div>
                </div>

                <button type="submit" class="terminal-button" disabled={loading}>
                    <span class="btn-icon">⚡</span>
                    <span class="btn-text">
                        {loading ? 'UPDATING_ACCESS...' : 'EXECUTE_PASSWORD_UPDATE'}
                    </span>
                </button>
            </form>
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
<script>
    import { showError, showSuccess } from '$lib/stores';
    import { api } from '$lib/api';

    let email = '';
    let loading = false;
    let submitted = false;

    async function handleSubmit() {
        console.log('Starting password reset request');
        
        if (!email.trim()) {
            console.log('Email validation failed: empty email');
            showError('Please enter your email');
            return;
        }

        console.log('Attempting password reset for email:', email);
        loading = true;

        try {
            console.log('Sending password reset request to API');
            const response = await api.requestPasswordReset({ email });
            console.log('Password reset request response:', response);
            
            submitted = true;
            showSuccess('Password reset link has been sent to your email');
            console.log('Password reset request successful');
        } catch (error) {
            console.error('Password reset request failed:', {
                error: error,
                message: error.message,
                status: error.status,
                details: error.details
            });
            
            // More specific error messages based on error type
            if (error.status === 404) {
                showError('Email not found');
            } else if (error.status === 429) {
                showError('Too many reset attempts. Please try again later.');
            } else {
                showError('Failed to send reset link. Please try again.');
            }
        } finally {
            console.log('Password reset request completed');
            loading = false;
        }
    }
</script>

<div class="container">
    <h1>Reset Password</h1>
    
    {#if !submitted}
        <p class="description">
            Enter your email address and we'll send you a link to reset your password.
        </p>
        <form on:submit|preventDefault={handleSubmit}>
            <input 
                type="email" 
                bind:value={email} 
                placeholder="Enter your email"
                required
                disabled={loading}
            >
            <button type="submit" disabled={loading}>
                {loading ? 'Sending...' : 'Send Reset Link'}
            </button>
        </form>
    {:else}
        <div class="success-message">
            <p>We've sent a password reset link to <strong>{email}</strong></p>
            <p>Please check your email and follow the instructions to reset your password.</p>
        </div>
    {/if}

    <a href="/login" class="back-link">Back to Login</a>
</div>

<style>
    .container {
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
        margin-bottom: 1rem;
    }

    .description {
        color: #666;
        margin-bottom: 1.5rem;
    }

    form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    input {
        width: 100%;
        padding: 0.75rem;
        font-size: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        box-sizing: border-box;
    }

    button {
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

    .success-message {
        background: #f0f7f0;
        padding: 1.5rem;
        border-radius: 4px;
        margin: 1.5rem 0;
    }

    .success-message p {
        margin: 0.5rem 0;
        color: #2c5282;
    }

    .back-link {
        display: inline-block;
        margin-top: 1rem;
        color: #666;
        text-decoration: none;
    }

    .back-link:hover {
        color: #333;
        text-decoration: underline;
    }
</style>
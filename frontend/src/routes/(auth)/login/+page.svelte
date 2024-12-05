<script>
    import { goto } from '$app/navigation';
    
    let email = '';
    let password = '';
    let errorMessage = '';
    let loading = false;

    async function handleSubmit(e) {
        loading = true;
        errorMessage = '';
        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
                const error = await response.json();
                if (response.status === 403) {
                    errorMessage = "Please verify your email before logging in. Check your inbox for the verification link.";
                } else {
                    throw new Error(error.error || "Login failed");
                }
                return;
            }

            const { access_token, refresh_token } = await response.json();
            localStorage.setItem("jwt", access_token);
            localStorage.setItem("refresh_token", refresh_token);
            window.location.href = "/";
        } catch (error) {
            errorMessage = error.message;
        } finally {
            loading = false;
        }
    }
</script>

<div class="container">
    <h1>Login</h1>
    <form on:submit|preventDefault={handleSubmit}>
        <div class="form-group">
            <input 
                type="email" 
                bind:value={email} 
                placeholder="Email" 
                required
                disabled={loading}
            >
        </div>
        <div class="form-group">
            <input 
                type="password" 
                bind:value={password} 
                placeholder="Password" 
                required
                disabled={loading}
            >
            <a href="/forgot-password" class="forgot-password">
                Forgot password?
            </a>
        </div>
        <button type="submit" disabled={loading}>
            {loading ? 'Logging in...' : 'Login'}
        </button>
    </form>
    {#if errorMessage}
        <p class="error">{errorMessage}</p>
    {/if}
    <p class="register-link">
        Don't have an account? <a href="/register">Register here</a>
    </p>
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

    .forgot-password {
        position: absolute;
        right: 0;
        top: 100%;
        font-size: 0.8rem;
        color: #666;
        text-decoration: none;
        margin-top: 0.25rem;
    }

    .forgot-password:hover {
        color: #333;
        text-decoration: underline;
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

    .error {
        color: #e74c3c;
        margin-top: 1rem;
        font-size: 0.9rem;
    }

    .register-link {
        margin-top: 1.5rem;
        color: #666;
    }

    .register-link a {
        color: #4CAF50;
        text-decoration: none;
    }

    .register-link a:hover {
        text-decoration: underline;
    }
</style>
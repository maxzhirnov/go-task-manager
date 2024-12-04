<script>
    let email = '';
    let password = '';
    let errorMessage = '';
    let successMessage = '';
    let isRegistered = false;

    async function handleSubmit(e) {
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
            successMessage = "Registration successful! Please check your email to verify your account.";
            errorMessage = '';
        } catch (error) {
            errorMessage = error.message;
            successMessage = '';
        }
    }

    async function handleResendVerification() {
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

            successMessage = "Verification email sent! Please check your inbox.";
            errorMessage = '';
        } catch (error) {
            errorMessage = error.message;
            successMessage = '';
        }
    }
</script>

<div class="container">
    <h1>Register</h1>
    
    {#if !isRegistered}
        <form on:submit|preventDefault={handleSubmit}>
            <input 
                type="email" 
                bind:value={email} 
                placeholder="Email" 
                required
            >
            <input 
                type="password" 
                bind:value={password} 
                placeholder="Password" 
                required
                minlength="6"
            >
            <button type="submit">Register</button>
        </form>
    {:else}
        <div class="success-message">
            <p>{successMessage}</p>
            <button on:click={handleResendVerification}>
                Resend Verification Email
            </button>
        </div>
    {/if}

    {#if errorMessage}
        <p class="error">{errorMessage}</p>
    {/if}
    
    {#if successMessage && !errorMessage}
        <p class="success">{successMessage}</p>
    {/if}

    <p>Already have an account? <a href="/login">Login here</a></p>
</div>

<style>
    .container {
        font-family: Arial, sans-serif;
        max-width: 400px;
        margin: 50px auto;
        text-align: center;
    }
    form {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }
    input {
        padding: 10px;
        font-size: 16px;
    }
    button {
        padding: 10px;
        font-size: 16px;
        background-color: #4CAF50;
        color: white;
        border: none;
        cursor: pointer;
    }
    button:hover {
        background-color: #45a049;
    }
    .error {
        color: red;
    }
    .success-message {
        background-color: #dff0d8;
        padding: 20px;
        border-radius: 5px;
        margin: 20px 0;
    }
    .success {
        color: #3c763d;
    }
</style>
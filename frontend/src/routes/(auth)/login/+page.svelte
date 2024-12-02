<script>
    let username = '';
    let password = '';
    let errorMessage = '';

    async function handleSubmit(e) {
        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || "Login failed");
            }

            const { access_token, refresh_token } = await response.json();
            localStorage.setItem("jwt", access_token);
            localStorage.setItem("refresh_token", refresh_token);
            window.location.href = "/";
        } catch (error) {
            errorMessage = error.message;
        }
    }
</script>

<div class="container">
    <h1>Login</h1>
    <form on:submit|preventDefault={handleSubmit}>
        <input type="text" bind:value={username} placeholder="Username" required>
        <input type="password" bind:value={password} placeholder="Password" required>
        <button type="submit">Login</button>
    </form>
    {#if errorMessage}
        <p class="error">{errorMessage}</p>
    {/if}
    <p>Don't have an account? <a href="/register">Register here</a></p>
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
</style>
<script>
    let username = '';
    let password = '';
    let errorMessage = '';

    async function handleSubmit(e) {
        try {
            const response = await fetch("/api/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Registration failed");
            }

            alert("Registration successful! Redirecting to login...");
            window.location.href = "/login";
        } catch (error) {
            errorMessage = error.message;
        }
    }
</script>

<div class="container">
    <h1>Register</h1>
    <form on:submit|preventDefault={handleSubmit}>
        <input type="text" bind:value={username} placeholder="Username" required>
        <input type="password" bind:value={password} placeholder="Password" required>
        <button type="submit">Register</button>
    </form>
    {#if errorMessage}
        <p class="error">{errorMessage}</p>
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
</style>
<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';

    let status = 'verifying';
    let message = 'Verifying your email...';

    onMount(async () => {
        const token = $page.url.searchParams.get('token');
        if (!token) {
            status = 'error';
            message = 'Verification token is missing';
            return;
        }

        try {
            const response = await fetch(`/api/verify-email?token=${token}`);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Verification failed');
            }

            status = 'success';
            message = 'Email verified successfully! You can now login.';
        } catch (error) {
            status = 'error';
            message = error.message;
        }
    });
</script>

<div class="container">
    <h1>Email Verification</h1>
    
    {#if status === 'verifying'}
        <div class="loading">
            <p>{message}</p>
        </div>
    {:else if status === 'success'}
        <div class="success">
            <p>{message}</p>
            <a href="/login" class="button">Go to Login</a>
        </div>
    {:else}
        <div class="error">
            <p>{message}</p>
            <button on:click={() => window.location.reload()}>Try Again</button>
        </div>
    {/if}
</div>

<style>
    .container {
        font-family: Arial, sans-serif;
        max-width: 400px;
        margin: 50px auto;
        text-align: center;
    }
    .loading, .success, .error {
        padding: 20px;
        border-radius: 5px;
        margin: 20px 0;
    }
    .success {
        background-color: #dff0d8;
        color: #3c763d;
    }
    .error {
        background-color: #f2dede;
        color: #a94442;
    }
    .button {
        display: inline-block;
        padding: 10px 20px;
        background-color: #4CAF50;
        color: white;
        text-decoration: none;
        border-radius: 5px;
        margin-top: 20px;
    }
    .button:hover {
        background-color: #45a049;
    }
</style>
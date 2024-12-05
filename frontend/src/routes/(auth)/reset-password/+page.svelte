<script>
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
            goto('/login');
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

</style>
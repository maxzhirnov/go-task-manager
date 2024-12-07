<script>
    import { onMount } from 'svelte';
    import { user, showError, showSuccess, parseJWT } from '$lib/stores';
    import { api } from '$lib/api';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import UserStatistics from '$lib/components/UserStatistics.svelte';
    import Logo12 from '$lib/components/logos/Logo12.svelte';
    import Logo11 from '$lib/components/logos/Logo11.svelte';

    let loading = false;

    

    let formData = {
        username: $user?.username || '',
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
        email: ''
    };

    let formErrors = {
        username: '',
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
        email: ''
    };


    async function updateUsername() {
        //reset form errors
        formErrors = {
            username: '',
            currentPassword: '',
            newPassword: '',
            confirmPassword: ''
        };

        console.log('Starting updateUsername function');
        if (!formData.username.trim()) {
            formErrors.username = 'Username cannot be empty';
            showError("Username cannot be empty");
            return;
        }

        if (!formData.currentPassword.trim()) {
            formErrors.currentPassword = 'Current password cannot be empty';
            showError("Current password cannot be empty");
            return;
        }

        // loading = true;
        try {
            // Update profile
            await api.updateProfile({ 
                username: formData.username,
                current_password: formData.currentPassword 
            });

            // Refresh token after successful update
            const newToken = await api.refreshToken();
            if (!newToken) {
                console.error('Failed to refresh token after username update');
                return;
            }

            // Update user store with new data
            const userData = parseJWT(newToken);
            if (userData) {
                user.set({
                    id: userData.user_id,
                    username: userData.username,
                    email: userData.email
                });
            }

            showSuccess("Username updated successfully");
        } catch (error) {
            console.error('Update failed:', error);
            
            if (error.status === 401) {
                formErrors.currentPassword = 'Invalid current password';
            } else if (error.status === 400) {
                if (error.details?.error) {
                    formErrors.username = error.details.error;
                }
            } else {
                showError("Failed to update username");
            }
        } finally {
            loading = false;
        }
    }

    async function updatePassword() {
        console.log('Starting updatePassword function');
        
        // Reset form errors
        formErrors = {
            username: '',
            currentPassword: '',
            newPassword: '',
            confirmPassword: ''
        };

        // Validate passwords match
        if (formData.newPassword !== formData.confirmPassword) {
            console.log('Password mismatch');
            formErrors.confirmPassword = "Passwords don't match";
            return;
        }

        // Validate password is not empty
        if (!formData.newPassword.trim()) {
            console.log('New password is empty');
            formErrors.newPassword = "New password cannot be empty";
            return;
        }

        // loading = true;

        try {
            console.log('Sending password update request');
            await api.updateProfile({ 
                new_password: formData.newPassword,
                current_password: formData.currentPassword 
            });
            
            console.log('Password update successful');
            showSuccess("Password updated successfully");
            
            // Reset form
            showPasswordForm = false;
            formData.currentPassword = '';
            formData.newPassword = '';
            formData.confirmPassword = '';
        } catch (error) {
            console.error('Password update failed:', error);
            
            if (error.status === 401) {
                formErrors.currentPassword = 'Invalid current password';
            } else if (error.status === 400) {
                if (error.details?.error) {
                    formErrors.newPassword = error.details.error;
                }
            } else {
                showError("Failed to update password");
            }
        } finally {
            console.log('Password update process completed');
            loading = false;
        }
    }

    async function updatePasswordThroughEmail(email) {
        
        if (!email.trim()) {
            formErrors.email = 'Please enter your email';
            return;
        }

        loading = true;

        try {
            const response = await api.requestPasswordReset({ email });
            showSuccess('Password reset link has been sent to your email');
        } catch (error) {
            formErrors.email = 'Failed to send reset link. Please try again.'
            showError('Failed to send reset link. Please try again.');
            if (error.status === 404) {
                formErrors.email = 'Email not found';
                showError('Email not found');
            } else if (error.status === 429) {
                formErrors.email = 'Too many reset attempts. Please try again later.';
                showError('Too many reset attempts. Please try again later.');
            } else {
                formErrors.email = 'Failed to send reset link. Please try again.';
                showError('Failed to send reset link. Please try again.');
            }
        } finally {
            formErrors.email = '';
            loading = false;
        }
    }
</script>

<!-- <div class="container">
    <div class="profile-header">
        <h1>My Profile</h1>
        <div class="email">{$user?.email}</div>
    </div>

    {#if loading}
        <div class="loading">
            <LoadingSpinner />
        </div>
    {:else}
        <div class="card">
            <div class="card-header">
                <h2>Username</h2>
            </div>
            <div class="card-content">
                <div class="form-group">
                    <input 
                        type="text" 
                        bind:value={formData.username} 
                        placeholder="Enter new username"
                        class:error={formErrors.username}
                    />
                    {#if formErrors.username}
                        <span class="error-message">{formErrors.username}</span>
                    {/if}
                </div>
                
                <div class="form-group">
                    <input 
                        type="password" 
                        bind:value={formData.currentPassword} 
                        placeholder="Enter current password"
                        class:error={formErrors.currentPassword}
                    />
                    {#if formErrors.currentPassword}
                        <span class="error-message">{formErrors.currentPassword}</span>
                    {/if}
                </div>
                
                <button class="primary-button" on:click={updateUsername}>
                    Update Username
                </button>
            </div>
        </div>

        <div class="card">
            <div class="card-header">
                <h2>Password</h2>
            </div>
            <div class="card-content">
                <div class="form-group">
                    <input 
                        type="email"
                        bind:value={formData.email} 
                        placeholder="Enter your login email"
                        class:error={formErrors.email}
                    />
                    {#if formErrors.email}
                        <span class="error-message">{formErrors.email}</span>
                    {/if}
                </div>
  
                <button class="primary-button" on:click={() => updatePasswordThroughEmail(formData.email)}>
                    Change Password
                </button>
            </div>
        </div>

        <Logo12 clickable={false} />
    {/if}
</div>

<style>
    .form-group {
        margin-bottom: 1.2rem;
        position: relative;
    }

    input.error {
        border-color: #e74c3c;
    }

    .error-message {
        color: #e74c3c;
        font-size: 0.8rem;
        margin-top: 0.25rem;
        position: absolute;
        bottom: -1.01rem;
        left: 0;
    }
    .container {
        max-width: 600px;
        margin: 0 auto;
        padding: 2rem;
    }

    .profile-header {
        text-align: center;
        margin-bottom: 2rem;
    }

    .profile-header h1 {
        margin: 0;
        color: #333;
        font-size: 2rem;
    }

    .email {
        color: #666;
        margin-top: 0.5rem;
    }

    .card {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        margin-bottom: 1.5rem;
        overflow: hidden;
    }

    .card-header {
        background: #f8f9fa;
        padding: 1rem 1.5rem;
        border-bottom: 1px solid #eee;
    }

    .card-header h2 {
        margin: 0;
        font-size: 1.2rem;
        color: #333;
    }

    .card-content {
        padding: 1.5rem;
    }

    input {
        width: 90%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
    }

    .button-group {
        display: flex;
        justify-content: flex-end;
        gap: 1rem;
        margin-top: 1rem;
    }

    .primary-button {
        background: #4CAF50;
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1rem;
    }

    .secondary-button {
        background: #f8f9fa;
        color: #333;
        border: 1px solid #ddd;
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1rem;
    }

    .text-button {
        background: none;
        border: none;
        color: #666;
        cursor: pointer;
        font-size: 1rem;
        padding: 0.75rem 1rem;
    }

    .loading {
        display: flex;
        justify-content: center;
        padding: 2rem;
    }

    @media (max-width: 640px) {
        .container {
            padding: 1rem;
        }
    }
</style> -->

<div class="container">
    <div class="terminal-box profile-header">
        <div class="terminal-header">
            <span class="terminal-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </span>
            <span class="terminal-title">USER_PROFILE.sys</span>
        </div>
        <div class="profile-content">
            <div class="user-id">ID: {$user?.email}</div>
            <div class="status-line">STATUS: AUTHENTICATED</div>
        </div>
    </div>

    {#if loading}
        <div class="loading">
            <LoadingSpinner />
        </div>
    {:else}
        <div class="terminal-box">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">USERNAME_CONFIG.exe</span>
            </div>
            <div class="terminal-content">
                <div class="form-group">
                    <div class="input-label">[USERNAME]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            type="text" 
                            bind:value={formData.username} 
                            placeholder="New username"
                            class:error={formErrors.username}
                        />
                    </div>
                    {#if formErrors.username}
                        <span class="error-message">[ERROR] {formErrors.username}</span>
                    {/if}
                </div>
                
                <div class="form-group">
                    <div class="input-label">[VERIFY]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            type="password" 
                            bind:value={formData.currentPassword} 
                            placeholder="Current password"
                            class:error={formErrors.currentPassword}
                        />
                    </div>
                    {#if formErrors.currentPassword}
                        <span class="error-message">[ERROR] {formErrors.currentPassword}</span>
                    {/if}
                </div>
                
                <button class="terminal-button" on:click={updateUsername}>
                    <span class="btn-icon">⚡</span>
                    <span class="btn-text">UPDATE_USERNAME</span>
                </button>
            </div>
        </div>

        <div class="terminal-box">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">PASSWORD_CONFIG.exe</span>
            </div>
            <div class="terminal-content">
                <div class="form-group">
                    <div class="input-label">[VERIFY_EMAIL]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            type="email"
                            bind:value={formData.email} 
                            placeholder="Confirm email"
                            class:error={formErrors.email}
                        />
                    </div>
                    {#if formErrors.email}
                        <span class="error-message">[ERROR] {formErrors.email}</span>
                    {/if}
                </div>
  
                <button class="terminal-button" on:click={() => updatePasswordThroughEmail(formData.email)}>
                    <span class="btn-icon">🔒</span>
                    <span class="btn-text">INITIATE_PASSWORD_RESET</span>
                </button>
            </div>
        </div>

        <Logo12 clickable={false} />
    {/if}
</div>

<style>
    .container {
        max-width: 800px;
        margin: 2rem auto;
        padding: 0 1rem;
        font-family: "JetBrains Mono", monospace;
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        margin-bottom: 2rem;
        position: relative;
        overflow: hidden;
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
        padding: 1rem;
    }

    .profile-content {
        padding: 1rem;
        color: #00b894;
    }

    .user-id {
        font-size: 0.9rem;
        margin-bottom: 0.5rem;
    }

    .status-line {
        font-size: 0.8rem;
        color: #0984e3;
    }

    .form-group {
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

    .error-message {
        color: #ff6b6b;
        font-size: 0.8rem;
        margin-top: 0.5rem;
    }

    .terminal-button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.8rem 1rem;
        border-radius: 3px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: inherit;
        font-size: 0.8rem;
        transition: all 0.3s ease;
        width: 100%;
        justify-content: center;
    }

    .terminal-button:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    .loading {
        display: flex;
        justify-content: center;
        padding: 2rem;
    }

    @media (max-width: 600px) {
        .container {
            margin: 1rem auto;
        }

        .terminal-content {
            padding: 0.8rem;
        }
    }
</style>
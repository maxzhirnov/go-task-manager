<script>
    import { onMount } from 'svelte';
    import { user, showError, showSuccess, parseJWT } from '$lib/stores';
    import { api } from '$lib/api';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import UserStatistics from '$lib/components/UserStatistics.svelte';

    let loading = false;
    let showPasswordForm = false;
    let formData = {
        username: $user?.username || '',
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
    };

    let formErrors = {
        username: '',
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
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
</script>

<div class="container">
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
                {#if !showPasswordForm}
                    <button 
                        class="secondary-button" 
                        on:click={() => showPasswordForm = true}
                    >
                        Change Password
                    </button>
                {:else}
                    <form on:submit|preventDefault={updatePassword}>
                        <div class="form-group">
                            <input
                                type="password"
                                bind:value={formData.currentPassword}
                                placeholder="Current password"
                                class:error={formErrors.currentPassword}
                                required
                            />
                            {#if formErrors.currentPassword}
                                <span class="error-message">{formErrors.currentPassword}</span>
                            {/if}
                        </div>
        
                        <div class="form-group">
                            <input
                                type="password"
                                bind:value={formData.newPassword}
                                placeholder="New password"
                                class:error={formErrors.newPassword}
                                required
                            />
                            {#if formErrors.newPassword}
                                <span class="error-message">{formErrors.newPassword}</span>
                            {/if}
                        </div>
        
                        <div class="form-group">
                            <input
                                type="password"
                                bind:value={formData.confirmPassword}
                                placeholder="Confirm new password"
                                class:error={formErrors.confirmPassword}
                                required
                            />
                            {#if formErrors.confirmPassword}
                                <span class="error-message">{formErrors.confirmPassword}</span>
                            {/if}
                        </div>
        
                        <div class="button-group">
                            <button 
                                type="button" 
                                class="text-button"
                                on:click={() => showPasswordForm = false}
                            >
                                Cancel
                            </button>
                            <button type="submit" class="primary-button">
                                Update Password
                            </button>
                        </div>
                    </form>
                {/if}
            </div>
        </div>


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
</style>
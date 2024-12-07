<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    import ErrorMessage from '$lib/components/ErrorMessage.svelte';
    import Notification from '$lib/components/Notification.svelte';
    
    let isLoading = true;

    onMount(() => {
        const token = localStorage.getItem('jwt');
        const isResetPasswordPage = window.location.pathname.startsWith('/reset-password');
        
        if (token && !isResetPasswordPage) {
            goto('/tasks');
        } else {
            isLoading = false;
        }
    });
</script>

{#if isLoading}
    <LoadingSpinner/>
    
{:else}
    <div class="auth-layout">
        <ErrorMessage/>
        <Notification/>
        <slot />
    </div>
{/if}

<style>
    .auth-layout {
        font-family: 'JetBrains Mono', monospace;
    }
</style>
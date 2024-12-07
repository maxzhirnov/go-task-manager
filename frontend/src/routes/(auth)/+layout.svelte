<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    
    import LoadingSpinner from '$lib/components/Common/LoadingSpinner.svelte';
    import ErrorMessage from '$lib/components/Common/ErrorMessage.svelte';
    import Notification from '$lib/components/Common/Notification.svelte';
    
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

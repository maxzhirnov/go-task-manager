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
    :global(body) {
        color: white;
        max-width: 800px;
        margin: 0 auto 2rem;
        padding: 20px;
        background-color: #0f1215;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                #161b22 1px,
                transparent 1px
            );
        background-size: 24px 24px;
    }
</style>
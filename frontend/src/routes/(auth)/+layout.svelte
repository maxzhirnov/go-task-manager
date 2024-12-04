<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
    
    let isLoading = true;

    onMount(() => {
        const token = localStorage.getItem('jwt');
        if (token) {
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
        <slot />
    </div>
{/if}
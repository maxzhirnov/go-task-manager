<script>
    import ErrorMessage from '$lib/components/ErrorMessage.svelte';
    import Header from '$lib/components/Header.svelte';
    import { onMount } from 'svelte';
    import { initializeUser } from '$lib/auth.js';
    import { page } from '$app/stores'

    // Array of routes where header should be hidden
    const noHeaderRoutes = ['/login', '/register'];

    // Check if current path should show header
    $: showHeader = !noHeaderRoutes.includes($page.url.pathname);

    onMount(() => {
        // Only initialize user if we're not on login/register pages
        if (showHeader) {
            initializeUser();
        }
    })
</script>

{#if showHeader}
    <Header />
{/if}
<ErrorMessage />
<slot />

<style>
    :global(body) {
        font-family: Helvetica, Arial, sans-serif;
        box-sizing: border-box;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
    }
  </style>
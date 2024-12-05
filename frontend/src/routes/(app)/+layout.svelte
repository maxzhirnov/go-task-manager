<script>
    import ErrorMessage from '$lib/components/ErrorMessage.svelte';
    import Header from '$lib/components/Header.svelte';
    import { onMount } from 'svelte';
    import { initializeUser } from '$lib/stores.js';
    import Notification from '$lib/components/Notification.svelte';

    let appName = "ActionHub";
    let pageTitle = $$props.title || "Tasks";

    onMount(() => {
        // Only initialize user if we're not on login/register pages
        initializeUser();
    })
</script>

<svelte:head>
    <title>{pageTitle} | {appName}</title>
</svelte:head>

<Header />
<ErrorMessage />

<slot />

<Notification />

<style>
    :global(body) {
        font-family: Helvetica, Arial, sans-serif;
        box-sizing: border-box;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
    }
  </style>
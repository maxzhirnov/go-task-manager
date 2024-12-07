<script>
    import { onMount } from 'svelte';
    import { initializeUser } from '$lib/stores.js';

    import Header from '$lib/components/Header/Header.svelte';
    import ErrorMessage from '$lib/components/Common/ErrorMessage.svelte';
    import Notification from '$lib/components/Common/Notification.svelte';

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

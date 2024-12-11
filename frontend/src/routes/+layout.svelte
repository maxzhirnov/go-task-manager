<script>
    import '../styles/global.css'

    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { Analytics } from '$lib/analytics';
    import { user } from '$lib/stores.js';

    let isUserTracked = false;

    $: if ($user && !isUserTracked) {
        Analytics.trackUser($user);
        isUserTracked = true;
    }

    $: if (!$user) {
        isUserTracked = false;
    }

    // Track page views
    // $: if (typeof window !== 'undefined') {
    //     Analytics.pageView($page.url.pathname);
    // }

    onMount(() => {
        Analytics.track('App Loaded');
    });
</script>

<div class="container">
    <slot />
</div>
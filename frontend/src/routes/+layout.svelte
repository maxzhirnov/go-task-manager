<script>
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import '../styles/global.css'
    import { Analytics } from '$lib/analytics';

    $: if (typeof window !== 'undefined') {
        Analytics.pageView($page.url.pathname);
    }

    onMount(() => {
        Analytics.track('App Loaded');
    });
</script>

<div class="container">
    <slot />
</div>
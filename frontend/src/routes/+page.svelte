<script>
    import LandingPage from "$lib/components/LandingPage.svelte";
    import { browser } from '$app/environment';
    import { goto } from '$app/navigation';

    if (browser && localStorage.getItem('jwt')) {
        goto('/tasks');
    }

</script>
<LandingPage/>

<style>
    /* :global(body) {
        margin: 0;
    } */
</style>
<script>
    import { onMount } from 'svelte';
    import UserStatistics from '$lib/components/UserStatistics.svelte';
    import { goto } from '$app/navigation';
    import Header from '$lib/components/Header.svelte';

    let statistics = null;

    onMount(async () => {
        // Check if user is authenticated
        const jwt = localStorage.getItem('jwt');
        if (!jwt) {
            goto('/login');
            return;
        }
    });
</script>

<svelte:head>
    <title>Task Statistics</title>
</svelte:head>

<UserStatistics {statistics} />
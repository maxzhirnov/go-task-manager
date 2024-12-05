<script>
    import { user } from '../stores.js';
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';


    function logout() {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        window.location.href = "/login";
    }
</script>

<div class="header-container">
    <div class="title">
        <a href="/tasks" class="title-link" on:click|preventDefault={() => goto('/tasks')}>
            <h1>ActionHub</h1>
        </a>
    </div>
    <div class="nav-section">
        <nav class="nav-links">
            <a class="nav-link" class:active={$page.url.pathname === '/'} href="/" on:click|preventDefault={() => goto('/tasks')}>Tasks</a>
            <a class="nav-link" class:active={$page.url.pathname === '/statistics'} href="/statistics" on:click|preventDefault={() => goto('/statistics')}>Statistics</a>
        </nav>
    </div>
    <div class="user-info">
        {#if $user}
            <span>Welcome, 
                <a href="/profile" class="profile-link" on:click|preventDefault={() => goto('/profile')}>
                    {$user.username}
                </a>
            </span>
        {/if}
        <button class="logout-btn" on:click={logout}>Logout</button>
    </div>
</div>

<style>
    h1 {
        margin: 0;
        font-size: 24px;
    }

    .title-link {
        text-decoration: none;
        color: inherit;
        cursor: pointer;
    }

    .profile-link {
        /* text-decoration: none; */
        color: inherit;
        cursor: pointer;
    }

    .title-link:hover {
        opacity: 0.7;
    }

    .header-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 15px 0px;
        gap: 20px;
        width: 100%;
    }

    .user-info {
        display: flex;
        align-items: center;
        gap: 15px;
        white-space: nowrap;
    }

    .logout-btn {
        background-color: #f44336;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
    }

    .logout-btn:hover {
        background-color: #da190b;
    }

    .nav-section {
        display: flex;
        align-items: center;
        gap: 20px;
        flex: 1;
    }

    .nav-links {
        display: flex;
        gap: 15px;
        margin-left: 20px;
    }

    .nav-link {
        text-decoration: none;
        color: #666;
        padding: 5px 10px;
        border-radius: 4px;
        transition: all 0.2s;
    }

    .nav-link:hover {
        background-color: #f5f5f5;
    }

    .nav-link.active {
        color: #2196F3;
        font-weight: bold;
    }

    @media screen and (max-width: 600px) {
        .title { 
            order: 2; 
        }
        
        .header-container {
            flex-direction: column;
            text-align: center;
            margin-top: 0px;
            padding-top: 0px;
        }

        .user-info {
            order: 1;
            flex-direction: row;
            width: 100%;
            justify-content: space-between;
        }

        .logout-btn {
            width: 100px;
            padding: 10px;
        }

        .nav-section {
            order: 3;
            flex-direction: column;
            width: 100%;
        }

        .nav-links {
            margin: 10px 0;
            width: 100%;
            justify-content: center;
        }
    }
</style>
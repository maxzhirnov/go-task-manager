<!-- <script>
    import { user } from '../stores.js';
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import Logo14 from './logos/Logo14.svelte';


    function logout() {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        window.location.href = "/login";
    }
</script>

<div class="header-container">
    <div class="title">
        <Logo14 />
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
</style> -->

<script>
    import { user } from '../stores.js';
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import Logo14 from './logos/Logo14.svelte';

    function logout() {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        window.location.href = "/login";
    }
</script>

<div class="header-container">
    <div class="header-content">
        <div class="title">
            <Logo14 />
        </div>
        <div class="nav-section">
            <div class="terminal-box">
                <div class="terminal-header">
                    <span class="terminal-dots">
                        <span class="dot"></span>
                        <span class="dot"></span>
                        <span class="dot"></span>
                    </span>
                    <span class="terminal-title">NAVIGATION.sys</span>
                </div>
                <nav class="nav-links">
                    <a 
                        class="nav-link" 
                        class:active={$page.url.pathname === '/'} 
                        href="/" 
                        on:click|preventDefault={() => goto('/tasks')}
                    >
                        <span class="prompt">>_</span>
                        <span class="link-text">TASKS</span>
                    </a>
                    <a 
                        class="nav-link" 
                        class:active={$page.url.pathname === '/statistics'} 
                        href="/statistics" 
                        on:click|preventDefault={() => goto('/statistics')}
                    >
                        <span class="prompt">>_</span>
                        <span class="link-text">STATS</span>
                    </a>
                </nav>
            </div>
        </div>
        <div class="user-info">
            {#if $user}
                <div class="user-terminal">
                    <span class="user-status">USER::</span>
                    <a href="/profile" class="profile-link" on:click|preventDefault={() => goto('/profile')}>
                        {$user.username}
                    </a>
                    <button class="logout-btn" on:click={logout}>
                        <span class="btn-icon">⏻</span>
                        <span class="btn-text">LOGOUT</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .header-container {
        background: #1c1c1c;
        border-bottom: 1px solid #0984e3;
        padding: 1rem 0;
        position: relative;
        font-family: "JetBrains Mono", monospace;
    }

    .header-container::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.05) 1px,
                transparent 1px
            );
        background-size: 10px 10px;
        pointer-events: none;
    }

    .header-content {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 1rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 2rem;
    }

    .terminal-box {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
    }

    .terminal-header {
        background: #1c1c1c;
        padding: 0.3rem 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
    }

    .nav-links {
        display: flex;
        gap: 1rem;
        padding: 0.5rem;
    }

    .nav-link {
        text-decoration: none;
        color: #00b894;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.4rem 0.8rem;
        border-radius: 3px;
        transition: all 0.3s ease;
    }

    .nav-link:hover {
        background: rgba(0, 184, 148, 0.1);
    }

    .nav-link.active {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .prompt {
        color: #0984e3;
        font-size: 0.8rem;
    }

    .link-text {
        font-size: 0.8rem;
    }

    .user-info {
        display: flex;
        align-items: center;
    }

    .user-terminal {
        display: flex;
        align-items: center;
        gap: 1rem;
        color: #00b894;
        font-size: 0.8rem;
    }

    .user-status {
        color: #636e72;
    }

    .profile-link {
        color: #0984e3;
        text-decoration: none;
        transition: all 0.3s ease;
    }

    .profile-link:hover {
        color: #00b894;
        text-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    .logout-btn {
        background: transparent;
        border: 1px solid #ff6b6b;
        color: #ff6b6b;
        padding: 0.4rem 0.8rem;
        border-radius: 3px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
        transition: all 0.3s ease;
    }

    .logout-btn:hover {
        background: rgba(255, 107, 107, 0.1);
        box-shadow: 0 0 8px rgba(255, 107, 107, 0.3);
    }

    @media screen and (max-width: 768px) {
        .header-content {
            flex-direction: column;
            gap: 1rem;
        }

        .nav-links {
            justify-content: center;
        }

        .user-terminal {
            flex-direction: column;
            align-items: center;
            gap: 0.8rem;
        }

        .logout-btn {
            width: 100%;
            justify-content: center;
        }
    }
</style>
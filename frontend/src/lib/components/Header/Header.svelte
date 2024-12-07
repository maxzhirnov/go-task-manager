<script>
    import { user } from '$lib/stores.js';

    import Logo from '$lib/components/logos/Logo15.svelte';
    import Navigation from './Navigation2.svelte';
    import UserInfo from './UserInfo.svelte';
</script>

<div class="header-container">
    <div class="header-content">
        <div class="logo">
            <Logo />
        </div>

        <div class="nav-section">
            <Navigation />
        </div>

        <div class="user-info">
            <UserInfo {user} />
        </div>
    </div>
</div>

<style>
    .header-container {
        margin-bottom: 1rem;
        background: #1c1c1c;
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

    .user-info {
        display: flex;
        align-items: center;
    }

    @media screen and (max-width: 768px) {
        .header-content {
            flex-direction: column;
            gap: 1rem;
        }
        
        .nav-section {
            order: 2;
        }

        .user-info {
            order: 1;
        }
    }
</style>
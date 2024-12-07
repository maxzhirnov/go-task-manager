<script>
    export let user;
    import { goto } from '$app/navigation';

    function logout() {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        window.location.href = "/login";
    }
</script>

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

<style>
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
        .user-terminal {
            /* flex-direction: column; */
            align-items: center;
            gap: 0.8rem;
        }

        .logout-btn {
            width: 100%;
            justify-content: center;
        }
    }
</style>
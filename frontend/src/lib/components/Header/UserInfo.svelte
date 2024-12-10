<script>
    export let user;
    import { goto } from '$app/navigation';
    import { Analytics } from '$lib/analytics';

    function logout() {
        Analytics.track('User Logged Out');
        Analytics.clearUser();

        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        goto('/login');
    }
</script>

{#if $user}
    <div class="user-terminal">
        <span class="user-status">USER::</span>
        <a href="/profile" class="profile-link" on:click|preventDefault={() => goto('/profile')}>
            {$user.username}
        </a>
        <button class="logout-btn" on:click={logout}>
            <svg class="btn-icon" viewBox="0 0 24 24" width="16" height="16">
                <!-- Outer circle -->
                <circle 
                    cx="12" 
                    cy="12" 
                    r="10" 
                    fill="none" 
                    stroke="currentColor" 
                    stroke-width="1"
                    opacity="0.5"
                    class="outer-circle"
                />
                <!-- Power circle -->
                <circle 
                    cx="12" 
                    cy="12" 
                    r="6" 
                    fill="none" 
                    stroke="currentColor" 
                    stroke-width="1.5"
                    class="power-circle"
                />
                <!-- Power line -->
                <path 
                    d="M12 8v4" 
                    stroke="currentColor" 
                    stroke-width="1.5" 
                    stroke-linecap="round"
                    class="power-line"
                />
            </svg>
            <span class="btn-text">LOGOUT</span>
        </button>
    </div>
{/if}

<style>
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

    .btn-icon {
        transition: transform 0.3s ease;
    }

    .outer-circle, .power-circle, .power-line {
        transition: all 0.3s ease;
    }

    .logout-btn:hover .btn-icon {
        transform: rotate(180deg);
    }

    .logout-btn:hover .power-circle {
        stroke-width: 2;
    }

    .logout-btn:hover .power-line {
        stroke-width: 2;
    }

    .logout-btn:hover .outer-circle {
        opacity: 0.8;
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
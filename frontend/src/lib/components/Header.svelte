<script>
    import { user } from '../stores.js';

    function logout() {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        window.location.href = "/login";
    }
</script>

<div class="header-container">
    <a href="/" class="title-link" on:click={window.location.reload()}>
        <h1>Task Manager</h1>
    </a>
    <div class="user-info">
        <span>Welcome, {$user}</span>
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

    @media screen and (max-width: 600px) {
        .title-link { 
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
    }
</style>